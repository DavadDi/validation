/*
Package validation implements a simple library for struct tag validte.

	1. Use interface for Validater
	2. Support User define Validater
	3. Support Struct define validater interface
	4. Support slice/array/pointer and netestd struct validate. Not for map now!
*/
package validation

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
)

// ex01 simple use
//
//	import (
//		"fmt"
//		"validation"
//	)
//
//	type Person struct {
//		Name     string   `valid:"required"`
//		Email    string   `valid:"required;email"`
//		Age      int      `valid:"-"`
//		Sex      int      ``
//		WebSites []string `valid:"url"`
//	}
//
//	func main() {
//		web1 := "http://www.do1618.com"
//		web2 := "www.baidu.com"
//
//		person1 := &Person{
//			Name:     "dave",
//			Email:    "dwh0403@163.com",
//			WebSites: []string{web1, web2},
//		}
//
//		validater := validation.NewValidation()
//		res := validater.Validate(person1)
//
//		if !res {
//			fmt.Printf("Person1 validate failed. %s", validater.ErrMsg())
//		} else {
//			fmt.Println("Person1 validate succeed!")
//		}
//
//		validater.Reset()
//		person2 := &Person{
//			Email:    "dwh0403@163.com",
//			WebSites: []string{web1, web2},
//		}
//
//		res = validater.Validate(person2)
//
//		if !res {
//			fmt.Printf("Person2 validate failed. %s\n", validater.ErrMsg())
//		} else {
//			fmt.Println("Person2 validate succeed!")
//		}
//	}

// Validater Interface For use define stuct, without params
// If validate struct has this interface, we will call this interface firstly
type Validater interface {
	Validater() error
}

// ValidaterFunc type
type ValidaterFunc func(v interface{}) error

// Interface for Validators
//type Validater interface {
//	Validater(v interface{}) error
//}

// Global var for validation pkg
const (
	ValidTag      = "valid"    // Validater tag name
	FuncSeparator = ";"        // Func sparator "required;email"
	ValidIgnor    = "-"        // Igore for validater
	RequiredKey   = "required" // required key for not empty value
)

var (
	// Init by this pkg. no need rwlock
	validatorsMap = map[string]ValidaterFunc{
		RequiredKey: requiredChecker,
		"email":     emailChecker,
		"url":       urlChecker,
	}

	// Using rwlock avoid race
	customValidatorsMap = newCustomValidators()

	// Output debug info or not
	debugFlag = false
)

// Pkg debug function, just a wrapper for log
func debug(arg interface{}) {
	if debugFlag {
		log.Print(arg)
	}
}

func debugf(format string, args ...interface{}) {
	if debugFlag {
		log.Printf(format, args...)
	}
}

// EnableDebug enable validation debug log
func EnableDebug(flag bool) {
	debugFlag = flag
}

// Return a new custom validater manager
func newCustomValidators() *CustomValidators {
	return &CustomValidators{
		validatorsMap: make(map[string]ValidaterFunc),
	}
}

// AddValidater Function export to add user define Validater
func AddValidater(name string, validater ValidaterFunc) error {
	if customValidatorsMap == nil {
		customValidatorsMap = newCustomValidators()
	}

	return customValidatorsMap.AddValidater(name, validater)
}

// CustomValidators Because user can add user define validater, avoid data race, add rwlock
type CustomValidators struct {
	validatorsMap map[string]ValidaterFunc
	sync.RWMutex
}

// AddValidater If name conflict with CustomValidators, replace.
// conflict with validatorsMap, return err
func (cvm *CustomValidators) AddValidater(name string, validater ValidaterFunc) error {
	// check validater
	if validater == nil {
		return ErrValidater
	}

	// check name conflict
	if validatorsMap[name] != nil {
		return ErrValidaterExists
	}

	cvm.Lock()
	cvm.validatorsMap[name] = validater
	cvm.Unlock()

	return nil
}

// Return user define validater for namego
func (cvm *CustomValidators) findValidater(name string) (v ValidaterFunc, ok bool) {
	cvm.RLock()
	v, ok = cvm.validatorsMap[name]
	cvm.RUnlock()

	return v, ok
}

// Validation err list
type Validation struct {
	Errors []*Error
}

// NewValidation create a new validation
func NewValidation() *Validation {
	return &Validation{}
}

// Errs Return Error list
func (mv *Validation) Errs() []*Error {
	return mv.Errors
}

// ErrMsg Return msg detail Error message
func (mv *Validation) ErrMsg() string {
	buf := bytes.NewBufferString("")

	for _, err := range mv.Errors {
		str := err.String()
		buf.WriteString(str)
	}

	return buf.String()
}

// Reset reset state to init
func (mv *Validation) Reset() {
	mv.clear()
}

// HasError Check has errors or not
func (mv *Validation) HasError() bool {
	return len(mv.Errors) != 0
}

// Validate entry function.
// True: Validiton passed.
// False: Validate don't passed, mv.ErrMsg() contains the detail info.
func (mv *Validation) Validate(obj interface{}) bool {
	if obj == nil {
		debug("obj == nil")
		return true
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	t := v.Type()

	// Here only accept structs
	if v.Kind() != reflect.Struct {
		err := &ErrOnlyStrcut{Type: v.Type()}
		mv.addError("Object", obj, err)
		return false
	}

	debugf("Check struct [%s]", t.Name())

	objvk, ok := obj.(Validater)
	if ok {
		err := objvk.Validater()
		if err != nil {
			mv.addError("Object", obj, err)
		}
	}

	for i := 0; i < v.NumField(); i++ {
		tf := t.Field(i) // type field
		vf := v.Field(i) // vaule field

		// Skip Anonymous and private field
		if !tf.Anonymous && len(tf.PkgPath) > 0 {
			continue
		}

		fns := mv.getValidFuns(tf, ValidTag)

		// Already skip ValidIgnor flag, such as "-"
		if len(fns) == 0 {
			continue
		}

		mv.typeCheck(vf, tf, v, false)
	}

	if mv.HasError() {
		return false
	}

	return true
}

func (mv *Validation) checkRequire(v reflect.Value, t reflect.StructField) error {
	rck := validatorsMap[RequiredKey]
	return rck(v.Interface())
}

func (mv *Validation) typeCheckPtr(v reflect.Value, t reflect.StructField, o reflect.Value) {

}

// Valid struct field type, if typeCheck is ptr, wo just check ptr for required, not for element
func (mv *Validation) typeCheck(v reflect.Value, t reflect.StructField, o reflect.Value, ignoreRequired bool) {
	fns := mv.getValidFuns(t, ValidTag)

	// skip
	if len(fns) == 0 {
		return
	}

	if ignoreRequired {
		delete(fns, RequiredKey)
	}

	// First check all field for required
	if fns[RequiredKey] != nil {
		if err := mv.checkRequire(v, t); err != nil {
			mv.addError(t.Name, v.Interface(), err)
		}

		delete(fns, RequiredKey)
	}

	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		debugf("\tCheck field [%s]", t.Name)

		for fname := range fns {
			debugf("CheckerName: [%s]", fname)

			// find custom map and pkg map
			var vck ValidaterFunc
			var find bool

			if vck, find = customValidatorsMap.findValidater(fname); !find {
				vck = validatorsMap[fname]
			}

			if vck == nil {
				err := fmt.Errorf("can't find checker for [%s]", fname)
				mv.addError(t.Name, v.Interface(), err)
				debugf("can't find checker for [%s]", fname)
				continue
			}

			err := vck(v.Interface())
			if err != nil {
				mv.addError(t.Name, v.Interface(), err)
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).Kind() != reflect.Struct {
				mv.typeCheck(v.Index(i), t, o, false)
			} else {
				mv.Validate(v.Index(i).Interface())
			}
		}

	case reflect.Interface:
		// If the value is an interface then encode its element
		if !v.IsNil() {
			mv.Validate(v.Elem())
		}

	case reflect.Ptr:
		// only check
		// If the value is a pointer then check its element
		if !v.IsNil() {
			mv.typeCheck(v.Elem(), t, o, true)
		}

	case reflect.Struct:
		mv.Validate(v.Interface())

	// case reflect.Map: // don't support map now
	default:
		err := fmt.Errorf("UnspportType %s", v.Type())
		mv.addError(t.Name, v.Interface(), err)
	}

	return
}

// Clear error
func (mv *Validation) clear() {
	mv.Errors = nil
}

// Apend error to validtion
func (mv *Validation) addError(key string, v interface{}, err error) {
	errtmp := &Error{FieldName: key, Value: v, Err: err}
	mv.Errors = append(mv.Errors, errtmp)
}

// Return fun names and params
func (mv *Validation) getValidFuns(tf reflect.StructField, tag string) map[string]interface{} {
	out := make(map[string]interface{})

	opt, ok := tf.Tag.Lookup(tag)
	if !ok || len(opt) == 0 || opt == ValidIgnor {
		return nil
	}

	for _, value := range strings.Split(opt, FuncSeparator) {
		// omit func has params for now
		value = strings.TrimSpace(value)
		out[value] = struct{}{}
	}

	return out
}
