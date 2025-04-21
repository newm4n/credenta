package credenta

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	IdTypeUserId      IdType = "USERID"
	IdTypeUserEmail   IdType = "EMAIL"
	IdTypeUserPhoneNo IdType = "PHONENO"
)

type IdType string

type CUser struct {
	Realm      string       `json:"realm"`
	Id         string       `json:"id"`
	IDType     IdType       `json:"idType"`
	Groups     []string     `json:"groups,omitempty"`
	Attributes []*Attribute `json:"attributes"`
	RoleMasks  []uint64     `json:"roleMasks"`

	VerificationMethod VerificationMethod `json:"method"`
	VerificationHash   string             `json:"hash"`

	Enable bool `json:"enable"`
	Active bool `json:"active"`
}

func toUint64ByBit(roleSquence int) (uintSeq int, bitno int) {
	left := math.Mod(float64(roleSquence), 64.0)
	div := roleSquence / 64
	return div, int(left)
}

func isBitFlagOn(currentBit uint64, bitSequence int) bool {
	flipper := uint64(1) << bitSequence
	return flipper == currentBit|flipper
}

func setBitFlagOn(currentBit uint64, bitSequence int) uint64 {
	flipper := uint64(1) << bitSequence
	return currentBit | flipper
}

func setBitFlagOff(currentBit uint64, bitSequence int) uint64 {
	flipper := uint64(1) << bitSequence
	notFlipper := 0xFFFFFFFF ^ flipper
	return currentBit & notFlipper
}

func (user *CUser) AddRole(roleSquence int) {
	seq, bit := toUint64ByBit(roleSquence)
	user.RoleMasks[seq] = setBitFlagOn(user.RoleMasks[seq], bit)
}

func (user *CUser) RemoveRole(roleSquence int) {
	seq, bit := toUint64ByBit(roleSquence)
	user.RoleMasks[seq] = setBitFlagOff(user.RoleMasks[seq], bit)
}

func (user *CUser) HasRole(roleSquence int) bool {
	seq, bit := toUint64ByBit(roleSquence)
	return isBitFlagOn(user.RoleMasks[seq], bit)
}

func (user *CUser) ClearRole() {
	for i := 0; i < len(user.RoleMasks); i++ {
		user.RoleMasks[i] = 0
	}
}

func (user *CUser) SortAttributes() []*Attribute {
	if user.Attributes == nil {
		return make([]*Attribute, 0)
	}
	copyAttribute := user.Attributes
	sort.Slice(copyAttribute, func(i, j int) bool {
		return copyAttribute[i].Seq < copyAttribute[j].Seq
	})
	return copyAttribute
}

func (user *CUser) GetAttributeList() []string {
	if user.Attributes == nil {
		return make([]string, 0)
	}

	names := make([]string, len(user.Attributes))
	for i, attr := range user.SortAttributes() {
		names[i] = attr.Name
	}
	return names
}
func (user *CUser) HasAttribute(name string) bool {
	if user.Attributes == nil {
		return false
	}
	for _, attr := range user.Attributes {
		if strings.EqualFold(attr.Name, name) {
			return true
		}
	}
	return false
}
func (user *CUser) RemoveAttribute(name string) {
	if user.Attributes != nil {
		for i, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				user.Attributes = append(user.Attributes[:i], user.Attributes[i+1:]...)
			}
		}
		for i, attr := range user.SortAttributes() {
			attr.Seq = i
		}
	}
}
func (user *CUser) RemoveAllAttributes() {
	if user.Attributes != nil {
		user.Attributes = user.Attributes[:0]
	}
}
func (user *CUser) GetsAttribute(name string) (string, error) {
	if user.Attributes != nil {
		for _, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.StringValue, nil
			}
		}
	}
	return "", errors.New("attribute not found")
}
func (user *CUser) GetiAttribute(name string) (int, error) {
	if user.Attributes != nil {
		for _, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.IntegerValue, nil
			}
		}
	}
	return -1, errors.New("attribute not found")
}
func (user *CUser) GetfAttribute(name string) (float64, error) {
	if user.Attributes != nil {
		for _, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.FloatValue, nil
			}
		}
	}
	return -1, errors.New("attribute not found")
}
func (user *CUser) GetbAttribute(name string) (bool, error) {
	if user.Attributes != nil {
		for _, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.BoolValue, nil
			}
		}
	}
	return false, errors.New("attribute not found")
}

func (user *CUser) SetsAttribute(name, value string) error {
	if user.Attributes != nil {
		user.Attributes = make([]*Attribute, 0)
	}
	if user.HasAttribute(name) {
		return errors.New("attribute already exists")
	}
	user.Attributes = append(user.Attributes, &Attribute{
		Name:         name,
		Seq:          len(user.Attributes),
		StringValue:  value,
		IntegerValue: 0,
		FloatValue:   0,
		BoolValue:    false,
	})
	return nil
}
func (user *CUser) SetiAttribute(name string, value int) error {
	if user.Attributes != nil {
		user.Attributes = make([]*Attribute, 0)
	}
	if user.HasAttribute(name) {
		fmt.Print("NOK")
		return errors.New("attribute already exists")
	}
	fmt.Print("OK")
	user.Attributes = append(user.Attributes, &Attribute{
		Name:         name,
		Seq:          len(user.Attributes),
		StringValue:  "",
		IntegerValue: value,
		FloatValue:   0,
		BoolValue:    false,
	})
	return nil
}
func (user *CUser) SetfAttribute(name string, value float64) error {
	if user.Attributes != nil {
		user.Attributes = make([]*Attribute, 0)
	}
	if user.HasAttribute(name) {
		return errors.New("attribute already exists")
	}
	user.Attributes = append(user.Attributes, &Attribute{
		Name:         name,
		Seq:          len(user.Attributes),
		StringValue:  "",
		IntegerValue: 0,
		FloatValue:   value,
		BoolValue:    false,
	})
	return nil
}
func (user *CUser) SetbAttribute(name string, value bool) error {
	if user.Attributes != nil {
		user.Attributes = make([]*Attribute, 0)
	}
	if user.HasAttribute(name) {
		return errors.New("attribute already exists")
	}
	user.Attributes = append(user.Attributes, &Attribute{
		Name:         name,
		Seq:          len(user.Attributes),
		StringValue:  "",
		IntegerValue: 0,
		FloatValue:   0,
		BoolValue:    value,
	})
	return nil
}

func (user *CUser) String() string {
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Sprintf("error %v", err)
	}
	return string(jsonBytes)
}
