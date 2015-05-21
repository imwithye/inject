package inject

import (
	"errors"
	"reflect"
)

type Inject interface {
	Injector
	Invoker
	TypeMapper
}

type Injector interface {
	Inject(interface{}) error
	InjectTag(interface{}, string) error
}

type Invoker interface {
	Invoke(interface{}) ([]reflect.Value, error)
}

type TypeMapper interface {
	Map(interface{}) TypeMapper
	MapTag(interface{}, string) TypeMapper
	Get(t reflect.Type) (reflect.Value, error)
	GetTag(string) (reflect.Value, error)
}

func TypeOf(val interface{}) reflect.Type {
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

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

func New() Inject {
	return NewTag("inject")
}

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
