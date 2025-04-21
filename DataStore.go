package credenta

import (
	"errors"
	"fmt"
)

func NewCredentaDB() *CredentaDB {
	return &CredentaDB{
		DefaultRealm: "DEFAULT",
		PassPolicy: &PassphrasePolicy{
			WordCount:               1,
			LetterCountPerWord:      8,
			LetterCountMinimumTotal: 8,
			MustHaveUpperAlphabet:   false,
			MustHaveNumeric:         false,
			MustHaveSymbol:          false,
		},
		Users:  make(map[string]*CUser),
		Groups: make(map[string]*CGroup),
	}
}

type CredentaDB struct {
	DefaultRealm string             `json:"defaultRealm"`
	PassPolicy   *PassphrasePolicy  `json:"passPolicy"`
	Users        map[string]*CUser  `json:"users"`
	Groups       map[string]*CGroup `json:"groups"`
}

func (store *CredentaDB) NewGroup(realm, name, parentGroup string) (*CGroup, error) {
	if realm == "" || name == "" {
		return nil, errors.New("realm and name is required")
	}

	if store.Groups == nil {
		store.Groups = make(map[string]*CGroup)
	}

	if _, ok := store.Groups[fmt.Sprintf("%s IN %s", name, realm)]; ok {
		return nil, fmt.Errorf("group %s already exists", name)
	}

	store.Groups[name] = &CGroup{
		Realm:        realm,
		Name:         name,
		ParentGroups: make([]string, 0),
		Attributes:   make([]*Attribute, 0),
	}

	return store.Groups[name], nil
}

func (store *CredentaDB) NewUser(realm, id, password string, groups []string, idType IdType, vMethod VerificationMethod) (*CUser, error) {
	if realm == "" || id == "" || password == "" {
		return nil, fmt.Errorf("realm, id and password is required")
	}
	valid, err := store.PassPolicy.IsPasswordValid(password)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid email password")
	}
	if _, ok := store.Users[fmt.Sprintf("%s IN %s", id, realm)]; ok {
		return nil, fmt.Errorf("user with id %s already exists", id)
	}
	hash, err := MakeVerification(vMethod, password)
	if err != nil {
		return nil, err
	}
	store.Users[fmt.Sprintf("%s IN %s", id, realm)] = &CUser{
		Realm:              realm,
		Id:                 fmt.Sprintf("%s IN %s", id, realm),
		IDType:             idType,
		Groups:             groups,
		Attributes:         make([]*Attribute, 0),
		RoleMasks:          make([]uint64, 10),
		VerificationMethod: vMethod,
		VerificationHash:   hash,
		Enable:             true,
		Active:             false,
	}
	return store.Users[fmt.Sprintf("%s IN %s", id, realm)], nil
}

func (store *CredentaDB) GetUser(realm, id string) (*CUser, error) {
	if realm == "" || id == "" {
		return nil, errors.New("realm and id are required")
	}
	if usr, ok := store.Users[fmt.Sprintf("%s IN %s", id, realm)]; ok {
		return usr, nil
	}
	return nil, fmt.Errorf("user with id %s not found in %s realm", id, realm)
}

func (store *CredentaDB) GetUserWithAuth(realm, id, password string) (*CUser, error) {
	if realm == "" || id == "" || password == "" {
		return nil, errors.New("realm and id and password are required")
	}
	user, err := store.GetUser(realm, id)
	if err != nil {
		return user, err
	}
	if MatchVerification(user.VerificationMethod, password, user.VerificationHash) {
		return user, nil
	}
	return nil, fmt.Errorf("invalid authentication")
}
