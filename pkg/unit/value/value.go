package value

type Type string

const (
	TypeNone   Type = "none"
	TypeString      = "string"
	TypeInt         = "int"
	TypeFloat       = "float"
	TypeBool        = "bool"
)

type Value struct {
	typ   Type
	value interface{}
}

type Input <-chan Value
type Output chan<- Value

func (v Value) Null() bool {
	return v.typ == ""
}

func None() Value {
	return Value{
		typ: TypeNone,
	}
}

func Bool(v bool) Value {
	return Value{
		typ:   TypeBool,
		value: v,
	}
}

func String(v string) Value {
	return Value{
		typ:   TypeString,
		value: v,
	}
}

func Float(v float64) Value {
	return Value{
		typ:   TypeFloat,
		value: v,
	}
}

func Int(v int64) Value {
	return Value{
		typ:   TypeInt,
		value: v,
	}
}

func (v Value) Type() Type {
	return v.typ
}

func (v Value) Bool() bool {
	return v.value.(bool)
}

func (v Value) String() string {
	return v.value.(string)
}

func (v Value) Float() float64 {
	return v.value.(float64)
}

func (v Value) Int() int64 {
	return v.value.(int64)
}
