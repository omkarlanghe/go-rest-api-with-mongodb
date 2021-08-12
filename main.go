package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type Student struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Age  string `json:"age,omitempty" bson:"age,omitempty"`
	Sex  string `json:"sex,omitempty" bson:"sex,omitempty"`
	City string `json:"city,omitempty" bson:"city,omitempty"`
}

var client *mongo.Client

// Method to get list of all students
func getAllStudentsEndpoint(response http.ResponseWriter, request *http.Request) {
	// Add headers in response interface
	response.Header().Add("content-type", "application/json")

	// create an object of type struct student
	var students []Student

	// getting collection
	collection := client.Database("student-records").Collection("students")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// quering on a collection which will return cursor
	cursor, err := collection.Find(ctx, bson.M{})

	// error handling
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	defer cursor.Close(ctx)

	// iterating over a cursor
	for cursor.Next(ctx) {
		var student Student
		cursor.Decode(&student)
		fmt.Printf("%+v\n", student)
		students = append(students, student)
	}

	// error handling
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(students)
}

// Method to insert a student in database
func insertStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	// Add headers in response interface
	response.Header().Add("content-type", "application/json")

	// creating an object of type struct student
	var student Student

	// error handling
	err := json.NewDecoder(request.Body).Decode(&student)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	// reading collection
	collection := client.Database("student-records").Collection("students")

	// inseting a record in a collection object
	ctx, err := collection.InsertOne(context.TODO(), student)

	// error handling
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	// printing a message
	fmt.Println("Inserted a single document: ", ctx)

	json.NewEncoder(response).Encode(ctx.InsertedID)
}

// main method
func main() {
	fmt.Println("Starting the Go Application Server running on port 8000...")

	// connecting to mongodb
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	// register route and endpoint
	router := mux.NewRouter()
	router.HandleFunc("/students", getAllStudentsEndpoint).Methods("GET")
	router.HandleFunc("/students", insertStudentEndpoint).Methods("POST")

	http.ListenAndServe(":8000", router)
}
