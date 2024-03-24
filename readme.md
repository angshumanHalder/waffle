# About

Waffle is a dynamic programming language. With minimal data types and syntax.

# Installation 
1. Clone the repo.
2. Download and install go compiler.
3. cd into the repo and run `go build main.go`
4. Once the build process completes just run the binary.

# Syntax

### Declarations 
Waffle doesn't support constants. Every variable is declared with let keyword.
```
let age = 1;
let name = "Tuna";
puts(age); // 1
puts(name); // Tuna

name = "Bob";
puts(name); // Bob

```
### Arrays
Arrays are just collection of values (values can be of any type).
```
let arr = [1, 2, "Hello World"]
puts(arr); // [1, 2, "Hello World"]
puts(arr[0]); // 1
```

### Strings
Everything inside `""` is considered a string.
```
let name = "Bob";
puts(name); // Bob
```


### Booleans
```
puts(true); // true
puts(false): // false

```
### Numbers
```
let a = 1;
let b = 2;
puts(a / b) // 0

let a = 1.1;
let b = 2;
puts(a / b); // 0.55
```

### Objects
Objects supports integers, booleans and strings as keys.
```
let myHash = {"name": "Jimmy", "age": 72, "band": "Led Zeppelin", 99: "integer", true: "true"};
puts(myHash["name"]); // Jimmy

let key = "age";
myHash[key] = 32;
puts(myHash[key]); // 32
```

### Equality
```
puts(1 == 1); // true
puts(1 == 2); // false
```

### Conditionals
Conditionals works the same way as they do in other programming languages.
`if else` is not supported rather the else block can have as many if blocks inside it.
```
let x = 2;
if (x > 10) { 
  puts("everything okay!");
} else {
  if (x < 5) {
    puts("x is too low!"); 
  } else {
    puts("x is low!")
  }
}
```

### Functions
Functions are first class functions in Waffle.
```
let func = fn(x, y) {
  return x + y;
}
puts(func(1, 2)) // 3

fn() {
  puts("hello world")
}() // "hello world"

let function = fn() {
  return fn(x, y) { return x > y; };
}

let checkGreater = function();
checkGreater(1, 2); // false

```

### Builtin functions
Waffle has some basic builtin functions.
```
let a = [1, 2, 3, 4];
puts(len(a)); // 4
puts(first(a)); // 1
puts(last(a)); // 4
puts(rest(a)); // [2, 3, 4]

let a = push(arr, 5);
puts(arr); // [1, 2, 3, 4, 5]

let name = "Bob";
puts(len(name)); // 3

```
