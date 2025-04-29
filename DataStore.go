package credenta

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	credentaDB             *CredentaDB
	authDB                 *AuthDB
	reloadUserChannel      = make(chan string)
	reloadGroupChannel     = make(chan string)
	saveUpdateUserChannel  = make(chan string)
	saveUpdateGroupChannel = make(chan string)
	killChannel            = make(chan bool)
	userMutext             = sync.Mutex{}
	groupMutext            = sync.Mutex{}

	DBStarted bool = false
)

func init() {
	credentaDB = NewCredentaDB()
	aDB, err := NewAuthDB(".", "/data/user", "/data/group")
	if err != nil {
		fmt.Errorf("FATAL NewAuthDB error: %v", err)
	}
	authDB = aDB
}

func StopDB() {
	DBStarted = false
	killChannel <- true
}

func StartDB() error {
	if authDB == nil {
		return fmt.Errorf("authDB is not started")
	}
	if DBStarted {
		return fmt.Errorf("already started")
	}
	DBStarted = true
	go func() {
		for {
			select {
			case userKeyToLoad := <-reloadUserChannel:
				userMutext.Lock()
				defer userMutext.Unlock()

				// TODO do not forget to pass on the realm
				err := authDB.ReloadUser(userKeyToLoad)
				if err != nil {
					fmt.Println(err)
				}
			case groupKeyToLoad := <-reloadGroupChannel:
				groupMutext.Lock()
				defer groupMutext.Unlock()

				// TODO do not forget to pass on the realm
				err := authDB.ReloadGroup(groupKeyToLoad)
				if err != nil {
					fmt.Println(err)
				}
			case userKeyToUpdate := <-saveUpdateUserChannel:
				userMutext.Lock()
				defer userMutext.Unlock()

				// TODO do not forget to pass on the realm
				err := authDB.UpdateOrSaveUser(userKeyToUpdate)
				if err != nil {
					fmt.Println(err)
				}
			case groupKeyToUpdate := <-saveUpdateGroupChannel:
				groupMutext.Lock()
				defer groupMutext.Unlock()

				// TODO do not forget to pass on the realm
				err := authDB.UpdateOrSaveGroup(groupKeyToUpdate)
				if err != nil {
					fmt.Println(err)
				}
			case <-killChannel:
				fmt.Println("Called killChannel")
				DBStarted = false
				return
			default:
				// Optional: Add a default case to avoid blocking
				// if there's no data in the channel
				// fmt.Println("No data yet")
				time.Sleep(100 * time.Millisecond) // avoid busy-waiting
			}
		}
		fmt.Println("Go func returned")
	}()
	return nil
}

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
	}
}

type CredentaDB struct {
	DefaultRealm string            `json:"defaultRealm"`
	PassPolicy   *PassphrasePolicy `json:"passPolicy"`
}

func NewAuthDB(baseFolder, userFolder, groupFolder string) (*AuthDB, error) {
	adb := &AuthDB{
		Users:       make(map[string]*CUser),  // KEY is UID_IN_REALM
		Groups:      make(map[string]*CGroup), // KEY is NAME_IN_REALM
		BaseFolder:  baseFolder,
		UserFolder:  userFolder,
		GroupFolder: groupFolder,
	}
	err := adb.ReloadAuthDB()
	if err != nil {
		return nil, err
	}
	return adb, nil
}

type AuthDB struct {
	Users       map[string]*CUser  `json:"users"`  // KEY is UID_IN_REALM
	Groups      map[string]*CGroup `json:"groups"` // KEY is NAME_IN_REALM
	BaseFolder  string             `json:"baseFolder"`
	UserFolder  string             `json:"userFolder"`
	GroupFolder string             `json:"groupFolder"`
}

func (store *AuthDB) ReloadAuthDB() error {
	userMutext.Lock()
	groupMutext.Lock()
	defer func() {
		groupMutext.Unlock()
		userMutext.Unlock()
	}()

	userMap, err := LoadAllUserData(store.BaseFolder, store.UserFolder)
	if err != nil {
		return fmt.Errorf("in ReloadAuthDB function error: %v", err)
	}
	for realm, userArr := range userMap {
		for _, theUser := range userArr {
			store.Users[fmt.Sprintf("%s_IN_%s", theUser.Id, realm)] = theUser
		}
	}

	groupMap, err := LoadAllGroupData(store.BaseFolder, store.GroupFolder)
	if err != nil {
		return fmt.Errorf("in ReloadAuthDB function error: %v", err)
	}
	for realm, groupArr := range groupMap {
		for _, theGroup := range groupArr {
			store.Groups[fmt.Sprintf("%s_IN_%s", theGroup.Name, realm)] = theGroup
		}
	}

	return nil
}

// TODO Finish this function Update the user file in CredentaDB/BaseFolder/UserFolder/REALM@UID.json
func (store *AuthDB) UpdateOrSaveUser(mapKey string) error {
	fmt.Println("Called UpdateOrSaveUser")
	return nil
}

// TODO Finish this function Update the group file in CredentaDB/BaseFolder/GroupFolder/REALM@Name.json
func (store *AuthDB) UpdateOrSaveGroup(mapKey string) error {
	fmt.Println("Called UpdateOrSaveGroup")
	return nil
}

// TODO Finish this function Read the user file in CredentaDB/BaseFolder/UserFolder/REALM@UID.json
func (store *AuthDB) ReloadUser(mapKey string) error {
	fmt.Println("Called ReloadUser" + mapKey)
	return nil
}

// TODO Finish this function Read the group file in CredentaDB/BaseFolder/GroupFolder/REALM@Name.json
func (store *AuthDB) ReloadGroup(mapKey string) error {
	fmt.Println("Called ReloadGroup")
	return nil
}

func (store *AuthDB) NewDefaultGroup(name string, parentGroup []string) (*CGroup, error) {
	return store.NewGroup(credentaDB.DefaultRealm, name, parentGroup)
}

func (store *AuthDB) NewGroup(realm, name string, parentGroup []string) (*CGroup, error) {
	if realm == "" || name == "" {
		return nil, errors.New("realm and name is required")
	}

	if store.Groups == nil {
		store.Groups = make(map[string]*CGroup)
	}

	if _, ok := store.Groups[fmt.Sprintf("%s_IN_%s", name, realm)]; ok {
		return nil, fmt.Errorf("group %s already exists", name)
	}

	theGroup := &CGroup{
		Realm:        realm,
		Name:         name,
		ParentGroups: parentGroup,
		Attributes:   make([]*Attribute, 0),
	}

	store.Groups[fmt.Sprintf("%s_IN_%s", name, realm)] = theGroup

	return theGroup, nil
}

func (store *AuthDB) NewDefaultUser(id, password string, groups []string, idType IdType, vMethod VerificationMethod) (*CUser, error) {
	return store.NewUser(credentaDB.DefaultRealm, id, password, groups, idType, vMethod)
}

func (store *AuthDB) NewUser(realm, id, password string, groups []string, idType IdType, vMethod VerificationMethod) (*CUser, error) {
	if realm == "" || id == "" || password == "" {
		return nil, fmt.Errorf("realm, id and password is required")
	}

	valid, err := credentaDB.PassPolicy.IsPasswordValid(password)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid email password")
	}
	if _, ok := store.Users[fmt.Sprintf("%s_IN_%s", id, realm)]; ok {
		return nil, fmt.Errorf("user %s already exists", id)
	}
	hash, err := MakeVerification(vMethod, password)
	if err != nil {
		return nil, err
	}
	theUser := &CUser{
		Realm:              realm,
		Id:                 id,
		IDType:             idType,
		Groups:             groups,
		Attributes:         make(map[string]*Attribute),
		RoleMasks:          make([]uint64, 10),
		VerificationMethod: vMethod,
		VerificationHash:   hash,
		Enable:             true,
		Active:             false,
	}

	store.Users[fmt.Sprintf("%s_IN_%s", id, realm)] = theUser

	return theUser, nil
}

func (store *AuthDB) GetDefaultUser(id string) (*CUser, error) {
	return store.GetUser(credentaDB.DefaultRealm, id)
}

func (store *AuthDB) GetUser(realm, id string) (*CUser, error) {
	if realm == "" || id == "" {
		return nil, errors.New("realm and id are required")
	}
	if usr, ok := store.Users[fmt.Sprintf("%s_IN_%s", id, realm)]; ok {
		return usr, nil
	}
	return nil, fmt.Errorf("user with id %s not found in %s realm", id, realm)
}

func (store *AuthDB) GetDefaultUserWithAuth(id, password string) (*CUser, error) {
	return store.GetUserWithAuth(credentaDB.DefaultRealm, id, password)
}

func (store *AuthDB) GetUserWithAuth(realm, id, password string) (*CUser, error) {
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
