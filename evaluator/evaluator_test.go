package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"strconv"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"50 % 2", 0},
		{"50 % 3 * 1", 2},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.1", 5.1},
		{"-5.1", -5.1},
		{"-10.2", -10.2},
		{"5.1 + 5 + 5 + 5 - 10", 10.100000000000001},
		{"-50 + 100.1", 50.099999999999994},
		{"20 + 2 * -10.0", 0.0},
		{"(5 + 10 * 2 + 15 % 3.0) * 2 + -10", 40.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	resultStr := strconv.FormatFloat(result.Value, 'f', -1, 64)
	expectedStr := strconv.FormatFloat(expected, 'f', -1, 64)
	if resultStr != expectedStr {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1.1 != 1.1", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		expected interface{}
		input    string
	}{
		{10, "if (true) { 10 }"},
		{nil, "if (false) { 10 }"},
		{10, "if (1) { 10 }"},
		{10, "if (1 < 2) { 10 }"},
		{nil, "if (1 > 2) { 10 }"},
		{20, "if (1 > 2) { 10 } else { 20 }"},
		{10, "if (1 < 2) { 10 } else { 20 }"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`if (11 > 1) {
        if (10 > 1) {
        return 10;
        }
      return 1;
      }`, 10,
		},
		{
			`
      let f = fn(x) {
        return x;
        x + 10;
      };
      f(10);`,
			10,
		},
		{
			`
      let f = fn(x) {
        let result = x + 10;
        return result;
        return 10;
      };
      f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"1 = 2;",
			"invalid identifier: 1",
		},
		{
			"let a = b;",
			"identifier not found: b",
		},
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5.0 + true; 5;",
			"type mismatch: FLOAT + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
      if (10 > 1) {
        if (10 > 1) {
          return true + false;
        }
        return 1;
      }
      `,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }]`,
			"unusable as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let a = 5; let b = 12; a = b; a;", 12},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
  let newAdder = fn(x) {
    fn(y) { x + y };
  }
  let addTwo = newAdder(2);
  addTwo(2);
  `
	testIntegerObject(t, testEval(input), 4)
}

func TestLoopExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let i = 0; let num = 0; loop (i < 2) { let num = num + 1; let i = i + 1; }; num;", 2},
		{"let iterator = fn(num) { let i = 0; loop (i < 10) { let i = i + 1; num = num + 1 }; num; }; iterator(0);", 10},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: `"Hello World!"`, expected: "Hello World!"},
		{input: `"Hello \t World!"`, expected: "Hello \t World!"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (+%v)", evaluated, evaluated)
		}

		if str.Value != tt.expected {
			t.Errorf("String has wrong value. got=%q", str.Value)
		}
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (+%v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			`"Hello" == "Hello"`,
			true,
		},
		{
			`"Hello" != "hello"`,
			true,
		},
		{
			`"Hello" == "World"`,
			false,
		},
		{
			`"Hello" != "Hello"`,
			false,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		out, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Fatalf("object is not String. got=%T (+%v)", evaluated, evaluated)
		}

		val := out.Value

		if val != tt.expected {
			t.Errorf("object has wrong value. got=%t, want=%t", val, tt.expected)
		}
	}
}

func TestBuiltInFunctions(t *testing.T) {
	tests := []struct {
		expected interface{}
		input    string
	}{
		{0, `len("")`},
		{4, `len("four")`},
		{11, `len("hello world")`},
		{"argument to `len` not supported, got INTEGER", `len(1)`},
		{"wrong number of arguments. got=2, want=1", `len("one", "two")`},
		{3, `len([1, 2, 3])`},
		{0, `len([])`},
		{1, `first([1, 2, 3])`},
		{nil, `first([])`},
		{"argument to `first` must be an ARRAY, got INTEGER", `first(1)`},
		{3, `last([1, 2, 3])`},
		{"argument to `last` must be an ARRAY, got INTEGER", `last(1)`},
		{nil, `last([])`},
		{[]int{2, 3}, `rest([1, 2, 3])`},
		{nil, `rest([])`},
		{[]int{1}, `push([], 1)`},
		{"argument to `push` must be an ARRAY, got INTEGER", `push(1, 1)`},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		expected interface{}
		input    string
	}{
		{
			1,
			"[1, 2, 3][0]",
		},
		{
			2,
			"[1, 2, 3][1]",
		},
		{
			3,
			"[1, 2, 3][2]",
		},
		{
			1,
			"let i = 0; [1][i]",
		},
		{
			3,
			"[1, 2, 3][1 + 1]",
		},
		{
			3,
			"let myArray = [1, 2, 3]; myArray[2];",
		},
		{
			6,
			"let  myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2]",
		},
		{
			2,
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];",
		},
		{
			nil,
			"[1, 2, 3][3]",
		},
		{
			nil,
			"[1, 2, 3][-1]",
		},
		{
			8,
			"let myArray = [1, 2, 3]; myArray[0] = 8; myArray[0];",
		},
		{
			1,
			"let a = [1]; a[0] = a; a[0][0];",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
  {
  "one": 10 - 9,
  "two": 1 + 1,
  "thr" + "ee": 6 / 2,
  4: 4,
  true: 5,
  false: 6
  }`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. got=%d", len(result.Pairs))
	}

	for expectedK, expectedV := range expected {
		pair, ok := result.Pairs[expectedK]
		if !ok {
			t.Error("no pair given for key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedV)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		expected interface{}
		input    string
	}{
		{
			5,
			`{"foo": 5}["foo"]`,
		},
		{
			nil,
			`{"foo": 5}["bar"]`,
		},
		{
			5,
			`let key = "foo"; {"foo": 5}[key]`,
		},
		{
			5,
			`{5: 5}[5]`,
		},
		{
			5,
			`{true: 5}[true]`,
		},
		{
			5,
			`{false: 5}[false]`,
		},
		{
			5,
			`let a = {"foo": 1}; a["foo"] = 5; a["foo"]`,
		},
		{
			5,
			`let a = {}; a["foo"] = 5; a["foo"]`,
		},
		{
			23,
			`let a = {"foo": 23}; a["foo"] = a; a["foo"]["foo"]`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
