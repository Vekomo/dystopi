package main

import (
    "context"
    "fmt"
    "os"
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

//User type for collection
type User struct {
    Username   string
    Password   string
    Rating     float64
    Influence  int
    Judgements map[string]int
    Score      map[string]int
    RatedBy    int

}
// Fields struct for when a judgement occurs
type JudgementFields struct {
  Judge       string
  Target      string
  RatingGiven float64
}

// Set client options
var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
// Connect to MongoDB
var client, err = mongo.Connect(context.TODO(), clientOptions)
//can now get a handle for the Users collection
// will later use this handle to query the collection
var collection = client.Database("dystopi").Collection("users")


//|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  findOptions := options.Find()
  //findOptions.SetLimit(10) // arbitrarily set to 10
  //storing decoded docs here
  var results[]*User
  // bson.D as the filter will match all documents
  cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
  if err != nil {
    log.Fatal(err)
  }
  //now need to iterate through cursor to decode docs one at a time
  for cur.Next(context.TODO()) {

    var elem User // value to decode single doc into

    err := cur.Decode(&elem)
    if err != nil {
      log.Println(err)
    }

    results = append(results, &elem)

  }
  cur.Close(context.TODO())

  json.NewEncoder(w).Encode(results)
}

//|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| Get single Users
func getUser(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r) //Get the params
  filter := bson.D{{"username", params["username"]}}

  var result User

  err = collection.FindOne(context.TODO(), filter).Decode(&result)
  if err != nil {
    log.Println(err)
  }

  json.NewEncoder(w).Encode(&result)

}

//|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| Create User
func createUser(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  var user User
  _ = json.NewDecoder(r.Body).Decode(&user)
  filter := bson.D{{"username", user.Username}}
  var result User

  findErr := collection.FindOne(context.TODO(), filter).Decode(&result)

  if findErr != nil {
    if findErr == mongo.ErrNoDocuments {
      insertResult, err := collection.InsertOne(context.TODO(), user)
      if err != nil {
        log.Println(err)
      }
      log.Println("Inserted a single document: ", insertResult.InsertedID)
      json.NewEncoder(w).Encode(user)
      return
    }
    log.Println(findErr)
  }
  log.Println("Username already taken, send PUT to update instead.")

}

//|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| Delete User
func deleteUser(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r) //Get the params
  filter := bson.D{{"username", params["username"]}}

  result, err := collection.DeleteOne(context.TODO(), filter)
  if err != nil {
    log.Println(err)
    return
  }
  if result.DeletedCount == 0 {
    log.Println("User not found in collection.")
    return
  }

  log.Printf("Deleted user: %v from the Users collection.\n", params["username"])

}

//|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| Update User
func updateUser(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  var fields JudgementFields
   _ = json.NewDecoder(r.Body).Decode(&fields)
  //judgeFilter := bson.D{{"username", fields.Judge}}
  //targetFilter := bson.D{{"username", fields.Target}}
  ratingGiven := fields.RatingGiven
  log.Println("Judge: " + fields.Judge)
  log.Println("Target: " + fields.Target)
  log.Println("Rating given: " ,ratingGiven)

  return
  //Using judges username, get influence and
  //Add targets username to judges judgements if found. If not exit operation.
  //
  //Apply that influence an

  /** Incrementing influence as a quick test.

  update := bson.D{
    {"$inc", bson.D{
      {"influence", 1},
      }},
  }

  updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
  if err != nil {
    log.Println(err)
    return
  }
  if updateResult.MatchedCount == 0 {
    log.Println("Could not match any documents.")
    return
  }

  log.Printf("Matched %v documents and updated %v documents.\n",
              updateResult.MatchedCount,
              updateResult.ModifiedCount)

  **/

}

func main() {
  //|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| CONNECTION TO MONGODB
    //checking the connection
  err = client.Ping(context.TODO(), nil)
  if err != nil {
    log.Fatal(err)
  }
  //add command line for option to drop collection on run.
  fmt.Println("Connected to MongoDB...")
  // if length requirement not there, no need to check if drop, if it is make sure
  // it is actually drop and not a type.
  if(len(os.Args) > 1 && os.Args[1] == "drop") {
    fmt.Println("Dropping collection...")
    collection.Drop(context.TODO())
  }
  fmt.Println("Ready for requests...")
  //|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||| --CONNECTION
  //routing
  // new router
  // r is our mux variable
  r := mux.NewRouter()

  // Route handles & endpoints
  r.HandleFunc("/users", getUsers).Methods("GET")
  r.HandleFunc("/users/{username}", getUser).Methods("GET")
  r.HandleFunc("/users", createUser).Methods("POST")
  r.HandleFunc("/users", updateUser).Methods("PUT")
  r.HandleFunc("/users/{username}", deleteUser).Methods("DELETE")

  //Starting the server
  log.Fatal(http.ListenAndServe(":3000", r))

} // MAIN





















//space
