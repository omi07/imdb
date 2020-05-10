package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../app"
	"../model"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

//Register  User in System
func Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var result model.User
	collection := app.GetDBCollection(app.Mongoconn, "imdb", "userdetails")

	//Checks whether username exists
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				res.Error = "Error While Hashing Password, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			user.Password = string(hash)
			//Gets max id from documents
			user.Uid, _ = app.GetUserid(collection)

			fmt.Printf("USERID :%v \n", user.Uid)
			//Inserting a new userdetails in Mongo
			err = app.MongoInsert(collection, user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again " + err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}
			res.Result = "Registration Successful"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "Username already Exists!!"
	json.NewEncoder(w).Encode(res)
	return
}

//Login Handler

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	collection := app.GetDBCollection(app.Mongoconn, "imdb", "userdetails")
	var result model.User
	var res model.ResponseResult

	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}

	//Password check
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	//

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  result.Username,
		"firstname": result.FirstName,
		"lastname":  result.LastName,
		"role":      result.Role,
		"uid":       result.Uid,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	json.NewEncoder(w).Encode(result)

}

//AddMovies insert movie data into DB
func AddMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	role, err := app.GetUserRole(r.Header.Get("Authorization"))
	var res model.ResponseResult
	if role != "admin" || err != nil {
		res.Error = "Only Admin can add Movies " + err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var movie model.Movie
	body, _ := ioutil.ReadAll(r.Body)
	jsonerr := json.Unmarshal(body, &movie)
	if jsonerr != nil {
		log.Fatal(err)
	}

	collection := app.GetDBCollection(app.Mongoconn, "imdb", "movies")

	//Inserting a new userdetails in Mongo
	err = app.MongoInsert(collection, movie)
	if err != nil {
		res.Error = "Error While Inserting Movie " + err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Result = "Movie added Successfully"
	json.NewEncoder(w).Encode(res)
	return

}

func RateCommentMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	userid, err := app.VerifyUser(r.Header.Get("Authorization"))

	var res model.ResponseResult
	if userid <= 0 || err != nil {
		res.Error = "Only Registered User can rate or add comment to movies" + err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var data map[string]interface{}
	jsonerr := json.Unmarshal(body, &data)
	if jsonerr != nil {
		log.Fatal(err)
	}
	collection := app.GetDBCollection(app.Mongoconn, app.Dbname, app.UserCollection)
	//Checks for rating to calculate average rating in movies collection
	if _, ok := data["rating"]; ok {
		//Check if rating for movie exists
		exists := app.CheckRatingExists(collection, data["movieid"].(string), userid)
		if exists == true {
			res.Error = "Movie Rating already Exists"
			json.NewEncoder(w).Encode(res)
			return
		}

		mvcollection := app.GetDBCollection(app.Mongoconn, "imdb", "movies")
		mctx := context.TODO()
		mvfilter := bson.D{{"movieid", data["movieid"]}}

		mvdata, mgerr := app.MongoFindOne(mvcollection, mvfilter, nil)
		if mgerr != nil {
			res.Error = "Error Finding Data for movieid " + data["movieid"].(string) + " " + mgerr.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		updaterating := (mvdata["rating"].(float64) + data["rating"].(float64)) / (float64(mvdata["totalusers"].(int32)) + 1)

		mvupdate := bson.D{{"$set", bson.D{{"rating", updaterating}, {"totalusers", mvdata["totalusers"].(int32) + 1}}}}

		aggerr := mvcollection.FindOneAndUpdate(mctx, mvfilter, mvupdate)
		if aggerr.Err() != nil {
			//log.Printf("Error Updating rating for the movie %v", aggerr.Err())
			res.Error = "Error Updating rating for the movie " + aggerr.Err().Error()
			json.NewEncoder(w).Encode(res)
			return
		}
	}
	//Updates rating and comments across the user
	//collection := app.GetDBCollection(app.Mongoconn, "imdb", "userdetails")

	filter := bson.D{{"uid", userid}}

	update := bson.M{"$push": bson.M{"ratedmovies": data}}
	result := app.FindAndUpdate(collection, filter, update)
	if result == false {
		res.Error = "Updating rating & Comments Failed"
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Result = "Updated Successfully"
	json.NewEncoder(w).Encode(res)
	return

}

//SearchMovie searches the movie by movieid
func SearchMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	movieid := r.URL.Query().Get("movieid")
	var res model.ResponseResult
	mvcollection := app.GetDBCollection(app.Mongoconn, "imdb", "movies")
	filter := bson.M{"movieid": movieid}
	projection := bson.M{"_id": 0}
	result, err := app.MongoFindOne(mvcollection, filter, projection)
	if err != nil {
		res.Error = "Failed to get Movie " + err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Result = result
	json.NewEncoder(w).Encode(res)
	return
}

//GetMovies Fetchs all movies from collection
func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mctx := context.TODO()
	var res model.ResponseResult

	mvcollection := app.GetDBCollection(app.Mongoconn, "imdb", "movies")

	cur, err := mvcollection.Find(mctx, bson.D{{}}, options.Find().SetProjection(bson.D{{"_id", 0}}))

	if err != nil {
		fmt.Printf("Mongo Find Error %v", err.Error())
	}
	var data []map[string]interface{}
	for cur.Next(mctx) {
		var result map[string]interface{}
		cur.Decode(&result)
		data = append(data, result)
	}
	res.Result = data
	json.NewEncoder(w).Encode(res)
	return

}
