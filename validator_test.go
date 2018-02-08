package validation

import (
	"fmt"
	"reflect"
	"testing"
)

// Person for test struct
type Person struct {
	Name     *string   `valid:"required"`
	Email    string    `valid:"required;email"`
	Age      int       `valid:"-"`
	Sex      int       ``
	WebSites []*string `valid:"url"`
}

func (p *Person) TValidater() error {
	return nil
}

var (
	web1   = "http://www.do1618.com"
	web2   = "www"
	name   = "dave"
	person = &Person{Name: &name, Email: "dwh0403@163.com", WebSites: []*string{&web1, &web2}}
)

func TestGetFuns(t *testing.T) {
	v := reflect.ValueOf(*person)

	tests := []struct {
		Name       string
		ExpectNum  int
		ExpectStrs []string
	}{
		{"Name", 1, []string{"required"}},
		{"Email", 2, []string{"required", "email"}},
		{"Age", 0, nil}, // - skip count
		{"Sex", 0, nil},
	}

	typ := v.Type()
	valider := NewValidation()

	for _, test := range tests {
		field, ok := typ.FieldByName(test.Name)
		if !ok {
			t.Fatalf("Get field by [%s] failed", test.Name)
		}

		// Test exist tag
		// funct list
		fns := valider.getValidFuns(field, ValidTag)

		if len(fns) != test.ExpectNum {
			t.Errorf("Get tag vliad failed. should get [%d], got [%d] %v",
				test.ExpectNum, len(fns), fns)
		}

		// Test fun name list
		for _, name := range test.ExpectStrs {
			if _, ok := fns[name]; !ok {
				t.Errorf("Func name [%s] should in list. but got %v", name, fns)
			}
		}
	}
}

func TestValidation(t *testing.T) {
	valider := NewValidation()
	res := valider.Validate(person)

	// expect www failed.
	if res {
		t.Errorf("Validate person should failed for www:\n")
	}
}

// Test struct interface
func TestValidationIf(t *testing.T) {
	validor := NewValidation()
	res := validor.Validate(person)
	if res {
		t.Errorf("TestValidationIf failed:\n %s", validor.ErrMsg())
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		Email  string
		Expect bool
	}{
		{"dwh0403@163.com", true},
		{"aaa@", false},
		{"dwh0403", false},
		{"13455", false},
	}

	for _, test := range tests {
		err := emailChecker(test.Email)
		res := (err == nil)

		if res != test.Expect {
			t.Errorf("Email check failed. [%s] should [%t]",
				test.Email, test.Expect)
		}
	}
}

type AddFile struct {
	FileName string `valid:"required"`
	URL      string `valid:"-"`
}

func TestNestedStruct(t *testing.T) {
	// EnableDebug(true)

	person := struct {
		Name     string     `valid:"required"`
		Email    string     `valid:"required;email"`
		Age      int        `valid:"-"`
		Sex      int        ``
		WebSites []*string  `valid:"required;url"`
		FileAdd  []*AddFile `valid:"required"`
		FileDel  []AddFile  `valid:"required"`
	}{
		Name:    "dave",
		Email:   "aa@aa.com",
		FileAdd: []*AddFile{{FileName: "file1"}, {FileName: "file2"}},
		FileDel: []AddFile{{FileName: "file1"}, {FileName: "file2"}},
	}

	validor := NewValidation()
	res := validor.Validate(person)

	if !res {
		debugf("TestNestedStruct: Error %s", validor.ErrMsg())
	}

	if res {
		t.Errorf("TestNestedStruct should failed\n")
	}

	EnableDebug(false)
}

func upperChecker(v interface{}) error {
	name, ok := v.(string)
	if !ok {
		return NewErrWrongType("string", v)
	}

	first := name[0]
	if !(first > 'A' && first < 'Z') {
		return fmt.Errorf("in name supper checker, name frist letter should upper")
	}

	return nil
}

func Bool(a bool) *bool {
	return &a
}

func TestPtrRequired(t *testing.T) {

	person := struct {
		Name     string    `valid:"required"`
		Email    string    `valid:"required;email"`
		Admin    *bool     `valid:"required"`
		Sex      int       ``
		WebSites []*string `valid:"url"`
	}{
		Name:  "Dave",
		Email: "aa@aa.com",
		Admin: Bool(false),
	}

	validor := NewValidation()
	res := validor.Validate(person)

	if !res {
		t.Errorf("TestPtrRequired should succeed %s\n", validor.ErrMsg())
	}
}

func TestCustomValidater(t *testing.T) {
	AddValidater("upper", upperChecker)

	person := struct {
		Name     string    `valid:"required;upper"`
		Email    string    `valid:"required;email"`
		Age      int       `valid:"-"`
		Sex      int       ``
		WebSites []*string `valid:"url"`
	}{
		Name:  "dave",
		Email: "aa@aa.com",
	}

	validor := NewValidation()
	res := validor.Validate(person)

	if res {
		t.Errorf("TestCustomValidater should failed\n")
	}
}

func TestCustomValidaterConflict(t *testing.T) {
	err := AddValidater("upper", upperChecker)

	// should replace
	if err != nil {
		t.Errorf("AddValidater should succeed. but got %s\n", err.Error())
	}

	err = AddValidater("email", upperChecker)
	if err != ErrValidaterExists {
		t.Errorf("AddValidater should failed [ErrValidaterExists]. but got %s\n", err.Error())
	}
}
