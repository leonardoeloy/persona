package main

import (
	"net/http"
	"time"
	"database/sql"
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"strconv"
)

type Project struct {
	Name string
	StartDate time.Time
	EndDate time.Time
}

type CostCenter struct {
	Id int64
	Name string
}

type PersonAllocation struct {
	Id int64
	Person
	CostCenter
	Allocation int
}

type Person struct {
	Id int64
	Name  string
	Email string
}

func HomeIndex(rw http.ResponseWriter, r *http.Request) {
	rdr.HTML(rw, http.StatusOK, "index", nil)
}

func SimulateGet(rw http.ResponseWriter, r *http.Request) {
	rdr.HTML(rw, http.StatusOK, "simulate", nil)
}

func PeopleGet(rw http.ResponseWriter, r *http.Request) {
	people := make([]Person, 0)

	Query("select p.id, p.name, p.email from person p", func (rows *sql.Rows) {
		person := Person{}

		err := rows.Scan(&person.Id, &person.Name, &person.Email)
		if err != nil {
			log.Fatalf("Error reading from row [%s]", err)
		}

		people = append(people, person)
	})

	rdr.HTML(rw, http.StatusOK, "people", people)
}

func PeopleGetOne(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	person := new(Person)

	err := db.QueryRow("select p.id, p.name, p.email from person p where id =" + vars["id"]).Scan(&person.Id, &person.Name, &person.Email)
	if err != nil {
		log.Fatalf("Error reading from row [%s]", err)
	}

	rdr.HTML(rw, http.StatusOK, "person", person)
}

func PeoplePostOne(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Error parsing form [%s]", err)
	}

	person := new(Person)
	decoder := schema.NewDecoder()
	err = decoder.Decode(person, r.PostForm)
	if err != nil {
		log.Fatalf("Error decoding form [%s]", err)
	}

	stmt, err := db.Prepare("update person set name=?, email=? where id=?")
	if err != nil {
		log.Fatalf("Error preparing query [%s]", err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(person.Name, person.Email, vars["id"])
	if err != nil {
		log.Fatalf("Error updating entity [%s]", err)
	}

	http.Redirect(rw, r, "/people/" + vars["id"], http.StatusMovedPermanently)
}

func PeopleGetNew(rw http.ResponseWriter, r *http.Request) {
	rdr.HTML(rw, http.StatusOK, "person", nil)
}

func PeoplePost(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Error parsing form [%s]", err)
	}

	person := new(Person)
	decoder := schema.NewDecoder()
	err = decoder.Decode(person, r.PostForm)
	if err != nil {
		log.Fatalf("Error decoding form [%s]", err)
	}

	stmt, err := db.Prepare("insert into person (id, name, email) values (null, ?, ?)")
	if err != nil {
		log.Fatalf("Error preparing query [%s]", err)
	}

	defer stmt.Close()
	result, err := stmt.Exec(person.Name, person.Email)
	if err != nil {
		log.Fatalf("Error inserting entity [%s]", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting last inserted id [%s]", err)
	}

	log.Printf("Last id=%d", id)

	http.Redirect(rw, r, "/people/" + strconv.FormatInt(id, 10), http.StatusMovedPermanently)
}

func PeoplePostRemove(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	stmt, err := db.Prepare("delete from person where id = ?")
	if err != nil {
		log.Fatalf("Error preparing query [%s]", err)
		http.StatusText(http.StatusInternalServerError)
		return
	}

	defer stmt.Close()
	_, err = stmt.Exec(vars["id"])
	if err != nil {
		log.Fatalf("Error removing entity [%s]", err)
		http.StatusText(http.StatusInternalServerError)
		return
	}

	http.StatusText(http.StatusNoContent)
}

func PeopleAllocateGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	person := new(Person)
	person.Id , _ = strconv.ParseInt(vars["id"], 10, 64)

	err := db.QueryRow("select p.name from person p where id =" + vars["id"]).Scan(&person.Name)
	if err != nil {
		log.Fatalf("Error reading from row [%s]", err)
	}

	personAllocations := make([]PersonAllocation, 0)
	Query("select pa.id, pa.allocation, pa.cost_center_id, cc.name from person_allocation pa, cost_center cc where cc.id = pa.cost_center_id and pa.persona_id=" + vars["id"],
		func (rows *sql.Rows) {
			costCenter := new(CostCenter)
			personAllocation := new(PersonAllocation)

			err := rows.Scan(&personAllocation.Id, &personAllocation.Allocation, &costCenter.Id, &costCenter.Name)
			if err != nil {
				log.Fatalf("Error reading from row [%s]", err)
			}

			personAllocation.Person = *person
			personAllocation.CostCenter = *costCenter

			personAllocations = append(personAllocations, *personAllocation)
		})

	rdr.HTML(rw, http.StatusOK, "personAllocations", map[string]interface{}{
		"person": person,
		"allocations": personAllocations,
	})
}