package models

type Filter struct {
	fields []field
}

type field struct {
	Name     string
	Value    string
	Operator string
}

func (f *Filter) GetFields() []field {
	return f.fields
}