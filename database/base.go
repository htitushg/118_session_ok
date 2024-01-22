package database

import (
	"118_session/assets"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var Db *sql.DB

func InitDB() {
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	Db, err = sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/sessiondb")
	if err != nil {
		panic(err)
	}
	CreateTables("users")
	defer Db.Close()
}
func CreateTables(nomtable string) {

	_, err := Db.Exec("create table IF NOT EXISTS users (id int auto_increment primary key, name varchar(255) not null, password varchar(255) not null)")
	assets.CheckError(err)

	creds := &assets.Credentials{
		Pseudo:   "johndoe",
		Password: "mysecurepassword",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	creds.Password = string(hashedPassword)
	assets.CheckError(err)
	query := "INSERT INTO users (name,Password) VALUES (?, ?)"
	_, err = Db.Exec(query, creds.Pseudo, creds.Password)
	if err != nil {
		log.Fatalf("impossible insert user: %s", err)
	}
	creds = &assets.Credentials{
		Pseudo:   "henry",
		Password: "1nhri96p",
	}
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	assets.CheckError(err)
	query = "INSERT INTO users (name,Password) VALUES (?, ?)"
	_, err = Db.Exec(query, creds.Pseudo, hashedPassword)
	if err != nil {
		log.Fatalf("impossible insert user: %s", err)
	}
	creds = &assets.Credentials{
		Pseudo:   "marie",
		Password: "marie123",
	}
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	assets.CheckError(err)
	query = "INSERT INTO users (name,Password) VALUES (?, ?)"
	_, err = Db.Exec(query, creds.Pseudo, hashedPassword)
	if err != nil {
		log.Fatalf("impossible insert user: %s", err)
	}
}
