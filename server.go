package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"text/template"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

type userDetails struct {
	Username string
	Password string
	Age      int
}

const (
	DBUser     = "postgres"
	DBPassword = "postgres"
	DBName     = "infra101"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func prettyPrintMyResults(allMyRows *sql.Rows) (results []userDetails) {
	var allData []userDetails
	for allMyRows.Next() {
		var username string
		var password string
		var age int

		err := allMyRows.Scan(&username, &password, &age)
		checkErr(err)

		myDataContainer := userDetails{
			Username: username,
			Password: password,
			Age:      age,
		}
		allData = append(allData, myDataContainer)
	}

	return allData
}

func showMyTable(db *sql.DB, tableName string) (results []userDetails) {
	fmt.Printf("Querying the table %s for data\n", tableName)
	mySelectQuery := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(mySelectQuery)
	checkErr(err)

	// just checking out the datatype
	fmt.Println(reflect.TypeOf(rows))

	return prettyPrintMyResults(rows)
}

func legacyServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":1000", nil)
}

func main() {
	r := mux.NewRouter()

	// simple data collection API
	r.HandleFunc("/infra/{person}/read/{item}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		personName := vars["person"]
		itemName := vars["item"]

		fmt.Fprintf(w, "Request received: %s item for %s", itemName, personName)
	})

	// templating data API
	r.HandleFunc("/infra/all", func(w http.ResponseWriter, r *http.Request) {
		containerFile := template.Must(template.ParseFiles("templates/container.html"))

		data := TodoPageData{
			PageTitle: "Satoshi's List",
			Todos: []Todo{
				{Title: "Bitcoin", Done: true},
				{Title: "Ethereum", Done: true},
				{Title: "HyperLedger", Done: true},
				{Title: "Corda", Done: false},
			},
		}

		containerFile.Execute(w, data)
	})

	// postgres initialization
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DBUser, DBPassword, DBName)
	db, err := sql.Open("postgres", dbInfo)
	checkErr(err)
	fmt.Println(reflect.TypeOf(db))

	// show all data API
	r.HandleFunc("/infra/db/view/{dbname}/all", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var dbResults = showMyTable(db, vars["dbname"])
		displayerFile := template.Must(template.ParseFiles("templates/displayer.html"))

		displayerFile.Execute(w, dbResults)
	})

	// input data -> store into postgres -> show API
	r.HandleFunc("/infra/db/insert/{dbname}/{username}/{password}/{age}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var insertedUsername string
		insertPayload := fmt.Sprintf("INSERT INTO %s(username, password, age) VALUES('%s', '%s', %s) returning username;", vars["dbname"], vars["username"], vars["password"], vars["age"])
		fmt.Printf("My query is: %s\n", insertPayload)
		err = db.QueryRow(insertPayload).Scan(&insertedUsername)
		checkErr(err)

		var dbResults = showMyTable(db, vars["dbname"])
		displayerFile := template.Must(template.ParseFiles("templates/displayer.html"))

		displayerFile.Execute(w, dbResults)
		fmt.Fprintf(w, "worked!") // don't ignore this
	})

	http.ListenAndServe(":2000", r)
}
