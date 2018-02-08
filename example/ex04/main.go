package main

import (
	"fmt"

	"github.com/DavadDi/validation"
)

// Bool convert bool -> *bool
func Bool(a bool) *bool {
	return &a
}

// Person for test struct
type Person struct {
	Name     string `valid:"required"`
	Email    string `valid:"required;email"`
	Age      int
	Sex      int
	IsAdmin  *bool    `valid:"required"`
	WebSites []string `valid:"url"`
}

func main() {
	// Turn on debug
	validation.EnableDebug(true)
	person1 := &Person{
		Name:  "dave",
		Email: "dwh0403@163.com",

		IsAdmin: Bool(false),
	}

	validater := validation.NewValidation()
	res := validater.Validate(person1)

	if res {
		fmt.Println("Person1 validate succeed!")
	} else {
		fmt.Printf("Person1 validate failed. %s\n", validater.ErrMsg())
	}
}
