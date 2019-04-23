package controllers

import (
	"database/sql"
	"fmt"
	"infra101/app"
	"reflect"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

// App struct definition
type App struct {
	*revel.Controller
}

type userDetails struct {
	Username string
	Password string
	Age      int
}

type registrationDetails struct {
	Uid            string
	Username       string
	SecurePassword string
	CreateDate     string
	AccountActive  bool
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// HashPassword function: hash a password usng bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash function: verify a password hash matches
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

// Index route: home
func (c App) Index() revel.Result {
	ref := c.Params.Get("ref")
	if len(ref) > 0 {
		ref = "You're coming from " + ref
	}

	return c.Render(ref)
}

// Register route: register
func (c App) Register() revel.Result {
	email := c.Params.Get("email")
	pwd := c.Params.Get("pwd")

	// alidation procedures
	c.Validation.Required(email).Message("Register error: username is a required parameter")
	c.Validation.Required(pwd).Message("Register error: password is a required parameter")

	c.Validation.MinSize(email, 3).Message("Register error: username length musn't be lesser than 3 chars")
	c.Validation.MinSize(pwd, 6).Message("Register error: password length musn't be lesser than 6 chars")

	emailExpMatcher := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@srmuniv.edu.in")
	if c.Validation.HasErrors() || !emailExpMatcher.MatchString(email) {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}

	// database write for registration
	timeNow := fmt.Sprintf(time.Now().Format("2006-01-02"))

	pwdHash, hashErr := HashPassword(pwd)
	checkErr(hashErr)

	allDetailsCollected := registrationDetails{
		Uid:            uuid.Must(uuid.NewRandom()).String(),
		Username:       email,
		SecurePassword: pwdHash,
		CreateDate:     timeNow,
		AccountActive:  true,
	}

	_ = allDetailsCollected
	fmt.Println(len(uuid.Must(uuid.NewRandom()).String()))
	fmt.Println(CheckPasswordHash("hello", pwdHash))

	insertPayload := fmt.Sprintf("INSERT INTO users(uid, username, createdate, accountactive, securepassword) VALUES('%s', '%s', '%v', %v, '%s');", allDetailsCollected.Uid, allDetailsCollected.Username, allDetailsCollected.CreateDate, allDetailsCollected.AccountActive, allDetailsCollected.SecurePassword)
	fmt.Printf("My query is: %s\n", insertPayload)
	fmt.Println(app.PQDB)
	app.PQDB.QueryRow(insertPayload)

	return c.Render(email)
}
