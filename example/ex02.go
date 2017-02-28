package main

import (
	"fmt"

	"github.com/DavadDi/validation"
)

type Person struct {
	Name     string   `valid:"required"`
	Email    string   `valid:"required;email"`
	Age      int      `valid:"required;age"`
	Sex      int      ``
	WebSites []string `valid:"url"`
}

type AgeChecker struct {
}

func (ac *AgeChecker) Validater(v interface{}) error {
	age, ok := v.(int)
	if !ok {
		return validation.NewErrWrongType("int", v)
	}

	if age <= 0 || age > 140 {
		return fmt.Errorf("age checke failed. should between [1-140], now %d", age)
	}

	return nil
}

func main() {
	validation.Debug(true)

	validation.AddValidater("age", &AgeChecker{})

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
