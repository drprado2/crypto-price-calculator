package observability

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
)

type (
	ContextFieldExtractor func(ctx context.Context) Attributes

	Attribute struct {
		Name  string
		Value string
	}

	Attributes []Attribute
)

func NewAttribute(name string, value interface{}) Attribute {
	if value == nil {
		return Attribute{Name: name, Value: ""}
	}

	return Attribute{Name: name, Value: fmt.Sprintf("%v", value)}
}

func NewAttributes() Attributes {
	return make([]Attribute, 0)
}

func (s Attributes) With(name string, value interface{}) Attributes {
	return append(s, NewAttribute(name, value))
}

func (s Attributes) Add(attrs ...Attribute) Attributes {
	for _, attr := range attrs {
		s = s.With(attr.Name, attr.Value)
	}

	return s
}

func (s Attributes) ToKeyValue() []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0)

	for _, attr := range s {
		if attr.Value != "" {
			attrs = append(attrs, attribute.String(attr.Name, attr.Value))
		}
	}

	return attrs
}

func (s Attributes) ToMap() map[string]interface{} {
	attrs := make(map[string]interface{}, 0)

	for _, attr := range s {
		if attr.Value != "" {
			attrs[attr.Name] = attr.Value
		}
	}

	return attrs
}
