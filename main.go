package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DeleteTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var tweet tweet
	DB.Find(&tweet, params["tweetid"])
	DB.Delete(&tweet)
	json.NewEncoder(w).Encode("deleted")

}

func GetTweetsOfUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var tweets []tweet
	DB.Where("user_id = ?", params["userid"]).Find(&tweets)
	json.NewEncoder(w).Encode(&tweets)

}

func GetFolloweesOfUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var followers []follows
	DB.Where("source_id = ?", params["userid"]).Find(&followers)

	json.NewEncoder(w).Encode(&followers)

}

func AddFollowee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var follow follows
	json.NewDecoder(r.Body).Decode(&follow)
	DB.Create(&follow)
	json.NewEncoder(w).Encode(&follow)

}

func AddTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tweet tweet
	json.NewDecoder(r.Body).Decode(&tweet)
	DB.Create(&tweet)
	json.NewEncoder(w).Encode(&tweet)

}
func AddUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	//create container for the incoming user
	var user user

	//take the data from the request body and put it in the empty container
	json.NewDecoder(r.Body).Decode(&user)
	//create record in table
	DB.Create(&user)

	//take the data and show it in browser
	json.NewEncoder(w).Encode(&user)
}

// source data name
const dsn = "root:@tcp(127.0.0.1:3306)/demodb?parseTime=true"

var DB *gorm.DB
var err error

func setupDB() {

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		panic("cannot connect to DB!!")
	}
	DB.AutoMigrate(&user{})
	DB.AutoMigrate(&follows{})
	DB.AutoMigrate(&tweet{})

	fmt.Println("Connected to DB and setup successfully!")
}

func setupRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/api/user", AddUser).Methods("POST")
	r.HandleFunc("/api/tweet", AddTweet).Methods("POST")
	r.HandleFunc("/api/user/tweets/{userid}", GetTweetsOfUser).Methods("GET")
	r.HandleFunc("/api/user/followees/{userid}", GetFolloweesOfUser).Methods("GET")
	r.HandleFunc("/api/follow", AddFollowee).Methods("POST")
	r.HandleFunc("/api/tweet/{tweetid}", DeleteTweet).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))

}

func main() {

	setupDB()
	setupRoutes()

}
