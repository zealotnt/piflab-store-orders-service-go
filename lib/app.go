package lib

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type App struct {
	router *mux.Router
	DB     *DB
	ENV    string
	PORT   string
}

type BodyMultipart struct {
	Buff        bytes.Buffer
	ContentType string
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.router.ServeHTTP(w, r)
}

func (app *App) Close() {
	app.DB.Close()
}

func (app *App) AddRoutes(routes Routes) {
	for _, route := range routes {
		handler := Logger(route.Handler(app))
		if route.Method == "OPTIONS" {
			app.router.
				// Match all url
				Methods(route.Method).
				Handler(handler)
		} else {
			app.router.
				// Name(route.Name).
				Methods(route.Method).
				Path(route.Pattern).
				Handler(handler)
		}

	}
}

func (app *App) Run() {
	log.Fatal(http.ListenAndServe(":"+app.PORT, app))
}

func (app *App) Request(method string, route string, body interface{}) *httptest.ResponseRecorder {
	var request = &http.Request{}
	switch t := body.(type) {
	case string:
		request, _ = http.NewRequest(method, route, strings.NewReader(t))
		if method == "POST" || method == "PUT" {
			if t != "" && t[0:1] == "{" {
				request.Header.Set("Content-Type", "application/json")
			} else {
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
			}
		}
	case BodyMultipart:
		request, _ = http.NewRequest(method, route, &t.Buff)
		request.Header.Set("Content-Type", t.ContentType)
	default:
		return nil
	}

	request.RemoteAddr = "127.0.0.1:8080"

	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	return response
}

func NewApp() *App {
	return &App{
		router: newRouter(),
		DB:     newDB(),
		PORT:   getPort(),
		ENV:    getEnv(),
	}
}

func getPort() string {
	env := os.Getenv("PORT")

	if env == "" {
		return "8080"
	}

	return env
}

func getEnv() string {
	env := os.Getenv("ENV")

	if env == "" {
		return "production"
	}

	return env
}

func newRouter() *mux.Router {
	return mux.NewRouter().StrictSlash(true)
}

func parseDB(database_url string) (driver string, gorm_arg string, err error) {
	re, _ := regexp.Compile(`(\w*):\/\/`)
	result := re.FindStringSubmatch(database_url)

	if result == nil {
		return "", "", errors.New("Can't find driver for DB")
	}

	driver = result[1]

	switch driver {
	case "postgres":
		return driver, database_url, nil
	case "mysql":
		re, _ := regexp.Compile(`(\w*):\/\/(.+@)([^/]+)(.+)`)
		result := re.FindStringSubmatch(database_url)
		return driver, result[2] + "tcp(" + result[3] + ")" + result[4] + "?parseTime=true", nil
	default:
		return "", "", errors.New("Driver is not supported")
	}
}

func newDB() *DB {
	driver, gorm_arg, err := parseDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(driver, gorm_arg)
	if err != nil {
		panic(err.Error())
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.DB().Ping()
	if err != nil {
		panic(err.Error())
	}

	return &DB{db.Debug()}
}
