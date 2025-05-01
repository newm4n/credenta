package credenta

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	IdTypeUserId      IdType = "USERID"
	IdTypeUserEmail   IdType = "EMAIL"
	IdTypeUserPhoneNo IdType = "PHONENO"
)

type IdType string

type CUser struct {
	FilePath string `json:"-"`

	Realm      string                `json:"realm"`
	Id         string                `json:"id"`
	IDType     IdType                `json:"idType"`
	Groups     []string              `json:"groups,omitempty"`
	Attributes map[string]*Attribute `json:"attributes"`
	RoleMasks  []uint64              `json:"roleMasks"`

	VerificationMethod VerificationMethod `json:"method"`
	VerificationHash   string             `json:"hash"`

	Enable bool `json:"enable"`
	Active bool `json:"active"`

	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

func (user *CUser) StoreOrSaveToFile(ctx context.Context) error {
	user.UpdatedBy = ctx.Value(ETX_USER).(string)
	user.UpdatedAt = time.Now()

	data, err := json.Marshal(user)

	if err != nil {
		return fmt.Errorf("in StoreOrSaveToFile function, error marshalling user: %w", err)
	}
	if _, err := os.Stat(user.FilePath); err == nil {
		f, err := os.Open(user.FilePath)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error opening file %s: %w", user.FilePath, err)
		}
		defer f.Close()
		err = f.Truncate(0)
		if err != nil {
			return fmt.Errorf("in writeDataToFile function. error truncate file %s: %w", user.FilePath, err)
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error seek in file %s: %w", user.FilePath, err)
		}
		_, err = f.Write(data)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error writing into file %s: %w", user.FilePath, err)
		}
	} else if os.IsNotExist(err) {
		file, err := os.Create(user.FilePath)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error creating file %s: %w", user.FilePath, err)
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf("in StoreOrSaveToFile function. error writing data to file %s: %w", user.FilePath, err)
		}
	} else {
		return fmt.Errorf("in StoreOrSaveToFile function. error obtaining stat of file %s: %w", user.FilePath, err)
	}
	return nil
}

func (user *CUser) ReloadFromFile(ctx context.Context) error {
	file, err := os.Open(user.FilePath)
	if err != nil {
		return fmt.Errorf("in ReloadFromFile function, error opening file %s: %w", user.FilePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buff := bytes.Buffer{}

	for scanner.Scan() {
		buff.Write(scanner.Bytes())
	}

	nUser := &CUser{}

	err = json.Unmarshal(buff.Bytes(), &nUser)
	if err != nil {
		return fmt.Errorf("in ReloadFromFile function, error unmarshaling data into CUser: %w", err)
	}

	user.Realm = nUser.Realm
	user.Id = nUser.Id
	user.IDType = nUser.IDType
	user.Groups = nUser.Groups
	user.Attributes = nUser.Attributes
	user.RoleMasks = nUser.RoleMasks

	user.VerificationMethod = nUser.VerificationMethod
	user.VerificationHash = nUser.VerificationHash

	user.Enable = nUser.Enable
	user.Active = nUser.Active

	user.CreatedAt = nUser.CreatedAt
	user.CreatedBy = nUser.CreatedBy
	user.UpdatedAt = nUser.UpdatedAt
	user.UpdatedBy = nUser.UpdatedBy

	return nil
}

func (user *CUser) DeleteFile(ctx context.Context) error {
	return os.Remove(user.FilePath)
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
	copyAttribute := make([]*Attribute, len(user.Attributes))
	i := 0
	for k := range user.Attributes {
		copyAttribute[i] = user.Attributes[k]
	}
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
	fmt.Printf("Looking for %s in %d\n", name, len(user.Attributes))
	for k, _ := range user.Attributes {
		fmt.Printf("Found %s\n", k)
		if strings.EqualFold(k, name) {
			return true
		}
	}
	return false
}
func (user *CUser) RemoveAttribute(name string) {
	if user.Attributes != nil {
		for key, _ := range user.Attributes {
			if strings.EqualFold(key, name) {
				delete(user.Attributes, key)
			}
		}
		for i, attr := range user.SortAttributes() {
			attr.Seq = i
		}
	}
}
func (user *CUser) RemoveAllAttributes() {
	user.Attributes = make(map[string]*Attribute)
}
func (user *CUser) GetsAttribute(name string) (string, error) {
	if user.Attributes != nil {
		for _, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.StringValue, nil
			}
		}
	}
	return "", fmt.Errorf("attribute not found")
}
func (user *CUser) GetiAttribute(name string) (int, error) {
	if user.Attributes != nil {
		for _, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.IntegerValue, nil
			}
		}
	}
	return -1, fmt.Errorf("attribute not found")
}
func (user *CUser) GetfAttribute(name string) (float64, error) {
	if user.Attributes != nil {
		for _, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.FloatValue, nil
			}
		}
	}
	return -1, fmt.Errorf("attribute not found")
}
func (user *CUser) GetbAttribute(name string) (bool, error) {
	if user.Attributes != nil {
		for _, attr := range user.Attributes {
			if strings.EqualFold(attr.Name, name) {
				return attr.BoolValue, nil
			}
		}
	}
	return false, fmt.Errorf("attribute not found")
}

func (user *CUser) SetsAttribute(name, value string) error {
	if user.Attributes == nil {
		user.Attributes = make(map[string]*Attribute)
	}
	if user.HasAttribute(name) {
		return fmt.Errorf("attribute already exists")
	}
	user.Attributes[name] = &Attribute{
		Name:         name,
		Seq:          len(user.Attributes),
		StringValue:  value,
		IntegerValue: 0,
		FloatValue:   0,
		BoolValue:    false,
	}
	return nil
}
func (user *CUser) SetiAttribute(name string, value int) error {
	if user.Attributes == nil {
		user.Attributes = make(map[string]*Attribute)
	}
	if user.HasAttribute(name) {
		return fmt.Errorf("attribute already exists")
	}
	user.Attributes[name] = &Attribute{
		Name:         name,
		Seq:          len(user.Attributes),
		StringValue:  "",
		IntegerValue: value,
		FloatValue:   0,
		BoolValue:    false,
	}
	return nil
}
func (user *CUser) SetfAttribute(name string, value float64) error {
	if user.Attributes == nil {
		user.Attributes = make(map[string]*Attribute)
	}
	if user.HasAttribute(name) {
		return fmt.Errorf("attribute already exists")
	}
	user.Attributes[name] = &Attribute{
		Name:         name,
		Seq:          len(user.Attributes),
		StringValue:  "",
		IntegerValue: 0,
		FloatValue:   value,
		BoolValue:    false,
	}
	return nil
}
func (user *CUser) SetbAttribute(name string, value bool) error {
	if user.Attributes == nil {
		user.Attributes = make(map[string]*Attribute)
	}
	if user.HasAttribute(name) {
		return fmt.Errorf("attribute already exists")
	}
	user.Attributes[name] = &Attribute{
		Name:         name,
		Seq:          len(user.Attributes),
		StringValue:  "",
		IntegerValue: 0,
		FloatValue:   0,
		BoolValue:    value,
	}
	return nil
}

func (user *CUser) String() string {
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Sprintf("error %v", err)
	}
	return string(jsonBytes)
}
