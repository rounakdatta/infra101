package controllers

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"

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

	// validate the presence of username and
	c.Validation.Required(email).Message("Register error: username is a required parameter")
	c.Validation.Required(pwd).Message("Register error: password is a required parameter")

	// username and password must be of minimum length
	c.Validation.MinSize(email, 3).Message("Register error: username length musn't be lesser than 3 chars")
	c.Validation.MinSize(pwd, 6).Message("Register error: password length musn't be lesser than 6 chars")

	emailExpMatcher := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@srmuniv.edu.in")

	if c.Validation.HasErrors() || !emailExpMatcher.MatchString(email) {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}

	return c.Render(email)
}
