package main

import (
    "context"
    "fmt"
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
// log used for showing requiest status while using API
// net/http package will start the server, and gorillamux for the routing

//Trainer type
type Trainer struct {
    Name string
    Age  int
    City string
}

// Set client options
var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
// Connect to MongoDB
var client, err = mongo.Connect(context.TODO(), clientOptions)
//can now get a handle for the trainers collection
// will later use this handle to query the collection
var collection = client.Database("test").Collection("trainers")


// Get all trainers
func getTrainers(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  findOptions := options.Find()
  //findOptions.SetLimit(10) // arbitrarily set to 10
  //storing decoded docs here
  var results[]*Trainer
  // bson.D as the filter will match all documents
  cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
  if err != nil {
    log.Fatal(err)
  }
  //now need to iterate through cursor to decode docs one at a time
  for cur.Next(context.TODO()) {

    var elem Trainer // value to decode single doc into

    err := cur.Decode(&elem)
    if err != nil {
      log.Println(err)
    }

    results = append(results, &elem)

  }
  cur.Close(context.TODO())

  json.NewEncoder(w).Encode(results)
}

// Get specific trainers
func getTrainer(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r) //Get the params
  filter := bson.D{{"name", params["name"]}}

  var result Trainer

  err = collection.FindOne(context.TODO(), filter).Decode(&result)
  if err != nil {
    log.Println(err)
  }

  json.NewEncoder(w).Encode(&result)

}

// Create trainer
func createTrainer(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  var trainer Trainer
  _ = json.NewDecoder(r.Body).Decode(&trainer)
  filter := bson.D{{"name", trainer.Name}}
  var result Trainer

  findErr := collection.FindOne(context.TODO(), filter).Decode(&result)

  if findErr != nil {
    if findErr == mongo.ErrNoDocuments {
      insertResult, err := collection.InsertOne(context.TODO(), trainer)
      if err != nil {
        log.Println(err)
      }
      log.Println("Inserted a single document: ", insertResult.InsertedID)
      json.NewEncoder(w).Encode(trainer)
      return
    }
    log.Println(findErr)
  }
  log.Println("Trainer name already taken, send PUT to update instead.")

}

//Delete trainer
func deleteTrainer(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r) //Get the params
  filter := bson.D{{"name", params["name"]}}

  result, err := collection.DeleteOne(context.TODO(), filter)
  if err != nil {
    log.Println(err)
    return
  }
  if result.DeletedCount == 0 {
    log.Println("Trainer not found in collection.")
    return
  }

  log.Printf("Deleted document with name: %v from the trainers collection.\n", params["name"])

}

func main() {
  //|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| CONNECTION TO MONGODB

    //checking the connection
  err = client.Ping(context.TODO(), nil)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Connected to MongoDB...")
  fmt.Println("Dropping collection...")
  collection.Drop(context.TODO())
  fmt.Println("Ready for requests...")
  //|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| --CONNECTION
  //routing
  // new router
  // r is our mux variable
  r := mux.NewRouter()

  // Route handles & endpoints
  r.HandleFunc("/trainers", getTrainers).Methods("GET")
  r.HandleFunc("/trainers/{name}", getTrainer).Methods("GET")
  r.HandleFunc("/trainers", createTrainer).Methods("POST")
  //r.HandleFunc("/trainers/{name}", updateTrainer).Methods("PUT")
  r.HandleFunc("/trainers/{name}", deleteTrainer).Methods("DELETE")

  //Starting the server
  log.Fatal(http.ListenAndServe(":3000", r))


/* ||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| BEGIN MONGODB WORK

// Set client options
  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

  // Connect to MongoDB
  client, err := mongo.Connect(context.TODO(), clientOptions)

  if err != nil {
    log.Fatal(err)
  }

  //checking the connection
  err = client.Ping(context.TODO(), nil)

  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Connected to MongoDB! owo")

  //can now get a handle for the trainers collection
  // will later use this handle to query the collection
  collection := client.Database("test").Collection("trainers")

  //|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| CREATE
  //create structs to insert into the DB
  //NAME(string), AGE(int), CITY(string)
  jeff := Trainer{"Jeff Winger", 40, "Greendale"}
  abed := Trainer{"Abed Nadir", 21, "Greendale"}
  troy := Trainer{"Troy Barnes", 21, "Greendale"}

  //inserting a single document
  insertResult, err := collection.InsertOne(context.TODO(), jeff)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Inserted a single document: ", insertResult.InsertedID)

  //inserting multiple, InsertMany() will take a slice of objects
  trainers := []interface{}{abed, troy}

  insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

  //|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| UPDATE

  // updating a single document. Filter document used to match doc,
  // and update doc describes the update operation. Here, using bson.D types

  jeffFilter := bson.D{{"name", "Jeff Winger"}}

  update := bson.D{
    {"$inc", bson.D{
      {"age", 1},
      }},
  }

  updateResult, err := collection.UpdateOne(context.TODO(), jeffFilter, update)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("Matched %v documents and updated %v documents.\n",
              updateResult.MatchedCount,
              updateResult.ModifiedCount)


  //|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| READ
  // to find one doc, need filter doc
  // and pointer to value in which the result can be decoded.

  // value into which result will be decoded
  var result Trainer

  err = collection.FindOne(context.TODO(), jeffFilter).Decode(&result)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("Found a single document %+v\n", result)

  //|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| DELETE

  //collection.Drop() to delete an entire collection

  //DeleteMany or DeleteOne
  deleteResult, err := collection.DeleteMany(context.TODO(), jeffFilter)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("Deleted %v documents in the trainers collection.\n", deleteResult)

  //closing the connection
  err = client.Disconnect(context.TODO())

  if err != nil {
      log.Fatal(err)
  }

  collection.Drop(context.TODO())
  fmt.Println("Connection to MongoDB closed. uwu")
  |||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| */
} // MAIN





















//space
