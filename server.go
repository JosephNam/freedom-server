package main

import (
	"encoding/json"

	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/JosephNam/freedom-server/dao"
	"github.com/JosephNam/freedom-server/models"
)

var session *mgo.Session

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func createResponse(data interface{}) bson.M {
	return bson.M{"data": data, "success": true}
}

func createFailureResponse(reason string) bson.M {
	return bson.M{"reason": reason, "success": false}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	c := session.DB("freedom").C("user")
	request := bson.M{}
	check(json.NewDecoder(r.Body).Decode(&request))
	doc := models.User{}
	readSuccess := dao.ReadOne(c, request, &doc)
	w.Header().Set("Content-Type", "application/json")
	if readSuccess {
		if doc.Username == request["username"] && doc.Password == request["password"] {
			response := createResponse(doc)
			json.NewEncoder(w).Encode(&response)
			return
		}
		failureResponse := createFailureResponse("Could not find a matching username and password combination")
		json.NewEncoder(w).Encode(&failureResponse)
		return
	}
	failureResponse := createFailureResponse("Something went wrong")
	json.NewEncoder(w).Encode(&failureResponse)
	return
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	c := session.DB("freedom").C("user")
	request := bson.M{}
	check(json.NewDecoder(r.Body).Decode(&request))
	writeSuccess := dao.Create(c, request)
	w.Header().Set("Content-Type", "application/json")
	if writeSuccess {
		response := createResponse(request)
		json.NewEncoder(w).Encode(&response)
		return
	}
}

func main() {
	var connErr error
	session, connErr = mgo.Dial("mongodb://dev:dev@ds061345.mlab.com:61345/freedom")
	check(connErr)

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	http.HandleFunc("/api/authenticate", handleLogin)
	http.HandleFunc("/api/register", handleRegister)
	http.ListenAndServe(":5000", nil)
}
