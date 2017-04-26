package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/JosephNam/freedom-server/dao"
	"github.com/JosephNam/freedom-server/models"

	"github.com/julienschmidt/httprouter"

	"github.com/stripe/stripe-go"
	//	"github.com/stripe/stripe-go/charge"
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

func createCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	c := session.DB("freedom").C("user")
	request := bsonify(r)
	doc := models.User{}
	writeSuccess := dao.Create(c, request)
	if writeSuccess {
		createResponse(request, &w)
		return
	}
	createFailureResponse("Something went wrong connecting your card", &w)
	return
}

func handleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func handleRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	c := session.DB("freedom").C("user")
	request := bsonify(r)
	writeSuccess := dao.Create(c, request)
	if writeSuccess {
		createResponse(request, &w)
		return
	}
	createFailureResponse("Something went wrong creating your user account", &w)
	return
}

// FreedomCharge ...
/*
type FreedomCharge struct {
	Amount   float64
	Name     string
	ExpMonth string
	ExpYear  string
	Number   string
	CVC      string
}
*/

func handleCharge(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// c := session.DB("freedom").C("tithe")
	request := bsonify(r)
	fmt.Println(request)
	createResponse(bson.M{"ch_id": "testing"}, &w)
	return
	/*
		c := FreedomCharge{
			Amount:   request["amount"].(float64),
			Name:     request["name"].(string),
			ExpMonth: request["expMonth"].(string),
			ExpYear:  request["expYear"].(string),
			Number:   request["number"].(string),
			CVC:      request["cvc"].(string),
		}
		chargeParams := &stripe.ChargeParams{
			Amount:   uint64(c.Amount),
			Currency: "usd",
			Desc:     "Charge for " + c.Name,
		}
		chargeParams.SetSource(&stripe.CardParams{
			Name:   c.Name,
			Number: c.Number,
			Month:  c.ExpMonth,
			Year:   c.ExpYear,
			CVC:    c.CVC,
		})
		chargeParams.AddMeta("key", "value")

		charge.New(chargeParams)
			check(err)
			fmt.Println(ch.ID)
	*/
}

func main() {
	var connErr error
	session, connErr = mgo.Dial("mongodb://dev:dev@ds061345.mlab.com:61345/freedom")
	check(connErr)

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	router := httprouter.New()

	stripe.Key = "sk_test_xW56HPAT1oXJ1YHC9tXIE8sB"

	router.POST("/api/authenticate", handleLogin)
	router.POST("/api/register", handleRegister)
	router.POST("/api/charge", handleCharge)

	http.ListenAndServe(":5000", router)
}
