package main

import (
	"database/sql"
	"log"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

func CheckDatabaseConnectivity(dataSourceName string, maxIdleConns int) {
	log.Println("Checking for database connectivity...")

	var db *sql.DB
	var err error
	var dataSource string = "(blank)"

	// Redacts password in datasource name
	reg, err := regexp.Compile("\\:.*@")
	if err != nil {
		log.Fatalf("Couldn't compile regex string to redact passwords on datasource name [%s]. Falling back to blank value.", err)
	} else {
		dataSource = reg.ReplaceAllString(dataSourceName, ":(redacted)@")
	}

	log.Printf("Opening database connection using [%s]", dataSource)
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error initializing database connection [%s]", err.Error())
	}

	db.SetMaxIdleConns(maxIdleConns)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database [%s]", err)
	}

	log.Println("Please check if there were errors. Even though connection might be unsuccessful, the service will start.")
	db.Close()
}
