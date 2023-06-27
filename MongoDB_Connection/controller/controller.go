package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shivansh/mongodbconnection/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

const connectionString = "mongodb+srv://shivansh:Shiv$$5955@cluster0.cisosfl.mongodb.net/?retryWrites=true&w=majority"
const dbName = "netflix"
const colName = "watchlist"

// important
var collection *mongo.Collection

// connect with mongoDB
func init() {
	//client options
	clientOptions := options.Client().ApplyURI(connectionString)

	//connect
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance is ready")
}

// MongoDB helpers ->

// 1. insert 1 record
func insertOneMovie(movie *model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted 1 movie in db with id", inserted.InsertedID)
}

// 2. update 1 record
func updateOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

// 3. delete 1 record
func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Movie got deleted with delete count:", deleteCount.DeletedCount)
}

// 4. delete all records
func deleteAllMovie() int64 {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.M{}, nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Movie got deleted with delete count:", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

// 5. get all movies from db
func getAllMovies() []bson.M {
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var movies []bson.M //we can also use []primitive.M

	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			log.Fatal()
		}

		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies
}

// Actual controllers - file

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	//w.Header().Set("Content-Type", "application/x-www-form-urlencode")

	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(&movie)

	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}
