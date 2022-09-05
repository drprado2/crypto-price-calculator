package observability

import (
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"testing"
)

func TestNewAttribute(t *testing.T) {
	assert.NotPanics(t, func() {
		NewAttribute("test", nil)
	})

	attr := NewAttribute("test", "")

	assert.Equal(t, attr.Name, "test")
	assert.Equal(t, attr.Value, "")
}

func TestNewAttributes(t *testing.T) {
	attrs := NewAttributes().
		With("Cid", "1").
		With("AnyValue", "2").
		With("ConfigId", "3")

	assert.Len(t, attrs, 3, "should have added 3 attributes")
	assert.Contains(t, attrs, Attribute{Name: "Cid", Value: "1"})
	assert.Contains(t, attrs, Attribute{Name: "AnyValue", Value: "2"})
	assert.Contains(t, attrs, Attribute{Name: "ConfigId", Value: "3"})
}

func TestToKeyValues(t *testing.T) {
	attrs := NewAttributes().
		With("Cid", "1").
		With("AnyValue", "2").
		With("ConfigId", "3").
		With("Test", nil)

	assert.Len(t, attrs, 4, "should have added 3 attributes")
	assert.Contains(t, attrs, Attribute{Name: "Cid", Value: "1"})
	assert.Contains(t, attrs, Attribute{Name: "AnyValue", Value: "2"})
	assert.Contains(t, attrs, Attribute{Name: "ConfigId", Value: "3"})
	assert.Contains(t, attrs, Attribute{Name: "Test", Value: ""})

	attrsKeyValue := attrs.ToKeyValue()

	assert.Len(t, attrsKeyValue, 3, "should have added 3 attributes")
	assert.Contains(t, attrsKeyValue, attribute.String("Cid", "1"))
	assert.Contains(t, attrsKeyValue, attribute.String("AnyValue", "2"))
	assert.Contains(t, attrsKeyValue, attribute.String("ConfigId", "3"))
	assert.NotContains(t, attrsKeyValue, attribute.String("Test", ""))
}

func TestToMap(t *testing.T) {
	attrs := NewAttributes().
		With("Cid", "1").
		With("AnyValue", "2").
		With("ConfigId", "3").
		With("Test", nil)

	assert.Len(t, attrs, 4, "should have added 3 attributes")
	assert.Contains(t, attrs, Attribute{Name: "Cid", Value: "1"})
	assert.Contains(t, attrs, Attribute{Name: "AnyValue", Value: "2"})
	assert.Contains(t, attrs, Attribute{Name: "ConfigId", Value: "3"})
	assert.Contains(t, attrs, Attribute{Name: "Test", Value: ""})

	attrsKeyValue := attrs.ToMap()

	assert.Len(t, attrsKeyValue, 3, "should have removed empty values")
	assert.Equal(t, attrsKeyValue, map[string]interface{}{"Cid": "1", "ConfigId": "3", "AnyValue": "2"})
}
