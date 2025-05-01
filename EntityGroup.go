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

	Realm        string       `json:"realm" :"realm"`
	Name         string       `json:"name" :"name"`
	ParentGroups []string     `json:"parentGroups,omitempty" :"parentGroups"`
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
