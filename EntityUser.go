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
	"time"
)

const (
	// IdTypeUserId a simple user ID.
	IdTypeUserId IdType = "USERID"
	// IdTypeUserEmail an ID from email address
	IdTypeUserEmail IdType = "EMAIL"
	// IdTypeUserPhoneNo an ID from a phone number
	IdTypeUserPhoneNo IdType = "PHONENO"
)

// IdType specify the type of Identification.
// this is useful to know if some kind of validation is required for the system for user.
// validation such as email box ownership or phone number ownership.
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

func (user *CUser) SortAttributeKeys() []string {
	if user.Attributes == nil {
		return make([]string, 0)
	}
	attnames := user.GetAttributeList()
	sort.Strings(attnames)
	return attnames
}

func (user *CUser) GetAttributeList() []string {
	if user.Attributes == nil {
		return make([]string, 0)
	}

	names := make([]string, 0)
	for k, _ := range user.Attributes {
		names = append(names, k)
	}
	return names
}
func (user *CUser) HasAttribute(name string) bool {
	if user.Attributes == nil {
		return false
	}
	_, exist := user.Attributes[name]
	return exist
}
func (user *CUser) RemoveAttribute(name string) {
	if user.Attributes != nil {
		delete(user.Attributes, name)
	}
}
func (user *CUser) RemoveAllAttributes() {
	user.Attributes = make(map[string]*Attribute)
}

// GetAttribute retrieve value of attribute with specified name, or error if problem during retrieval.
func (user *CUser) GetAttribute(name string) (valueType, valueString string, err error) {
	if user.Attributes != nil {
		if attrib, exist := user.Attributes[name]; exist {
			return attrib.ValueType, attrib.ValueString, nil
		}
	}
	return "", "", errors.New("attribute not found")
}

// SetAttribute set value attribute with specified name, type and the value in a string representation.
// return error if problem during storing.
func (user *CUser) SetAttribute(attributeName, valueType, valueString string) error {
	if user.Attributes == nil {
		user.Attributes = make(map[string]*Attribute)
	}
	if _, exist := user.Attributes[attributeName]; !exist {
		user.Attributes[attributeName] = &Attribute{
			Name:        attributeName,
			Seq:         len(user.Attributes),
			ValueType:   valueType,
			ValueString: valueString,
		}
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
