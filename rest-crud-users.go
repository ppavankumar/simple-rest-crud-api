package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var users Users

// Users struct which contains
// an array of users
type Users struct {
	Users []User `json:"users"`
}

// User struct which contains a name
// a type and a list of social links
type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Age    int    `json:"Age"`
	Social Social `json:"social"`
}

// Social struct which contains a
// list of links
type Social struct {
	Facebook string `json:"facebook"`
	Twitter  string `json:"twitter"`
}

func readUsersData() {
	jsonFile, err := os.Open("users.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &users)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllUsers")
	json.NewEncoder(w).Encode(users)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["id"])

	for _, User := range users.Users {
		if User.Id == key {
			json.NewEncoder(w).Encode(User)
		}
	}
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new User struct
	// append this to our Users array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var User User
	json.Unmarshal(reqBody, &User)
	// update our global Users array to include
	// our new User
	users.Users = append(users.Users, User)
	// Write the file again with new user array
	file, _ := json.MarshalIndent(users, "", "")
	_ = ioutil.WriteFile("users.json", file, 0644)

	json.NewEncoder(w).Encode(User)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for index, User := range users.Users {
		if User.Id == id {
			users.Users = append(users.Users[:index], users.Users[index+1:]...)
			// Write the file again with new user array
			file, _ := json.MarshalIndent(users, "", "")
			_ = ioutil.WriteFile("users.json", file, 0644)
		}
	}

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", returnAllUsers)
	myRouter.HandleFunc("/user", createNewUser).Methods("POST")
	myRouter.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{id}", returnSingleUser)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	readUsersData()
	handleRequests()
}
