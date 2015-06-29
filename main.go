package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/unrolled/render.v1"
	"html/template"
)

var (
	version = "1.0"
	bindingAddress = flag.String("binding", ":8180", "HTTP Server Binding Addres, i.e. 127.0.0.1:8080 or simply :8080")
	maxIdleConns = flag.Int("maxIdleConns", 100, "Max. number of idle database connections in the pool")
	dataSourceName = flag.String("dataSourceName", "root:password@/persona", "Data source connection string (only MySQL is supported right now)")
	i18nFile = flag.String("i18n", "i18n/en-US.all.json", "i18n JSON file with translations")

	rdr *render.Render
)

func main() {
	flag.Parse()

	log.Printf("Starting Persona v" + version)

	Initializei18n()
	CheckDatabaseConnectivity(*dataSourceName, *maxIdleConns)
	StartHttpServer()
}

func Initializei18n() {
	log.Printf("Loading i18n from [%s]", *i18nFile)

	i18n.MustLoadTranslationFile(*i18nFile)
}

func StartHttpServer() {
	T, _ := i18n.Tfunc("en-US")
	rdr = render.New(render.Options{
		Directory: "templates",
		Layout: "layout",
		Funcs: []template.FuncMap{{"T": T}},
		Extensions: []string{".html"},
		IsDevelopment: true,
	})

	routes := mux.NewRouter().StrictSlash(false)
	SetRoutes(routes)

	log.Println("Listening on " + *bindingAddress)
	if err := http.ListenAndServe(*bindingAddress, routes); err != nil {
		log.Fatalf("Couldn't bind HTTP server [%s]", err)
	}

	log.Println("Persona has left the building.\n")
}

func SetRoutes(r *mux.Router) {
	personas := r.Path("/personas").Subrouter()
	personas.Methods("GET").HandlerFunc(PersonasIndexHandler)

	r.PathPrefix("/s/").Handler(http.StripPrefix("/s/", http.FileServer(http.Dir("./public"))))
	r.HandleFunc("/", HomeIndex)

	http.Handle("/", r)
}




