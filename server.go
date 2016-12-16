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

func createResponse(data interface{}, w *http.ResponseWriter) {
	response := bson.M{"data": data, "success": true}
	json.NewEncoder(*w).Encode(&response)
	return
}

func createFailureResponse(reason string, w *http.ResponseWriter) {
	response := bson.M{"reason": reason, "success": false}
	json.NewEncoder(*w).Encode(&response)
}

func bsonify(r *http.Request) bson.M {
	body := bson.M{}
	check(json.NewDecoder(r.Body).Decode(&body))
	return body
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	c := session.DB("freedom").C("user")
	request := bsonify(r)
	doc := models.User{}
	readSuccess := dao.ReadOne(c, request, &doc)
	if readSuccess {
		if doc.Username == request["username"] && doc.Password == request["password"] {
			createResponse(doc, &w)
			return
		}
		createFailureResponse("Could not find a matching username and password combination", &w)
		return
	}
	createFailureResponse("Something went wrong", &w)
	return
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	c := session.DB("freedom").C("user")
	request := bsonify(r)
	writeSuccess := dao.Create(c, request)
	w.Header().Set("Content-Type", "application/json")
	if writeSuccess {
		createResponse(request, &w)
		return
	}
	createFailureResponse("Something went wrong creating your user account", &w)
	return
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
