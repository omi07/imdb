package app

import (
	"context"
	"fmt"
	"log"

	"../model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo connects to mongo
func ConnectMongo(mctx context.Context, mongohost string) (*mongo.Client, error) {
	var client *mongo.Client
	var err error
	for i := 0; i <= MAXRETRY; i++ {
		client, err = mongo.Connect(mctx, options.Client().ApplyURI(mongohost), options.Client().SetMaxPoolSize(MAXCONNPOOL))
		if err == nil {
			break
		}
		log.Printf("Mongo Connect Retry Attempt :%v", i)
		log.Println(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, err
}

//Get collection of mongodb
func GetDBCollection(client *mongo.Client, dbname string, collectionname string) *mongo.Collection {
	mctx := context.TODO()
	err := client.Ping(mctx, nil)
	if err != nil {
		client, mgerr = ConnectMongo(mctx, MONGOHOST)
		if mgerr != nil {
			log.Fatalf("Mongo Connection Failed")
			return nil
		}
	}
	collection := client.Database(dbname).Collection(collectionname)
	return collection
}

func MongoFind(client *mongo.Client, user model.User, dbname string, collectionname string) (*mongo.Collection, error) {
	var result model.User
	collection := GetDBCollection(client, dbname, collectionname)

	err := collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	return collection, err

}

//InsertingToMongo inserts data into mongo
func MongoInsert(collection *mongo.Collection, data interface{}) error {
	mctx := context.TODO()
	_, err := collection.InsertOne(mctx, data)
	return err

}

func GetUserid(collection *mongo.Collection) (int64, error) {
	mctx := context.TODO()

	filter := bson.D{}
	projection := bson.D{{"uid", 1}, {"_id", 0}}
	cur, err := collection.Find(mctx, filter, options.Find().SetProjection(projection), options.Find().SetSort(bson.D{{"_id", -1}}), options.Find().SetLimit(1))

	defer cur.Close(mctx)
	var result map[string]int64
	for cur.Next(mctx) {
		err = cur.Decode(&result)
	}
	return result["uid"] + 1, err
}

func FindAndUpdate(collection *mongo.Collection, filter interface{}, update interface{}) bool {
	mctx := context.TODO()
	singleResult := collection.FindOneAndUpdate(mctx, filter, update)
	if singleResult.Err() != nil {
		log.Printf(" Find error: %v", singleResult.Err())
		return false
	}
	return true
}

func MongoFindOne(collection *mongo.Collection, filter interface{}, projection interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	mctx := context.TODO()

	mgerr := collection.FindOne(mctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)

	if mgerr != nil {
		fmt.Printf("ERROR %v  ", mgerr.Error())
		return nil, mgerr
	}
	return result, nil
}

func CheckRatingExists(collection *mongo.Collection, movieid string, userid int64) bool {

	filter := bson.D{{"ratedmovies.movieid", movieid}, {"uid", userid}}
	projection := bson.D{{"rating.rating", 1}, {"uid", 1}}

	_, err := MongoFindOne(collection, filter, projection)
	if err != nil && err.Error() == "mongo: no documents in result" {
		return false
	}
	return true

}
