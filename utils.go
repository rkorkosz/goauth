package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

var (
	templates = make(map[string]*template.Template)
	store     = sessions.NewCookieStore([]byte("Sdt43F0AJER[90sdfpSDF[FJAOa'sdpfo"))
)

const (
	DB_NAME      = "goauth"
	SESSION_NAME = "session"
)

func executeTemplate(w http.ResponseWriter, name string, data interface{}) {
	var tmpl *template.Template
	var ok bool
	if tmpl, ok = templates[name]; !ok {
		tmpl = template.Must(template.New(name).ParseFiles(
			"templates/base.html",
			fmt.Sprintf("templates/%s.html", name),
		))
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func getCollection(name string) *mgo.Collection {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	return session.Clone().DB(DB_NAME).C(name)
}

func getSession(r *http.Request) *sessions.Session {
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		panic(err)
	}
	return session
}

func loginRequired(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := getSession(r)
		if session.Values["user"] == nil {
			http.Redirect(w, r, "/auth/login/", http.StatusFound)
			return
		}
		c := getCollection("users")
		var user User
		err := c.FindId(bson.ObjectIdHex(session.Values["user"].(string))).One(&user)
		if err != nil {
			http.Redirect(w, r, "/auth/ogin/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
