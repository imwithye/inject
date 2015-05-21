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

func Test_InvokeTag(t *testing.T) {
	injector := New()
	injector.Map("TOKEN: 19940730")
	injector.MapTag("User", "name")
	injector.MapTag("Password", "password")
	injector.MapTag("Male", "gender")
	fn := func(name, password, gender, token string) {
		if name != "User" {
			t.Errorf("Expected %s - Got %s", "User", name)
		}
		if password != "Password" {
			t.Errorf("Expected %s - Got %s", "Password", password)
		}
		if gender != "Male" {
			t.Errorf("Expected %s - Got %s", "Male", gender)
		}
		if token != "TOKEN: 19940730" {
			t.Errorf("Expected %s - Got %s", "TOKEN: 19940730", gender)
		}
	}
	injector.InvokeTag("name", "password", "gender", fn)
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
	injector.Apply(&u)
	equal(t, u.Name, name)
	equal(t, u.Password, password)
	equal(t, u.Usertype, usertype)
}

func Test_InjectCustomTag(t *testing.T) {
	type user struct {
		Name     string `db:"name"`
		Password string `db:"password"`
		Usertype string
	}
	injector := NewTag("db")

	usertype := "normal"
	injector.Map(usertype)
	name := "Ciel"
	password := "123456"
	injector.MapTag(name, "name")
	injector.MapTag(password, "password")

	u := user{}
	injector.Apply(&u)
	equal(t, u.Name, name)
	equal(t, u.Password, password)
	equal(t, u.Usertype, usertype)
}
