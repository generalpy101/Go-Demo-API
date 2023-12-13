package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const NotFoundJSONReponse = `{
	"error": 404,
	"message": "Not found"
}`

// Model for course and author
// TODO: Move in separate file
type Course struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	// Store reference so that we get original author instead of copy
	Author *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// DB imitation using slice
var courses []Course

// middleware or helper functions; TODO: Move in separate file
func (c *Course) IsEmpty() bool {
	/*
		Checks if a course is empty or not
		Checks using id and name

		Returns:
			True if empty
		Params:
			c: *Course
	*/
	if c.Id == "" && c.Name == "" {
		return true
	}
	return false
}

func main() {
	// Init router
	router := mux.NewRouter()

	// Mock data
	author1 := Author{
		Fullname: "John Doe",
		Website:  "https://johndoe.com",
	}
	author2 := Author{
		Fullname: "Jane Doe",
		Website:  "https://janedoe.com",
	}

	courses = append(courses, Course{
		Id:     "1",
		Name:   "Course 1",
		Price:  10.99,
		Author: &author1,
	}, Course{
		Id:     "2",
		Name:   "Course 2",
		Price:  20.99,
		Author: &author2,
	})

	// Route handlers
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/courses", getCourses).Methods("GET")
	router.HandleFunc("/courses/{id}", getCourse).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", router))
}

// Utils
func raiseError(err error) {
	/*
		Raises error
	*/
	if err != nil {
		panic(err)
	}
}

// Controllers; TODO: Move in separate file

// home route
func homePage(w http.ResponseWriter, r *http.Request) {
	/*
		Home route

		Params:
			w: http.ResponseWriter
			r: *http.Request

		Returns:
			nil
	*/
	w.Write([]byte("Welcome to the api"))
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	/*
		Get all courses

		Params:
			w: http.ResponseWriter
			r: *http.Request

		Returns:
			nil
	*/
	// response := strings.Builder{}

	// jsonCourses, err := json.Marshal(courses)

	// raiseError(err)

	// response.Write(jsonCourses)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(response.String()))

	// Better way
	w.Header().Set("Content-Type", "application/json")
	// NewEncoder returns a new encoder that writes to w
	err := json.NewEncoder(w).Encode(courses)
	raiseError(err)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	/*
		Get a single course

		Params:
			w: http.ResponseWriter
			r: *http.Request
		Returns:
			nil
	*/
	// Get id from url
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	// Using slice so have to use loop
	for _, course := range courses {
		if course.Id == params["id"] {
			err := json.NewEncoder(w).Encode(course)
			raiseError(err)
		}
	}
	// If not found
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(NotFoundJSONReponse))
}
