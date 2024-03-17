package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "hello world"}
	hello2 := &String{Value: "hello world"}
	diff1 := &String{Value: "my name is tuna"}
	diff2 := &String{Value: "my name is tuna"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Error("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Error("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Error("strings with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Error("trues do not have same hash key")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Error("falses do not have same hash key")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Error("true has same hash key as false")
	}
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}
	two1 := &Integer{Value: 2}
	two2 := &Integer{Value: 2}

	if one1.HashKey() != one2.HashKey() {
		t.Error("integers with same content have twoerent hash keys")
	}

	if two1.HashKey() != two2.HashKey() {
		t.Error("integers with same content have twoerent hash keys")
	}

	if one1.HashKey() == two1.HashKey() {
		t.Error("integers with twoerent content have same hash keys")
	}
}
