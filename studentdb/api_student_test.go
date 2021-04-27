package studentdb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlerStudent_notImplemented(t *testing.T) {
	// instantiate mock HTTP server
	// register our handlerStudent <-- actual logic
	ts := httptest.NewServer(http.HandlerFunc(HandlerStudent))
	defer ts.Close()

	// create a request to our mock HTTP server
	//    in our case it means to create DELETE request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	// check if the response from the handler is what we expect
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusNotImplemented, resp.StatusCode)
	}
}

func Test_handlerStudent_malformedURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandlerStudent))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/student/id/extra",
		ts.URL + "/stud/",
	}
	for _, tstring := range testCases {
		resp, err := http.Get(tstring)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d", tstring,
				http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}

// GET /student/
// empty array back
func Test_handlerStudent_getAllStudents_empty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandlerStudent))
	defer ts.Close()
	Global_db = &StudentsDB{}
	Global_db.Init()

	resp, err := http.Get(ts.URL + "/student/")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}

	var a []interface{}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	if len(a) != 0 {
		t.Errorf("Excpected empty array, got %s", a)
	}
}

// GET /student/
// single Tom student back
func Test_handlerStudent_getAllStudents_Tom(t *testing.T) {
	Global_db = &StudentsDB{}
	Global_db.Init()
	testStudent := Student{"Tom", 21, "id0"}
	Global_db.Add(testStudent)

	ts := httptest.NewServer(http.HandlerFunc(HandlerStudent))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/student/")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}

	var a []Student
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	if len(a) != 1 {
		t.Errorf("Excpected array with one element, got %v", a)
	}

	if a[0].Name != testStudent.Name || a[0].Age != testStudent.Age || a[0].StudentID != testStudent.StudentID {
		t.Errorf("Students do not match! Got: %v, Expected: %v\n", a[0], testStudent)
	}
}

// GET /student/id0
// single Tom student back
func Test_handlerStudent_getStudent_Tom(t *testing.T) {
	Global_db = &StudentsDB{}
	Global_db.Init()
	testStudent := Student{"Tom", 21, "id0"}
	Global_db.Add(testStudent)

	ts := httptest.NewServer(http.HandlerFunc(HandlerStudent))
	defer ts.Close()

	// --------------
	resp, err := http.Get(ts.URL + "/student/id1")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusNotFound, resp.StatusCode)
		return
	}

	// --------------
	resp, err = http.Get(ts.URL + "/student/id0")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}

	var a Student
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	if a.Name != testStudent.Name || a.Age != testStudent.Age || a.StudentID != testStudent.StudentID {
		t.Errorf("Students do not match! Got: %v, Expected: %v\n", a, testStudent)
	}
}

func Test_handlerStudent_POST(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandlerStudent))
	defer ts.Close()

	Global_db = &StudentsDB{}
	Global_db.Init()

	// Testing empty body
	resp, err := http.Post(ts.URL+"/student/", "text/plain", strings.NewReader(" "))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Testing proper JSON body
	tom := "{ \"name\": \"Tom\", \"age\": 21, \"id\": \"id0\"}"

	resp, err = http.Post(ts.URL+"/student/", "application/json", strings.NewReader(tom))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, Body: %s",
			http.StatusOK, resp.StatusCode, all)
	}

	// Trying to add Tom second time
	resp, err = http.Post(ts.URL+"/student/", "application/json", strings.NewReader(tom))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, Body: %s",
			http.StatusBadRequest, resp.StatusCode, all)
	}

	// Testing malformed JSON body
	wrongTom := "{ \"namee\": \"Tom\", \"agee\": 21, \"id\": \"id0\"}"

	resp, err = http.Post(ts.URL+"/student/", "application/json", strings.NewReader(wrongTom))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, Body: %s",
			http.StatusBadRequest, resp.StatusCode, all)
	}
}
