package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/mail"
)

func AuthRouter(r *mux.Router) {
	r.HandleFunc("/register/", RegisterForm).Methods("GET")
	r.HandleFunc("/register/", Register).Methods("POST")
	r.HandleFunc("/logout/", Logout).Methods("GET")
	r.HandleFunc("/login/", LoginForm).Methods("GET")
	r.HandleFunc("/login/", Login).Methods("POST")
}

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
}

func RegisterForm(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "register", map[string]interface{}{csrf.TemplateTag: csrf.TemplateField(r)})
}

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	errors := make(map[string]interface{})
	if r.FormValue("password") != r.FormValue("confirm_password") {
		errors["Password"] = "Passwords does not match"
	}
	address, err := mail.ParseAddress(r.FormValue("email"))
	if err != nil {
		errors["Email"] = "Email is invalid"
	}
	if len(errors) > 0 {
		errors[csrf.TemplateTag] = csrf.TemplateField(r)
		executeTemplate(w, "register", errors)
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user := User{
		ID:       bson.NewObjectId(),
		Email:    address.Address,
		Password: string(pass[:]),
	}
	c := getCollection("users")
	err = c.Insert(&user)
	if err != nil {
		panic(err)
	}
	session := getSession(r)
	session.Values["user"] = user.ID.Hex()
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	session.Values["user"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/auth/login/", http.StatusFound)
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "login", map[string]interface{}{csrf.TemplateTag: csrf.TemplateField(r)})
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	errors := make(map[string]interface{})
	var user User
	c := getCollection("users")
	err := c.Find(bson.M{"email": r.FormValue("email")}).One(&user)
	if err != nil {
		errors["email"] = "User with this email not found"
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.FormValue("password")))
	if err != nil {
		errors["password"] = "Wrong password"
	}
	if len(errors) > 0 {
		errors[csrf.TemplateTag] = csrf.TemplateField(r)
		executeTemplate(w, "login", errors)
		return
	}
	session := getSession(r)
	session.Values["user"] = user.ID.Hex()
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
