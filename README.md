validation
=====================
 
The Simpler Validation in go. 

* Use **func(v interface{}) error** for Validater
* Support User define Validater
* Support Struct define **Validater() error** interface
* Support slice/array/pointer and netestd struct validate. Not for map now!

[![Build Status](http://img.shields.io/travis/DavadDi/validation.svg?style=flat-square)](https://travis-ci.org/DavadDi/validation)
[![Coverage Status](http://img.shields.io/coveralls/DavadDi/validation.svg?style=flat-square)](https://coveralls.io/r/DavadDi/validation)

## Install and tests

Install:

```
$go get github.com/DavadDi/validation
```

Test:

```
$go test github.com/DavadDi/validation
```


## Simple Usage

```go

package main

import (
	"fmt"

	"github.com/DavadDi/validation"
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

	if res {
		fmt.Println("Person1 validate succeed!")
	} else {
		fmt.Printf("Person1 validate failed. %s", validater.ErrMsg())
	}

	validater.Reset()
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

#### Struct Tag Functions:
	required
	email
	url

### Output:
	Person1 validate succeed!
	Person2 validate failed. [Name] check failed [field can't be empty or zero] [""]

## Add Use Define Validater

Use define validater **func(v interface{}) error** and add it to validation.

```go
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

func ageChecker(v interface{}) error {
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
	validation.AddValidater("age", ageChecker)

	person1 := &Person{
		Name:  "dave",
		Email: "dwh0403@163.com",
	}

	validater := validation.NewValidation()
	res := validater.Validate(person1)

	if res {
		fmt.Println("Person1 validate succeed!")
	} else {
		fmt.Printf("Person1 validate failed. %s\n", validater.ErrMsg())
	}
}

```

### Output:

	Person1 validate failed. [Age] check failed [field can't be empty or zero] [0]

## Collaborate with Struct Interface
Use struct define validater need impl the interface **Validater() error**.

```go

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

func (p *Person) Validater() error {
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

	if res {
		fmt.Println("Person1 validate succeed!")
	} else {
		fmt.Printf("Person1 validate failed. %s\n", validater.ErrMsg())
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

## LICENSE

BSD License http://creativecommons.org/licenses/BSD/

