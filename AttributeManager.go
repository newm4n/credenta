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

	// GetsAttribute retrieve string value of attribute with specified name, or error if problem during retrieval.
	GetsAttribute(name string) (string, error)
	// GetiAttribute retrieve int value of attribute with specified name, or error if problem during retrieval.
	GetiAttribute(name string) (int, error)
	// GetfAttribute retrieve float64 value of attribute with specified name, or error if problem during retrieval.
	GetfAttribute(name string) (float64, error)
	// GetbAttribute retrieve bool value of attribute with specified name, or error if problem during retrieval.
	GetbAttribute(name string) (bool, error)

	// SetsAttribute set string value attribute with specified name return error if problem during storing.
	SetsAttribute(name, value string) error
	// SetiAttribute set int value attribute with specified name return error if problem during storing.
	SetiAttribute(name string, value int) error
	// SetfAttribute set float64 value attribute with specified name return error if problem during storing.
	SetfAttribute(name string, value float64) error
	// SetbAttribute set bool value attribute with specified name return error if problem during storing.
	SetbAttribute(name string, value bool) error
}

// Attribute the attribute element that can be contained within an object that implements Attributable
type Attribute struct {
	Name         string  `json:"name"`
	Seq          int     `json:"seq,omitempty"`
	StringValue  string  `json:"stringValue,omitempty"`
	IntegerValue int     `json:"integerValue,omitempty"`
	FloatValue   float64 `json:"floatValue,omitempty"`
	BoolValue    bool    `json:"boolValue,omitempty"`
}

func (attr *Attribute) String() string {
	if attr.StringValue != "" {
		return fmt.Sprintf("%s(%s)", attr.Name, attr.StringValue)
	}
	if attr.IntegerValue != 0 {
		return fmt.Sprintf("%s(%d)", attr.Name, attr.IntegerValue)
	}
	if attr.FloatValue != 0 {
		return fmt.Sprintf("%s(%f)", attr.Name, attr.FloatValue)
	}
	return fmt.Sprintf("%s(%v)", attr.Name, attr.BoolValue)
}
