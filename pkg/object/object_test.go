package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello, world!"}
	hello2 := &String{Value: "Hello, world!"}
	diff1 := &String{Value: "My name is Johnny"}
	diff2 := &String{Value: "My name is Johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	meaning1 := &Integer{Value: 42}
	meaning2 := &Integer{Value: 42}
	fifty1 := &Integer{Value: 50}
	fifty2 := &Integer{Value: 50}

	if meaning1.HashKey() != meaning2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if fifty1.HashKey() != fifty2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if meaning1.HashKey() == fifty1.HashKey() {
		t.Errorf("integers with different content have same hash keys")
	}
}
