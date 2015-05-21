// Package inject is a dependency injection library for golang.
//
// It is original created by Jeremy Saenz and this is a fork version of
// https://github.com/codegangsta/inject
//
// This fork add tag name injection and custom tag support.
//
//  package main
//
//  import "github.com/imwithye/inject"
//
//  func main() {
//  	type user struct {
//  		Name     string
//  		Password string `inject:"password"`
//  		Usertype string `inject:"usertype"`
//  	}
//  	injector := New()
//
//  	name := "Ciel"
//  	injector.Map(name)
//  	password := "123456"
//  	injector.MapTag(password, "password")
//  	usertype := "normal"
//   	injector.MapTag(usertype, "usertype")
//
//  	u := user{}
//  	injector.Inject(&u)
//  	// u.Name == "Ciel"
//  	// u.Password == "Ciel"
//  	// u.Usertype == "normal"
//
//  	fn := func(name string) {
//  		// name == "Ciel"
//  	}
//  	injector.Invoke(fn)
//  }
package inject

import (
	"errors"
	"reflect"
)

// Inject represents an interface for mapping and injecting dependencies into structs
// and function arguments.
type Inject interface {
	Injector
	Invoker
	TypeMapper
}

// Injector represents an interface for mapping dependencies to a struct.
type Injector interface {
	Inject(interface{}) error
	InjectTag(interface{}, string) error
}

// Invoker represents an interface for calling functions via reflection.
type Invoker interface {
	Invoke(interface{}) ([]reflect.Value, error)
	InvokeTag(...interface{}) ([]reflect.Value, error)
}

// TypeMapper represents an interface for mapping interface{} values based
// on type or struct tag.
type TypeMapper interface {
	Map(interface{}) TypeMapper
	MapTag(interface{}, string) TypeMapper
	Get(t reflect.Type) (reflect.Value, error)
	GetTag(string) (reflect.Value, error)
}

// TypeOf returns a type for a given value. If the given type is a pointer, it
// will return its object's type.
func TypeOf(val interface{}) reflect.Type {
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// ValueOf returns a value for a given value. If the given type is a pointer, it
// will return its object's value.
func ValueOf(val interface{}) reflect.Value {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

type injector struct {
	tag     string
	typeMap map[reflect.Type]reflect.Value
	tagMap  map[string]reflect.Value
}

// New creates a new Inject with "inject" tag. Struct can use "inject" tag if
// it wants to be injected.
func New() Inject {
	return NewTag("inject")
}

// NewTag creates a new Inject with given tag. Struct can use the given tag if
// it wants to be injected.
func NewTag(tag string) Inject {
	return &injector{
		tag:     tag,
		typeMap: make(map[reflect.Type]reflect.Value),
		tagMap:  make(map[string]reflect.Value),
	}
}

func (i *injector) Inject(stc interface{}) error {
	return i.InjectTag(stc, i.tag)
}

func (i *injector) InjectTag(stc interface{}, tag string) error {
	v := ValueOf(stc)
	t := TypeOf(stc)
	if v.Kind() != reflect.Struct {
		return nil
	}
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		stcField := t.Field(j)
		if f.CanSet() {
			var (
				value reflect.Value
				err   error
			)
			if string(stcField.Tag) == tag || stcField.Tag.Get(tag) != "" {
				value, err = i.GetTag(stcField.Tag.Get(tag))
			} else {
				value, err = i.Get(f.Type())
			}
			if err != nil {
				return err
			}
			f.Set(value)
		}
	}
	return nil
}

func (i *injector) Invoke(fn interface{}) ([]reflect.Value, error) {
	t := reflect.TypeOf(fn)
	var in = make([]reflect.Value, t.NumIn())

	for j := 0; j < t.NumIn(); j++ {
		argType := t.In(j)
		val, err := i.Get(argType)
		if err != nil {
			return nil, err
		}
		in[j] = val
	}

	return reflect.ValueOf(fn).Call(in), nil
}

func (i *injector) InvokeTag(vals ...interface{}) ([]reflect.Value, error) {
	if len(vals) == 0 {
		return nil, nil
	}
	if len(vals) == 1 {
		return i.Invoke(vals[0])
	}

	fn := vals[len(vals)-1]
	tags := vals[:len(vals)-1]

	t := reflect.TypeOf(fn)
	if len(tags) > t.NumIn() {
		tags = tags[:t.NumIn()]
	}
	var in = make([]reflect.Value, t.NumIn())

	for j := 0; j < len(tags); j++ {
		argTag := tags[j].(string)
		val, err := i.GetTag(argTag)
		if err != nil {
			return nil, err
		}
		in[j] = val
	}

	for j := len(tags); j < t.NumIn(); j++ {
		argType := t.In(j)
		val, err := i.Get(argType)
		if err != nil {
			return nil, err
		}
		in[j] = val
	}

	return reflect.ValueOf(fn).Call(in), nil
}

func (i *injector) Map(val interface{}) TypeMapper {
	i.typeMap[TypeOf(val)] = reflect.ValueOf(val)
	return i
}

func (i *injector) MapTag(val interface{}, tag string) TypeMapper {
	i.tagMap[tag] = reflect.ValueOf(val)
	return i
}

func (i *injector) Get(t reflect.Type) (reflect.Value, error) {
	val := i.typeMap[t]
	if val.IsValid() {
		return val, nil
	}

	if t.Kind() == reflect.Interface {
		for k, v := range i.typeMap {
			if k.Implements(t) {
				val = v
				break
			}
		}
	}

	if !val.IsValid() {
		return val, errors.New("value is not valid")
	}

	return val, nil
}

func (i *injector) GetTag(tag string) (reflect.Value, error) {
	val := i.tagMap[tag]
	var err error
	if !val.IsValid() {
		err = errors.New("value is not valid")
	}
	return val, err
}
