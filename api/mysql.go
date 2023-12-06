package api

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Users struct {
	UserId       int       `json:"USER_ID"`
	PW           string    `json:"PW"`
	Name         string    `json:"NAME"`
	Email        string    `json:"EMAIL"`
	Tel          string    `json:"TEL"`
	LastAccessDt time.Time `json:"LAST_ACCESS_DT"`
	UpdateDt     time.Time `json:"UPDATE_DT"`
	CreateDt     time.Time `json:"CREATE_DT"`
}

// CREATE TABLE IF NOT EXISTS ghldnjs.USERS (USER_ID varchar(45) NOT NULL, PW binary NOT NULL, NAME varchar(45) NOT NULL, EMAIL varchar(100) NOT NULL, TEL varchar(45) DEFAULT NULL, LAST_ACCESS_DT datetime DEFAULT NULL, UPDATE_DT datetime DEFAULT NULL, CREATE_DT datetime DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (USER_ID) ) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;
func NewTable() error {
	db := NewDb()
	defer db.Close()
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS USERS (USER_ID int NOT NULL, PW binary NOT NULL, NAME varchar(45) NOT NULL, EMAIL varchar(100) NOT NULL, TEL varchar(45) DEFAULT NULL, LAST_ACCESS_DT datetime DEFAULT NULL, UPDATE_DT datetime DEFAULT NULL, CREATE_DT datetime DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (USER_ID) ) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;")
	if err != nil {
		log.Println("create table fail.")
		return errors.New("create table fail")
	}
	return nil
}

func NewDb() *sql.DB {
	db, err := sql.Open("mysql", "root:ROOTroot12#$@tcp(mysql)/ghldnjs?parseTime=true")
	if err != nil {
		log.Println(err)
	}

	return db
}

func CloseDb(db *sql.DB) {
	db.Close()
}

func queryPw(db *sql.DB, name string) (int, string) {
	var pw string
	var id int
	err := db.QueryRow("SELECT USER_ID, PW FROM USERS where name = ?", name).Scan(&id, &pw)
	if err != nil {
		log.Println(err)
	}
	return id, pw
}
