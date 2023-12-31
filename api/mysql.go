package api

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"ynp/env"

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
	_, err := env.MyDB.Exec("CREATE TABLE IF NOT EXISTS USERS (USER_ID int NOT NULL, PW binary(100) NOT NULL, NAME varchar(45) NOT NULL, EMAIL varchar(100) NOT NULL, TEL varchar(45) DEFAULT NULL, LAST_ACCESS_DT datetime DEFAULT NULL, UPDATE_DT datetime DEFAULT NULL, CREATE_DT datetime DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (USER_ID) ) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;")
	if err != nil {
		return errors.New("create table fail")
	}
	return nil
}

func NewDb() *sql.DB {
	db, err := sql.Open("mysql", "root:ROOTroot12#$@tcp(127.0.0.1:3306)/ghldnjs?parseTime=true")
	if err != nil {
		log.Println(err)
	}
	return db
}

func CloseDb(db *sql.DB) {
	db.Close()
}

func queryPw(db *sql.DB, name, password string) (int, error) {
	var id int

	err := db.QueryRow("SELECT USER_ID FROM USERS WHERE NAME = ? AND PW = ?;", name, password).Scan(&id)
	return id, err
}

func queryId(db *sql.DB, name string) (int, error) {
	var id int

	err := db.QueryRow("SELECT USER_ID FROM USERS WHERE NAME = ?;", name).Scan(&id)
	return id, err
}

func InsertNums(db *sql.DB, round, nums string) error {
	_, err := db.Exec("INSERT INTO lotto(round, nums) VALUES(?, ?);", round, nums)
	return err
}

func getNums(db *sql.DB, round string) (string, error) {
	var nums string

	err := db.QueryRow("SELECT nums FROM lotto WHERE round = ?;", round).Scan(&nums)
	return nums, err
}

func getLatestNums(db *sql.DB) (string, error) {
	var nums string

	err := db.QueryRow("SELECT nums FROM lotto ORDER BY CAST(round AS DECIMAL) DESC LIMIT 1;").Scan(&nums)
	return nums, err
}
