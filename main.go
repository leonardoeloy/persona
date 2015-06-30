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
	lang = flag.String("lang", "en-US", "Change default UI language. JSON translation must be in i18n/<lang>.all.json")

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
	i18nFile := "i18n/" + *lang + ".all.json"
	log.Printf("Loading i18n from [%s]", i18nFile)

	i18n.MustLoadTranslationFile(i18nFile)
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

	log.Println("Persona has left the building.")
}

func SetRoutes(r *mux.Router) {
	simulate := r.Path("/simulate").Subrouter()
	simulate.Methods("GET").HandlerFunc(SimulateGet)

	people := r.Path("/people").Subrouter()
	people.Methods("GET").HandlerFunc(PeopleGet)
	people.Methods("POST").HandlerFunc(PeoplePost)

	r.Path("/people/new").Methods("GET").HandlerFunc(PeopleGetNew)

	person := r.PathPrefix("/people/{id:[0-9]+}").Subrouter()
	person.Methods("GET").HandlerFunc(PeopleGetOne)
	person.Methods("POST").HandlerFunc(PeoplePostOne)
	person.Methods("DELETE").HandlerFunc(PeoplePostRemove)
	person.Path("/allocations").Methods("GET").HandlerFunc(PeopleAllocateGet)

	r.PathPrefix("/s/").Handler(http.StripPrefix("/s/", http.FileServer(http.Dir("./public"))))
	r.HandleFunc("/", HomeIndex)

	http.Handle("/", r)
}
