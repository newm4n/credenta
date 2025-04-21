package credenta

import "fmt"

type Attributable interface {
	GetAttributeList() []string
	HasAttribute(name string) bool
	RemoveAttribute(name string)
	RemoveAllAttributes()

	GetsAttribute(name string) (string, error)
	GetiAttribute(name string) (int, error)
	GetfAttribute(name string) (float64, error)
	GetbAttribute(name string) (bool, error)

	SetsAttribute(name, value string) error
	SetiAttribute(name string, value int) error
	SetfAttribute(name string, value float64) error
	SetbAttribute(name string, value bool) error
}

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
