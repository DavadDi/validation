package main

import (
	"fmt"
	"log"

	"github.com/DavadDi/validation"
)

type Person struct {
	Name     string `valid:"required"`
	Email    string `valid:"required;email"`
	Age      int
	Sex      int
	WebSites []string `valid:"url"`
}

func (p *Person) TValidater() error {
	log.Println("In our struct validater now")
	if p.Age <= 0 || p.Age > 140 {
		return fmt.Errorf("age checke failed. should between [1-140], now %d", p.Age)
	}

	return nil
}

func main() {
	// Turn on debug
	// validation.EnableDebug(true)
	person1 := &Person{
		Name:  "dave",
		Email: "dwh0403@163.com",
	}

	validater := validation.NewValidation()
	res := validater.Validate(person1)

	if !res {
		fmt.Printf("Person1 validate failed. %s\n", validater.ErrMsg())
	} else {
		fmt.Println("Person1 validate succeed!")
	}
}
