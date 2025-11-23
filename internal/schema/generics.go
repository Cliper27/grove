package schema

import (
	"fmt"

	"reflect"

	"gopkg.in/yaml.v3"
)

// NamedRule represents something with a Name and optional Rule
type NamedRule interface {
	Set(name string, rule any)
}

func cloneRule(r any) any {
	// r must be a pointer to struct
	if reflect.TypeOf(r).Kind() != reflect.Ptr {
		panic("cloneRule: ruleType must be a pointer")
	}
	return reflect.New(reflect.TypeOf(r).Elem()).Interface()
}

func unmarshalNamedRule(value *yaml.Node, target NamedRule, ruleType any) error {
	switch value.Kind {
	case yaml.ScalarNode:
		// plain string -> no rule
		target.Set(value.Value, nil)
		return nil

	case yaml.MappingNode:
		if len(value.Content) == 0 {
			return fmt.Errorf("empty mapping node")
		}

		// first key/value pair
		keyNode := value.Content[0]
		valNode := value.Content[1]

		// create a new rule object dynamically
		rule := cloneRule(ruleType)
		if err := valNode.Decode(rule); err != nil {
			return err
		}

		target.Set(keyNode.Value, rule)
		return nil

	default:
		return fmt.Errorf("unsupported node kind: %v", value.Kind)
	}
}
