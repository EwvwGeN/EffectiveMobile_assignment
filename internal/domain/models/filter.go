package models

type Filter struct {
	Fields []Field
}

type Field struct {
	UnionCondition string
	Name           string
	Value          string
	Operator       string
}