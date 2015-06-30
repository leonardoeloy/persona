package main

import (
	"database/sql"
	"log"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func CheckDatabaseConnectivity(dataSourceName string, maxIdleConns int) {
	log.Println("Checking for database connectivity...")
	var dataSource string = "(blank)"

	// Redacts password in datasource name
	reg, err := regexp.Compile("\\:.*@")
	if err != nil {
		log.Fatalf("Couldn't compile regex string to redact passwords on datasource name [%s]. Falling back to blank value.", err)
	} else {
		dataSource = reg.ReplaceAllString(dataSourceName, ":(redacted)@")
	}

	log.Printf("Opening database connection using [%s]", dataSource)
	TryToOpenDatabase()

	db.SetMaxIdleConns(maxIdleConns)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database [%s]", err)
	}

	log.Println("Please check if there were errors. Even though connection might be unsuccessful, the service will start.")
}

func TryToOpenDatabase() {
	var err error
	if db != nil {
		return
	}

	db, err = sql.Open("mysql", *dataSourceName)
	if err != nil {
		log.Fatalf("Error initializing database connection [%s]", err.Error())
	}
}

func Query(sql string, f func(rows *sql.Rows)) {
	TryToOpenDatabase()

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatalf("Error querying database [%s]", err)
	}

	defer rows.Close()
	for rows.Next() {
		f(rows)
	}

	err = rows.Err()
	if err != nil {
		log.Fatalf("Error after querying database [%s]", err)
	}
}