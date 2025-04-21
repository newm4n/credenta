package credenta

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type CGroup struct {
	Realm        string       `json:"realm" :"realm"`
	Name         string       `json:"name" :"name"`
	ParentGroups []string     `json:"parentGroups" :"parentGroups"`
	Attributes   []*Attribute `json:"attributes"`
}

func (grp *CGroup) GetAttributeList() []string {
	names := make([]string, len(grp.Attributes))
	for i, attr := range grp.SortAttributes() {
		names[i] = attr.Name
	}
	return names
}

func (grp *CGroup) HasAttribute(name string) bool {
	for _, attr := range grp.Attributes {
		if strings.EqualFold(attr.Name, name) {
			return true
		}
	}
	return false
}

func (grp *CGroup) RemoveAttribute(name string) {
	for i, attr := range grp.Attributes {
		if strings.EqualFold(attr.Name, name) {
			grp.Attributes = append(grp.Attributes[:i], grp.Attributes[i+1:]...)
		}
	}
	for i, attr := range grp.SortAttributes() {
		attr.Seq = i
	}
}

func (grp *CGroup) RemoveAllAttributes() {
	grp.Attributes = grp.Attributes[:0]
}

func (grp *CGroup) GetsAttribute(name string) (string, error) {
	for _, attr := range grp.Attributes {
		if strings.EqualFold(attr.Name, name) {
			return attr.StringValue, nil
		}
	}
	return "", errors.New("attribute not found")
}

func (grp *CGroup) GetiAttribute(name string) (int, error) {
	for _, attr := range grp.Attributes {
		if strings.EqualFold(attr.Name, name) {
			return attr.IntegerValue, nil
		}
	}
	return -1, errors.New("attribute not found")
}
func (grp *CGroup) GetfAttribute(name string) (float64, error) {
	for _, attr := range grp.Attributes {
		if strings.EqualFold(attr.Name, name) {
			return attr.FloatValue, nil
		}
	}
	return -1, errors.New("attribute not found")
}

func (grp *CGroup) GetbAttribute(name string) (bool, error) {
	for _, attr := range grp.Attributes {
		if strings.EqualFold(attr.Name, name) {
			return attr.BoolValue, nil
		}
	}
	return false, errors.New("attribute not found")
}

func (grp *CGroup) SetsAttribute(name, value string) error {
	if grp.HasAttribute(name) {
		return errors.New("attribute already exists")
	}
	grp.Attributes = append(grp.Attributes, &Attribute{
		Name:         name,
		Seq:          len(grp.Attributes),
		StringValue:  value,
		IntegerValue: 0,
		FloatValue:   0,
		BoolValue:    false,
	})
	return nil
}

func (grp *CGroup) SetiAttribute(name string, value int) error {
	if grp.HasAttribute(name) {
		fmt.Print("NOK")
		return errors.New("attribute already exists")
	}
	fmt.Print("OK")
	grp.Attributes = append(grp.Attributes, &Attribute{
		Name:         name,
		Seq:          len(grp.Attributes),
		StringValue:  "",
		IntegerValue: value,
		FloatValue:   0,
		BoolValue:    false,
	})
	return nil
}

func (grp *CGroup) SetfAttribute(name string, value float64) error {
	if grp.HasAttribute(name) {
		return errors.New("attribute already exists")
	}
	grp.Attributes = append(grp.Attributes, &Attribute{
		Name:         name,
		Seq:          len(grp.Attributes),
		StringValue:  "",
		IntegerValue: 0,
		FloatValue:   value,
		BoolValue:    false,
	})
	return nil
}

func (grp *CGroup) SetbAttribute(name string, value bool) error {
	if grp.HasAttribute(name) {
		return errors.New("attribute already exists")
	}
	grp.Attributes = append(grp.Attributes, &Attribute{
		Name:         name,
		Seq:          len(grp.Attributes),
		StringValue:  "",
		IntegerValue: 0,
		FloatValue:   0,
		BoolValue:    value,
	})
	return nil
}

func (grp *CGroup) SortAttributes() []*Attribute {
	copy := grp.Attributes
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Seq < copy[j].Seq
	})
	return copy
}
