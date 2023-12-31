package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Data Source Name Properties
	dsn := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "sakila",
		AllowNativePasswords: true,
	}

	// Get a database handle
	var err error
	db, err = sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	actors, err := GetActor(201)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Actor found: %v\n", actors)

	// actorID, err := addActor("JOE", "BERRY")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("ID of added actor: %v\n", actorID)
}

type Actor struct {
	actor_id int64
	first_name string
	last_name string
}

func GetActor(actorID int32) ([]Actor, error) {
	var actors []Actor

	result, err := db.Query("SELECT actor_id, first_name, last_name FROM actor WHERE actor_id = ?", actorID)
	if err != nil {
		return nil, fmt.Errorf("GetActor %v: %v", actorID, err)
	}
	defer result.Close()

	// loop through rows
	for result.Next() {
		var acts Actor
		if err := result.Scan(&acts.actor_id, &acts.first_name, &acts.last_name); err != nil {
			return nil, fmt.Errorf("GetActor %v: %v", actorID, err)
		}
		actors = append(actors, acts)

		if err := result.Err(); err != nil {
			return nil, fmt.Errorf("GetActor %v: %v", actorID, err)
		}
	}
	return actors, nil
}

// func addActor(firstname, lastname string) (int64, error) {
// 	result, err := db.Exec("INSERT INTO actor (first_name, last_name) VALUES (?, ?)", firstname, lastname)
// 	if err != nil {
// 		return 0, fmt.Errorf("addActor: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("addActor: %v", err)
// 	}
// 	return id, nil
// }