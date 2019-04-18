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
	username string
	password string
	age      int
}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "infra101"
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
			username: username,
			password: password,
			age:      age,
		}
		allData = append(allData, myDataContainer)
	}

	return allData
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
	r.HandleFunc("/infra/{person}/read/{item}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		personName := vars["person"]
		itemName := vars["item"]

		fmt.Fprintf(w, "Request received: %s item for %s", itemName, personName)
	})

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

	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbInfo)
	checkErr(err)
	fmt.Println(db)

	fmt.Println("Lets Query data")
	rows, err := db.Query("SELECT * FROM demoTable")
	checkErr(err)
	fmt.Println(reflect.TypeOf(rows))

	// don't raw print the rows
	// fmt.Println(rows)

	r.HandleFunc("/infra/db/insert/{dbname}/{username}/{password}/{age}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var insertedUsername string
		insertPayload := fmt.Sprintf("INSERT INTO %s(username, password, age) VALUES('%s', '%s', %s) returning username;", vars["dbname"], vars["username"], vars["password"], vars["age"])
		fmt.Printf("My query is: %s\n", insertPayload)
		err = db.QueryRow(insertPayload).Scan(&insertedUsername)
		checkErr(err)
		fmt.Fprintf(w, "worked!")
	})

	http.ListenAndServe(":2000", r)
}
