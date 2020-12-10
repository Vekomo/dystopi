package main

import (
    "context"
    "strconv"
    "fmt"
    "os"
    "math"
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
  RatingGiven int
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
  judgeFilter := bson.D{{"username", fields.Judge}}
  targetFilter := bson.D{{"username", fields.Target}}
  //find and get doc for both judge and target
  var judgeDoc User
  var targetDoc User
  //finding judge document
  err = collection.FindOne(context.TODO(), judgeFilter).Decode(&judgeDoc)
  if err != nil {
    log.Println(err)
    log.Println("Could not find judge user.")
    return
  }
  //finding target user document
  err = collection.FindOne(context.TODO(), targetFilter).Decode(&targetDoc)
  if err != nil {
    log.Println(err)
    log.Println("Could not find target user.")
    return
  }

  // **. First make sure target is not already in Judgers map
  // 1. Add targets user name and rating given to judges map
  // 2. Use judges influence to add to targets score map
  // 3. Re-calculate targets rating
  // 4. Re-calculate targets influence
  // 5. Increment targets rated by count


  //Adding target to judges judgements map
  _, ok := judgeDoc.Judgements[fields.Target]
  if ok == true {
    log.Println("Judge has already rated: ", fields.Target)
    return
  }
  judgeDoc.Judgements[fields.Target] = fields.RatingGiven
  trackUpdate := bson.D {
    {"$set", bson.D{
      {"Judgements", judgeDoc.Judgements},
      }},
  }
  _, err := collection.UpdateOne(context.TODO(), judgeFilter, trackUpdate)
  if err != nil {
    log.Println(err)
    log.Println("Could not match/update judges judgements map.")
    return
  }
  // Using judges influence update targets score map

  targetDoc.Score[strconv.Itoa(fields.RatingGiven)] += judgeDoc.Influence
  scoreUpdate := bson.D {
    {"$set", bson.D{
      {"Score", targetDoc.Score},
      }},
  }
  _, err = collection.UpdateOne(context.TODO(), targetFilter, scoreUpdate)
  if err != nil {
    log.Println(err)
    log.Println("Could not match/update targets score map.")
    return
  }

  //Calculating new influence/rating/and new rated by count
  ratingNume := 0
  totalWeight := 0
  for key, weight := range targetDoc.Score {
    scoreInt, _ := strconv.Atoi(key)
    ratingNume += scoreInt * weight
    totalWeight += weight
  }
  newRating := ratingNume/totalWeight
  newRatedBy := targetDoc.RatedBy + 1    //Find out which form is best
  newInfluence := int(math.Pow(float64(newRating), 2))
  newInfluence = int(newInfluence * (newRatedBy/150))
  // Threshold to get through until over a 1, get rated by a lot of people!
  if newInfluence < 1 {
    newInfluence = 1
  }
  // Updating rating
  ratingUpdate := bson.D {
    {"$set", bson.D{
      {"Rating", newRating},
      }},
  }
  _, err = collection.UpdateOne(context.TODO(), targetFilter, ratingUpdate)
  if err != nil {
    log.Println(err)
    log.Println("Could not update rating.")
    return
  }
  // Updating rated by count
  ratedByUpdate := bson.D {
    {"$set", bson.D{
      {"RatedBy", newRatedBy},
      }},
  }
  _, err = collection.UpdateOne(context.TODO(), targetFilter, ratedByUpdate)
  if err != nil {
    log.Println(err)
    log.Println("Could not update target rating count, ratedBy.")
    return
  }
  // Updating influence
  influenceUpdate := bson.D {
    {"$set", bson.D{
      {"Influence", newInfluence},
      }},
  }
  _, err = collection.UpdateOne(context.TODO(), targetFilter, influenceUpdate)
  if err != nil {
    log.Println(err)
    log.Println("Could not update target influence.")
    return
  }
  log.Println("Completed updating target/judger documents...")
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
