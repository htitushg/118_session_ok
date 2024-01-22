package assets

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"time"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Chemin     = filepath.Dir(filepath.Dir(b)) + "/"
)
var (
	Ctx context.Context
	Db  *sql.DB
)

const (
	Port = ":8080"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// Create a struct that models the structure of a user in the request body
type CredentialsR struct {
	Pseudo    string `json:"pseudo"` //go.mod, db:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`  //, db:"password"`
	Firstname string `json:"firstname"` //go.mod, db:"firstname"`
	Lastname  string `json:"lastname"`  //go.mod, db:"lastname"`
}
type Credentials struct {
	Pseudo      string `json:"pseudo"`
	Email       string `json:"email"`
	Password    string `json:"password"`  //, db:"password"`
	Password2   string `json:"password2"` //, db:"password2"`
	Firstname   string `json:"firstname"` //go.mod, db:"firstname"`
	Lastname    string `json:"lastname"`  //go.mod, db:"lastname"`
	Address     string `json:"address"`
	Town        string `json:"town"`
	ZipCode     string `json:"zipcode"`
	Country     string `json:"country"`
	Language    string `json:"language"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	Message     string
}

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var Sessions = map[string]Session{}

// each session contains the username of the user and the time at which it expires
type Session struct {
	Pseudo    string
	Expiry    time.Time
	Email     string
	Firstname string
	Lastname  string
	Address   string
	Town      string
	ZipCode   string
	Country   string
}
