package postgres

import (
	"fmt"
	"strings"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"
)

func AddFilter(filter *models.Filter, name, value string) error {
	operator := "="
	field := models.Field{
		UnionCondition: "AND",
		Name: name,
		Operator: operator,
		Value: value,
	}
	splited := strings.Split(value, ":")
	if len(splited) == 1 {
		filter.Fields = append(filter.Fields, field)
		return nil
	}
	// how to distinguish a missing value from an empty one?
	if len(splited) == 2 && splited[1] == "" {
		return fmt.Errorf("empty value of filters field")
	}
	operatorPos := 0
	switch strings.ToLower(splited[0]) {
	case "and":
		operatorPos++
	case "or":
		field.UnionCondition = "OR"
		operatorPos++
	default:
	}
	switch strings.ToLower(splited[operatorPos]) {
	case "eq":
	case "neq":
		field.Operator = "<>"
	case "gt":
		field.Operator = ">"
	case "get":
		field.Operator = ">="
	case "lt":
		field.Operator = "<"
	case "let":
		field.Operator = "<="
	case "like":
		field.Operator = "LIKE"
	default:
		return fmt.Errorf("not valid filter operator")			
	}
	field.Value = strings.Join(splited[operatorPos+1:], "")
	filter.Fields = append(filter.Fields, field)
	return nil
}