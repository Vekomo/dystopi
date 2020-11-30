package main

import (
    "context"
    "fmt"
    "log"

    //"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

//Trainer type
type Trainer struct {
    Name string
    Age  int
    City string
}

func main() {
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







  //closing the connection
  err = client.Disconnect(context.TODO())

  if err != nil {
      log.Fatal(err)
  }
  fmt.Println("Connection to MongoDB closed. uwu")

} // MAIN




















//space
