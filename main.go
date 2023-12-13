package main

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

}
