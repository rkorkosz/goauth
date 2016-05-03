package main

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func TaskRouter(r *mux.Router) {

}

type Task struct {
	ID      bson.ObjectId
	Subject string
	Project bson.ObjectId
	Author  string
}
