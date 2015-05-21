package inject

import (
	"reflect"
	"testing"
)

func equal(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_TypeOf(t *testing.T) {
	var str string
	var ptr *string
	equal(t, TypeOf(ptr), TypeOf(str))
}

func Test_Map(t *testing.T) {
	injector := New()
	injector.Map("Hello World")
	val, _ := injector.Get(TypeOf(""))
	equal(t, val.String(), "Hello World")
}

func Test_MapTag(t *testing.T) {
	injector := New()
	injector.MapTag("Hello World", "name")
	val, _ := injector.GetTag("name")
	equal(t, val.String(), "Hello World")
}

func Test_Invoke(t *testing.T) {
	injector := New()
	str := "Ciel"
	injector.Map(str)
	fn := func(name string) {
		if name != str {
			t.Errorf("Expected %s - Got %s", str, name)
		}
	}
	injector.Invoke(fn)
}

func Test_Inject(t *testing.T) {
	type user struct {
		Name     string `inject:"name"`
		Password string `inject:"password"`
		Usertype string
	}
	injector := New()

	usertype := "normal"
	injector.Map(usertype)
	name := "Ciel"
	password := "123456"
	injector.MapTag(name, "name")
	injector.MapTag(password, "password")

	u := user{}
	injector.Inject(&u)
	equal(t, u.Name, name)
	equal(t, u.Password, password)
	equal(t, u.Usertype, usertype)
}
