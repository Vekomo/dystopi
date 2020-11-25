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
  //collection := client.Database("test").Collection("trainers")



  //closing the connection
  err = client.Disconnect(context.TODO())

  if err != nil {
      log.Fatal(err)
  }
  fmt.Println("Connection to MongoDB closed. uwu")

}




















//space
