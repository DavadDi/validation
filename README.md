# validation

 
The Simpler Validation in go. 

* Use **interface** for Valiater
* Support User define Validater
* Support Struct define validater interface
* Support slice/array/pointer and netestd struct validate. Not for map now!


## Install

```
$go get https://www.github.com/DavadDi/validation
```

## Simple Useage

```go

package main

import (
	"fmt"
	"validation"
)

// ex01 simple use
type Person struct {
	Name     string   `valid:"required"`
	Email    string   `valid:"required;email"`
	Age      int      `valid:"-"`
	Sex      int      ``
	WebSites []string `valid:"url"`
}

func main() {
	web1 := "http://www.do1618.com"
	web2 := "www.baidu.com"

	person1 := &Person{
		Name:     "dave",
		Email:    "dwh0403@163.com",
		WebSites: []string{web1, web2},
	}

	validater := validation.NewValidation()
	res := validater.Validate(person1)

	if !res {
		fmt.Printf("Person1 validate failed. %s", validater.ErrMsg())
	} else {
		fmt.Println("Person1 validate succeed!")
	}

	validater.Clear()
	person2 := &Person{
		Email:    "dwh0403@163.com",
		WebSites: []string{web1, web2},
	}

	res = validater.Validate(person2)

	if !res {
		fmt.Printf("Person2 validate failed. %s\n", validater.ErrMsg())
	} else {
		fmt.Println("Person2 validate succeed!")
	}
}
```

### Output:
	Person1 validate succeed!
	Person2 validate failed. [Name] check failed [field can't be empty or zero] [""]

## Add Use Define Validater

Use define validater need impl the interface **Validater(v interface{}) error** and add it to validation.

```go
package main

import (
	"fmt"
	"validation"
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

	if age < 0 || age > 140 {
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
```

### Output:

	Person1 validate failed. [Age] check failed [field can't be empty or zero] [0]

## Collaborate with Struct Interface

```go
package main

import (
	"fmt"
	"log"
	"validation"
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

```

### Output

	2017/02/28 10:29:34 In our struct validater now
	Person1 validate failed. [Object] check failed [age checke failed. should between [1-140], now 0] [&main.Person{Name:"dave", Email:"dwh0403@163.com", Age:0, Sex:0, WebSites:[]string(nil)}]

## Debug

Turn on Debug

```go
validation.EnableDebug(true)
```

Turn off Debug

```go
validation.EnableDebug(false)
```
