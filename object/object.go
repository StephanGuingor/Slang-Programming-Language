package object

import (
	"bytes"
	"compiler-book/ast"
	"fmt"
	"strings"
)

const (
	INTEGER  ObjectType = "INTEGER"
	FLOAT    ObjectType = "FLOAT"
	STRING   ObjectType = "STRING"
	RUNE     ObjectType = "RUNE"
	BOOLEAN  ObjectType = "BOOLEAN"
	NULL     ObjectType = "NULL"
	FUNCTION ObjectType = "FUNCTION"

	RETURN_VALUE ObjectType = "RETURN_VALUE"
	ERROR        ObjectType = "ERROR"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType {
	return FLOAT
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING
}

func (s *String) Inspect() string {
	return s.Value
}

type Rune struct {
	Value rune
}

func (r *Rune) Type() ObjectType {
	return RUNE
}

func (r *Rune) Inspect() string {
	return fmt.Sprintf("%c", r.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// TODO: improve trace
type Trace struct {
	Column int
	Line   int
}

type Error struct {
	Message string
	Trace   Trace
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
