package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
)

type Project struct {
	ID bson.ObjectId `bson:"_id,omitempty" schema:"-"`
	Name string `bson:"name" schema:"name"`
}

func NewController(r *mux.Router, name string, model interface{}) {
	c := Controller{model, name}
	s := r.PathPrefix(fmt.Sprintf("/%s", name)).Subrouter()
	s.HandleFunc("/", c.List).Methods("GET")
	s.HandleFunc("/add/", c.Form).Methods("GET")
	s.HandleFunc("/add/", c.Create).Methods("POST")
	s.HandleFunc("/{id}/", c.Details).Methods("GET")
	s.HandleFunc("/{id}/edit/", c.Form).Methods("GET")
	s.HandleFunc("/{id}/edit/", c.Update).Methods("POST")
	s.HandleFunc("/{id}/delete/", c.Delete).Methods("POST")
}

type Controller struct {
	Model interface{}
	Name  string
}

func (self *Controller) makeOne() interface{} {
	model := reflect.New(reflect.TypeOf(self.Model))
	return model.Interface()
}

func (self *Controller) makeMany() interface{} {
	tp := reflect.TypeOf(self.Model)
	slice := reflect.MakeSlice(reflect.SliceOf(tp), 10, 10)
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	return x.Interface()
}

func (self *Controller) List(w http.ResponseWriter, r *http.Request) {
	objects := self.makeMany()
	c := getCollection(self.Name)
	err := c.Find(nil).All(objects)
	if err != nil {
		panic(err)
	}
	executeTemplate(w, fmt.Sprintf("%s_list", self.Name), objects)
}

func (self *Controller) Details(w http.ResponseWriter, r *http.Request) {
	object := self.makeOne()
	vars := mux.Vars(r)
	c := getCollection(self.Name)
	err := c.FindId(bson.ObjectIdHex(vars["id"])).One(object)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	executeTemplate(w, fmt.Sprintf("%s_details", self.Name), object)
}

func (self *Controller) Form(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{csrf.TemplateTag: csrf.TemplateField(r)}
	vars := mux.Vars(r)
	if oid, ok := vars["id"]; ok {
		object := self.makeOne()
		c := getCollection(self.Name)
		err := c.FindId(bson.ObjectIdHex(oid)).One(object)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		data["Object"] = &object
	}
	executeTemplate(w, fmt.Sprintf("%s_form", self.Name), data)
}

func (self *Controller) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	object := self.makeOne()
	decoder := schema.NewDecoder()
	err := decoder.Decode(object, r.PostForm)
	if err != nil {
		panic(err)
	}
	c := getCollection(self.Name)
	err = c.Insert(object)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/%s/", self.Name), http.StatusFound)
}

func (self *Controller) Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	update := bson.M{}
	decoder := schema.NewDecoder()
	err := decoder.Decode(update, r.PostForm)
	if err != nil {
		panic(err)
	}
	c := getCollection(self.Name)
	err = c.UpdateId(bson.ObjectIdHex(vars["id"]), update)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/%s/", self.Name), http.StatusFound)
}

func (self *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c := getCollection(self.Name)
	err := c.RemoveId(bson.ObjectIdHex(vars["id"]))
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/%s/", self.Name), http.StatusFound)
}
