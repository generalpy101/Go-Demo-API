package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const CourseIdLength = 10
const CourseIdPrefix = "co"

const IdCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Error response struct
type ErrorResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

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
	return c.Name == ""
}

func (c *Course) GenerateId() {
	/*
		Generates and sets id for a course
	*/
	c.Id = generateRandomStringOfLength(CourseIdPrefix, CourseIdLength)
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
	router.HandleFunc("/courses", createCourse).Methods("POST")
	router.HandleFunc("/courses/{id}", deleteCourse).Methods("DELETE")
	router.HandleFunc("/courses/{id}", updateCourse).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", router))
}

// Utils
func RaiseError(err error, w http.ResponseWriter) {
	/*
		Raises error
	*/
	if err != nil {
		// Send internal server error
		internalServerError := ErrorResponse{
			Error:   http.StatusInternalServerError,
			Message: "Internal server error",
			Detail:  "Something went wrong, server is shutting down",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerError)
		panic(err)
	}
}

func generateRandomStringOfLength(prefix string, length int) string {
	/*
		Generates random string of given length

		Params:
			prefix: string
			length: string
		Returns:
			string
	*/
	randomString := strings.Builder{}

	randomString.WriteString(prefix)

	for i := 0; i < length; i++ {
		// Get random index from charset
		randomString.WriteByte(IdCharset[rand.Intn(len(IdCharset))])
	}

	return randomString.String()
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
	RaiseError(err, w)
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
			RaiseError(err, w)
			return
		}
	}
	// If not found
	notFoundResponse := ErrorResponse{
		Error:   http.StatusNotFound,
		Message: "Course not found",
		Detail:  "Course with given id not found",
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(notFoundResponse)
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	/*
		Create a course

		Params:
			w: http.ResponseWriter
			r: *http.Request
		Returns:
			nil
	*/
	// Empty body
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		emptyBodyResponse := ErrorResponse{
			Error:   http.StatusBadRequest,
			Message: "Empty body",
			Detail:  "Provide body with course details in json",
		}
		json.NewEncoder(w).Encode(emptyBodyResponse)
		return
	}

	// Init course
	var course Course

	// Decode json body
	err := json.NewDecoder(r.Body).Decode(&course)
	RaiseError(err, w)

	if course.IsEmpty() {
		w.WriteHeader(http.StatusBadRequest)
		emptyBodyResponse := ErrorResponse{
			Error:   http.StatusBadRequest,
			Message: "Empty body",
			Detail:  "Course name and id are mandatory",
		}
		json.NewEncoder(w).Encode(emptyBodyResponse)
		return
	}

	// Generate id
	course.GenerateId()

	courses = append(courses, course)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}

func deleteCourse(w http.ResponseWriter, r *http.Request) {
	/*
		Delete a course

		Params:
			w: http.ResponseWriter
			r: *http.Request
		Returns:
			nil
	*/
	// Get Id from url
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id := params["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		emptyBodyResponse := ErrorResponse{
			Error:   http.StatusBadRequest,
			Message: "Empty id",
			Detail:  "Provide id in url",
		}
		json.NewEncoder(w).Encode(emptyBodyResponse)
		return
	}

	for index, course := range courses {
		if course.Id == id {
			// Storing course to return
			course = courses[index]
			// Remove course from slice
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	// If not found
	notFoundResponse := ErrorResponse{
		Error:   http.StatusNotFound,
		Message: "Course not found",
		Detail:  "Course with given id not found",
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(notFoundResponse)
}

func updateCourse(w http.ResponseWriter, r *http.Request) {
	/*
		Update a course
		Not using patch because we are updating whole course
		Body should contain details we want to update except id which is in url
		Can also get all items in body and update course includindg id

		Params:
			w: http.ResponseWriter
			r: *http.Request
		Returns:
			nil
	*/
	// Get Id from url
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id := params["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		emptyBodyResponse := ErrorResponse{
			Error:   http.StatusBadRequest,
			Message: "Empty id",
			Detail:  "Provide id in url to update course",
		}
		json.NewEncoder(w).Encode(emptyBodyResponse)
		return
	}

	for index, course := range courses {
		if course.Id == id {
			course := courses[index]
			// Remove course from slice
			courses = append(courses[:index], courses[index+1:]...)
			// Decode json body
			err := json.NewDecoder(r.Body).Decode(&course)
			RaiseError(err, w)
			// Set id to old course id since this is update
			course.Id = id
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	// If not found
	notFoundResponse := ErrorResponse{
		Error:   http.StatusNotFound,
		Message: "Course not found",
		Detail:  "Course with given id not found",
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(notFoundResponse)
}
