package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
	"os"
	"tweethour"
)

var s *server

const NoStatusCode = 0

func rootHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Try /hello/name")
}

func nameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	fmt.Fprintf(w, "Hello %s", html.EscapeString(name))
}

func histogramHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	username := vars["username"]

	hourTweets, err := s.h.Get(username)

	if err != nil {
		handleError(w, err)
		return
	}

	json, marshalErr := json.MarshalIndent(hourTweets, "", "    ")

	if marshalErr != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(json))
}

func handleRequests(port string) {

	router := mux.NewRouter().StrictSlash(true)
	// calling .Methods("GET") returns 404 not found, where as it should be returning 405 method not allowed
	// hence the wrapper function to return 405 in code
	router.HandleFunc("/", filter(rootHandler))
	router.HandleFunc("/hello/{name}", filter(nameHandler))
	router.HandleFunc("/histogram/{username}", filter(histogramHandler))

	port = ":" + port
	log.Fatal(http.ListenAndServe(port, router))
}

func disallowMethods(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "%s not allowed, try GET", r.Method)
}

type Description struct {
	Get GET `json:"GET"`
}

type GET struct {
	Message string `json:"description"`
}

func newGetDescription(message string) Description {
	g := GET{message}
	d := Description{g}
	return d
}

func describe(w http.ResponseWriter, r *http.Request) {
	desc := newGetDescription(fmt.Sprintf("GET: %s", r.URL.Path))
	s, _ := json.MarshalIndent(desc, "", "    ")
	fmt.Fprintf(w, string(s))

}

func filter(fn http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// setting allowed methods returns 404 not found, where as it should be returning 405 method not allowed
		// hence the wrapper function to return 405 method not allowed

		// "OPTIONS" - describe
		if r.Method == "GET" {
			fn(w, r)
		} else if r.Method == "OPTIONS" {
			describe(w, r)
		} else {
			disallowMethods(w, r)
		}
	}
}

type server struct {
	h tweethour.Histogram
}

func NewServer() *server {
	s := new(server)
	var err tweethour.Error
	s.h, err = tweethour.NewHistogram()

	if err != nil {
		tokenErr, ok := err.(*tweethour.TokenError)

		if ok {
			fmt.Println("Fatal error: Unable to obtain auth token")
			fmt.Println(tokenErr)
			os.Exit(1)
		}
	}
	return s
}

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("Usage: Tweethour <port-number>")
		return
	}

	port := args[1]

	fmt.Println("Setting up...")
	fmt.Println("Obtaining auth token...")

	// obtain the toekn once and retain it
	s = NewServer()

	fmt.Printf("Starting to listen on : %s", port)

	handleRequests(port)
}
