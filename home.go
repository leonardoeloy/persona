package main

import (
	"net/http"
)

func HomeIndex(rw http.ResponseWriter, r *http.Request) {
	rdr.HTML(rw, http.StatusOK, "index", nil)
}
