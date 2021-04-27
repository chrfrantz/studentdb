package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"studentdb/studentdb"
)

func handlerDefault(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Please use '/hello' or '/student' path.\n" +
		"If using '/student' path, use following POST payload structure:\n" +
		"{\n    \"name\": \"Student\",\n    \"age\": 21,\n    \"studentid\": \"someId\"\n}")
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		status := 400
		http.Error(w, http.StatusText(status), status)
		return
	}
	name := parts[2]
	fmt.Fprintln(w, parts)
	fmt.Fprintf(w, "Hello %s %s!\n", name, parts[3])
}

// -----------------

/*
Main function.

Test with POST invocation to <ip>:<port>/student

{
    "name": "Student",
    "age": 21,
    "studentid": "someId"
}

*/
func main() {

	if os.Getenv("DB_HOST") == "" {
          fmt.Println("Using in-memory database")
          // Using in-memory storage
          studentdb.Global_db = &studentdb.StudentsDB{}
        } else {
	  fmt.Println("Using MongoDB instance ", os.Getenv("DB_HOST"))
	  // Using MongoDB based storage
	  studentdb.Global_db = &studentdb.StudentsMongoDB{
		os.Getenv("DB_HOST"), //"mongodb://localhost",
		"studentsDB",
		"students",
	  }
        }

	studentdb.Global_db.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handlerDefault)
	http.HandleFunc("/hello/", handlerHello)
	http.HandleFunc("/student/", studentdb.HandlerStudent)
	fmt.Println("Listening on port "+port)
	http.ListenAndServe(":"+port, nil)
}
