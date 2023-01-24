package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// function to check is user exists for signing in
func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user, signedinuser user
	json.NewDecoder(r.Body).Decode(&user)
	rows := DB.Where("BINARY name = ? and password = ?", user.Name, user.Password).Find(&signedinuser).RowsAffected
	if rows != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// list of all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []user
	//select all records from users
	DB.Find(&users)
	json.NewEncoder(w).Encode(&users)

}

// unfollow friend
func DeleteFollowee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var followee follows
	DB.Delete(&followee, "BINARY source_user = ? and target_user = ?", params["username"], params["followeename"])
	json.NewEncoder(w).Encode("deleted followee")

}

// delete tweet
func DeleteTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var tweet tweet
	DB.Delete(&tweet, params["tweetid"])
	json.NewEncoder(w).Encode("deleted tweet")

}

//get all tweets of user

func GetTweetsOfUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var tweets []tweet
	fmt.Println(params["username"])
	DB.Where("BINARY user_name = ?", params["username"]).Find(&tweets)
	json.NewEncoder(w).Encode(&tweets)

}

// function to get all followees of user
func GetFolloweesOfUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//get params passed from request
	params := mux.Vars(r)

	var followers []follows

	//select the records from follow table where source user is current user
	DB.Where("BINARY source_user = ?", params["username"]).Find(&followers)

	json.NewEncoder(w).Encode(&followers)

}

// function to add a new followee (friend)
func AddFollowee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var follow, existing follows
	json.NewDecoder(r.Body).Decode(&follow)

	// check if user exists
	var user user
	rows := DB.Where("BINARY name = ?", follow.SourceUser).Find(&user).RowsAffected
	if rows != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check if the user is already following
	rows = DB.Where("BINARY source_user = ? and target_user = ?", follow.SourceUser, follow.TargetUser).Find(&existing).RowsAffected
	if rows == 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DB.Create(&follow)
	json.NewEncoder(w).Encode(&follow)

}

// function to check if a following relationship exists between 2 users
func CheckFollowing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var existing follows
	params := mux.Vars(r)

	// check if user exists
	var user user
	rows := DB.Where("BINARY name = ?", params["username"]).Find(&user).RowsAffected
	if rows != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//check number of rows returned for the 2 users
	rows = DB.Where("BINARY source_user = ? and target_user = ?", params["username"], params["followeename"]).Find(&existing).RowsAffected
	if rows == 1 {
		w.WriteHeader(http.StatusFound)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// function to create new tweet
func AddTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tweet tweet
	//get the tweet from request
	json.NewDecoder(r.Body).Decode(&tweet)

	//check if user exists
	var user user
	rows := DB.Where("BINARY name = ?", tweet.UserName).Find(&user).RowsAffected
	if rows != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//tweet validation
	if len(tweet.Content) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//craete the tweet and return json
	DB.Create(&tweet)
	json.NewEncoder(w).Encode(&tweet)

}
func AddUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	//create container for the incoming user
	var user user

	//take the data from the request body and put it in the empty container
	json.NewDecoder(r.Body).Decode(&user)

	//username and password length validation
	if len(user.Name) < 3 || len(user.Password) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//create record in table
	err := DB.Create(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//take the data and show it in browser
	json.NewEncoder(w).Encode(&user)
}

// source data name
const dsn = "root:@tcp(127.0.0.1:3306)/demodb?parseTime=true"

var DB *gorm.DB
var err error

// function to setup DB
func setupDB() {

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		panic("cannot connect to DB!!")
	}
	err := DB.AutoMigrate(&user{})
	if err != nil {
		panic("cannot initiate user table")
	}
	err = DB.AutoMigrate(&follows{})
	if err != nil {
		panic("cannot initiate followers table")
	}
	err = DB.AutoMigrate(&tweet{})
	if err != nil {
		panic("cannot initiate tweets table")
	}

	fmt.Println("Connected to DB and setup successfully!")
}
func setupRoutes() {
	r := mux.NewRouter().StrictSlash(true)

	//routes for the apis
	r.HandleFunc("/api/signin", SignIn).Methods("POST")
	r.HandleFunc("/api/user", GetAllUsers).Methods("GET")
	r.HandleFunc("/api/user", AddUser).Methods("POST")
	r.HandleFunc("/api/tweet", AddTweet).Methods("POST")
	r.HandleFunc("/api/user/tweets/{username}", GetTweetsOfUser).Methods("GET")
	r.HandleFunc("/api/user/followees/{username}", GetFolloweesOfUser).Methods("GET")
	r.HandleFunc("/api/follow", AddFollowee).Methods("POST")
	r.HandleFunc("/api/tweet/{tweetid}", DeleteTweet).Methods("DELETE")
	r.HandleFunc("/api/user/followees/{username}/{followeename}", DeleteFollowee).Methods("DELETE")
	r.HandleFunc("/api/user/followees/{username}/{followeename}", CheckFollowing).Methods("GET")

	//allowing CORS for the client
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
	})

	handler := c.Handler(r)
	//start server
	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		log.Fatal("cant start server")
	}

}

// main function
func main() {

	setupDB()
	setupRoutes()

}
