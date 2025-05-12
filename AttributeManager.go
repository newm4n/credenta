package credenta

import "fmt"

// Attributable defines methods for all object that can contain an attribute.
type Attributable interface {
	// GetAttributeList returns list of attribute names contained in the implementing object
	GetAttributeList() []string
	// HasAttribute returns boolean for checking if the object have an attribute with specified name
	HasAttribute(name string) bool
	// RemoveAttribute remove an attribute with specified name form list of attributes in the object
	RemoveAttribute(name string)
	// RemoveAllAttributes remove all attributes in the object
	RemoveAllAttributes()

	// GetAttribute retrieve value of attribute with specified name, or error if problem during retrieval.
	GetAttribute(name string) (valueType, valueString string, err error)

	// SetAttribute set value attribute with specified name, type and the value in a string representation.
	// return error if problem during storing.
	SetAttribute(attributeName, valueType, valueString string) error
}

// Attribute the attribute element that can be contained within an object that implements Attributable
type Attribute struct {
	Name        string `json:"name"`
	Seq         int    `json:"seq,omitempty"`
	ValueType   string `json:"valueType"`
	ValueString string `json:"valueString"`
}

func (attr *Attribute) String() string {
	return fmt.Sprintf("%s = %s(%s)", attr.Name, attr.ValueString, attr.ValueType)
}
