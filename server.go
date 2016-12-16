package main

import (
	"encoding/json"
	"fmt"

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
	result := models.User{}
	success := dao.ReadOne(c, request, result)
	w.Header().Set("Content-Type", "application/json")
	if success == true {
		fmt.Println("User: ", result.Username)
		res := createResponse(result)
		json.NewEncoder(w).Encode(res)
		return
	}
	failureResponse := createFailureResponse("Something went wrong")
	json.NewEncoder(w).Encode(failureResponse)

}

func main() {
	var connErr error
	session, connErr = mgo.Dial("mongodb://dev:dev@ds061345.mlab.com:61345/freedom")
	check(connErr)

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	http.HandleFunc("/api/authenticate", (handleLogin))
	http.ListenAndServe(":5000", nil)
}
