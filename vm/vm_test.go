package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"strconv"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElm := vm.LastPopppedStackElem()

		testExpectedObject(t, tt.expected, stackElm)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}

	case float64:
		err := testFloatObject(float64(expected), actual)
		if err != nil {
			t.Errorf("testFloatObject failed: %s", err)
		}

	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T(+%v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d want=%d", result.Value, expected)
	}

	return nil
}

func testFloatObject(expected float64, actual object.Object) error {
	result, ok := actual.(*object.Float)
	if !ok {
		return fmt.Errorf("object is not float. got=%T(+%v)", actual, actual)
	}

	resultStr := strconv.FormatFloat(result.Value, 'f', -1, 64)
	expectedStr := strconv.FormatFloat(expected, 'f', -1, 64)
	if resultStr != expectedStr {
		return fmt.Errorf("object has wrong value. got=%f want=%f", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T(+%v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t want=%t", result.Value, expected)
	}
	return nil
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"10 % 2", 0},
		{"(10 % 5) + 1", 1},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, tests)
}

func TestFloatArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1.0", 1.0},
		{"2.0", 2.0},
		{"1.0 + 2", 3.0},
		{"1 - 2.0", -1.0},
		{"1.0 * 2", 2.0},
		{"4.0 / 2", 2.0},
		{"50.0 / 2 * 2 + 10 - 5", 55.0},
		{"5 + 5 + 5 + 5.0 - 10", 10.0},
		{"2 * 2 * 2 * 2 * 2.0", 32.0},
		{"5 * (2 + 10.0)", 60.0},
		{"10.0 % 2", 0.0},
		{"(10 % 5) + 1.0", 1.0},
		{"5.1 + 5 + 5 + 5 - 10", 10.100000000000001},
		{"20 + 2 * -10.0", 0.0},
		{"(5 + 10 * 2 + 15 % 3.0) * 2 + -10", 40.0},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"1.1 != 1.1", false},
		{"1.121 > 1.1", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	runVmTests(t, tests)
}
