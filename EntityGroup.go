package credenta

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type CGroup struct {
	FilePath string `json:"-"`

	Realm        string       `json:"realm"`
	Name         string       `json:"name"`
	ParentGroups []string     `json:"parentGroups,omitempty"`
	Attributes   []*Attribute `json:"attributes,omitempty"`
	RoleMasks    []uint64     `json:"roleMasks,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

func (group *CGroup) StoreOrSaveToFile(ctx context.Context) error {
	group.UpdatedBy = ctx.Value(ETX_USER).(string)
	group.UpdatedAt = time.Now()

	/*
		TODO Make sure that the group to parentGroup relation will not cause a cyclic relation. If it form a cyclic, it should return an error message
	*/

	data, err := json.Marshal(group)
	if err != nil {
		return fmt.Errorf("in StoreOrSaveToFile function, error marshalling group: %w", err)
	}
	if _, err := os.Stat(group.FilePath); err == nil {
		f, err := os.Open(group.FilePath)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error opening file %s: %w", group.FilePath, err)
		}
		defer f.Close()
		err = f.Truncate(0)
		if err != nil {
			return fmt.Errorf("in writeDataToFile function. error truncate file %s: %w", group.FilePath, err)
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error seek in file %s: %w", group.FilePath, err)
		}
		_, err = f.Write(data)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error writing into file %s: %w", group.FilePath, err)
		}
	} else if os.IsNotExist(err) {
		file, err := os.Create(group.FilePath)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error creating file %s: %w", group.FilePath, err)
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error writing data to file %s: %w", group.FilePath, err)
		}
	} else {
		return fmt.Errorf("in StoreOrSaveToFile function. error obtaining stat of file %s: %w", group.FilePath, err)
	}
	return nil
}

func (group *CGroup) ReloadFromFile(ctx context.Context) error {
	file, err := os.Open(group.FilePath)
	if err != nil {
		return fmt.Errorf("in ReloadFromFile function, error opening file %s: %w", group.FilePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buff := bytes.Buffer{}

	for scanner.Scan() {
		buff.Write(scanner.Bytes())
	}

	nGroup := &CGroup{}

	err = json.Unmarshal(buff.Bytes(), &nGroup)
	if err != nil {
		return fmt.Errorf("in ReloadFromFile function, error unmarshaling data into CGroup: %w", err)
	}

	group.Realm = nGroup.Realm
	group.Name = nGroup.Name
	group.ParentGroups = nGroup.ParentGroups
	group.Attributes = nGroup.Attributes
	group.RoleMasks = nGroup.RoleMasks

	group.CreatedAt = nGroup.CreatedAt
	group.CreatedBy = nGroup.CreatedBy
	group.UpdatedAt = nGroup.UpdatedAt
	group.UpdatedBy = nGroup.UpdatedBy

	return nil
}

func (group *CGroup) DeleteFile(ctx context.Context) error {
	return os.Remove(group.FilePath)
}

func (group *CGroup) AddRole(roleSquence int) {
	seq, bit := toUint64ByBit(roleSquence)
	group.RoleMasks[seq] = setBitFlagOn(group.RoleMasks[seq], bit)
}

func (group *CGroup) RemoveRole(roleSquence int) {
	seq, bit := toUint64ByBit(roleSquence)
	group.RoleMasks[seq] = setBitFlagOff(group.RoleMasks[seq], bit)
}

func (group *CGroup) HasRole(roleSquence int) bool {
	seq, bit := toUint64ByBit(roleSquence)
	return isBitFlagOn(group.RoleMasks[seq], bit)
}

func (group *CGroup) ClearRole() {
	for i := 0; i < len(group.RoleMasks); i++ {
		group.RoleMasks[i] = 0
	}
}

func (group *CGroup) GetAttributeList() []string {
	names := make([]string, len(group.Attributes))
	for i, attr := range group.SortAttributes() {
		names[i] = attr.Name
	}
	return names
}

func (group *CGroup) HasAttribute(name string) bool {
	for _, attr := range group.Attributes {
		if strings.EqualFold(attr.Name, name) {
			return true
		}
	}
	return false
}

func (group *CGroup) RemoveAttribute(name string) {
	for i, attr := range group.Attributes {
		if strings.EqualFold(attr.Name, name) {
			group.Attributes = append(group.Attributes[:i], group.Attributes[i+1:]...)
		}
	}
	for i, attr := range group.SortAttributes() {
		attr.Seq = i
	}
}

func (group *CGroup) RemoveAllAttributes() {
	group.Attributes = group.Attributes[:0]
}

// GetAttribute retrieve value of attribute with specified name, or error if problem during retrieval.
func (group *CGroup) GetAttribute(name string) (valueType, valueString string, err error) {
	if group.Attributes != nil {
		for _, attr := range group.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.ValueType, attr.ValueString, nil
			}
		}
	}
	return "", "", errors.New("attribute not found")
}

// SetAttribute set value attribute with specified name, type and the value in a string representation.
// return error if problem during storing.
func (group *CGroup) SetAttribute(attributeName, valueType, valueString string) error {
	if group.Attributes == nil {
		group.Attributes = make([]*Attribute, 0)
	}
	if group.HasAttribute(attributeName) {
		return fmt.Errorf("attribute already exists")
	}
	group.Attributes = append(group.Attributes, &Attribute{
		Name:        attributeName,
		Seq:         len(group.Attributes),
		ValueType:   valueType,
		ValueString: valueString,
	})
	return nil
}

func (group *CGroup) SortAttributes() []*Attribute {
	copy := group.Attributes
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Seq < copy[j].Seq
	})
	return copy
}
