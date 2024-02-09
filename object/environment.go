package object

import "bytes"

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

// Debug purpose
func (e *Environment) String() string {
	var out bytes.Buffer

	out.WriteString("{\n")
	for k, v := range e.store {
		out.WriteString(k + ": " + v.Inspect() + ",\n")
	}
	out.WriteString("}")

	return out.String()
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
