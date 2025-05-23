package credenta

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	RoleMaskCount = 10
)

func getEnvVar(varName string, defaultValue string, options []string) string {
	value, present := os.LookupEnv(varName)
	if present {
		return value
	} else {
		if options == nil || len(options) == 0 {
			log.Printf("%s environment variable not set. Set to default value \"%s\".\n", varName, defaultValue)
		} else {
			log.Printf("%s environment variable not set. Set to default value \"%s\". Options are %s.\n", varName, defaultValue, strings.Join(options, ","))
		}
		return defaultValue
	}
}

func NewCredentaDB() (*CredentaDB, error) {

	baseFolder := getEnvVar("CREDENTA_BASE_DIR", ".", nil)
	userFolder := getEnvVar("CREDENTA_USER_DIR", "/data/user", nil)
	groupFolder := getEnvVar("CREDENTA_GROUP_DIR", "/data/group", nil)
	defaultRealm := getEnvVar("CREDENTA_REALM_DEFAULT", "DEFAULT", nil)
	passPolicy := getEnvVar("CREDENTA_PASS_POLICY", "SIMPLE", []string{"SIMPLE", "STRONG", "CLASSIC"})

	passphrasePolicy := &PassphrasePolicy{}

	switch passPolicy {
	case "STRONG":
		passphrasePolicy = StrongPasswordPolicy()
	case "CLASSIC":
		passphrasePolicy = ClassicPasswordPolicy()
	case "SIMPLE":
		passphrasePolicy = SimplePasswordPolicy()
	default:
		passphrasePolicy = SimplePasswordPolicy()
	}

	cDB := &CredentaDB{
		DefaultRealm: defaultRealm,
		PassPolicy:   passphrasePolicy,
		BaseFolder:   baseFolder,
		UserFolder:   userFolder,
		GroupFolder:  groupFolder,
	}

	if _, err := os.Stat(fmt.Sprintf("%s%s", baseFolder, userFolder)); err != nil {
		return nil, fmt.Errorf("could not find user directory \"%s%s\". Please create the directory or change the environment variable CREDENTA_BASE_DIR and/or CREDENTA_USER_DIR and try again ", baseFolder, userFolder)
	}

	if _, err := os.Stat(fmt.Sprintf("%s%s", baseFolder, groupFolder)); err != nil {
		return nil, fmt.Errorf("could not find user directory \"%s%s\". Please create the directory or change the environment variable CREDENTA_BASE_DIR and/or CREDENTA_GROUP_DIR and try again ", baseFolder, groupFolder)
	}

	return cDB, nil
}

type CredentaDB struct {
	DefaultRealm string            `json:"defaultRealm"`
	PassPolicy   *PassphrasePolicy `json:"passPolicy"`
	BaseFolder   string            `json:"baseFolder"`
	UserFolder   string            `json:"userFolder"`
	GroupFolder  string            `json:"groupFolder"`
}

func (store *CredentaDB) GetRoleMasksOfGroups(ctx context.Context, realm, group string) []uint64 {
	ret := make([]uint64, RoleMaskCount)
	theGroup, err := store.GetGroup(ctx, realm, group)
	if err != nil || theGroup == nil {
		return ret
	}

	if theGroup.ParentGroups == nil || len(theGroup.ParentGroups) == 0 {
		return theGroup.RoleMasks
	} else {
		ret = theGroup.RoleMasks
		for _, parentGroupName := range theGroup.ParentGroups {
			parentRoles := store.GetRoleMasksOfGroups(ctx, realm, parentGroupName)
			for i, rete := range ret {
				ret[i] = rete | parentRoles[i]
			}
		}
		return ret
	}
}

func (store *CredentaDB) NewDefaultGroup(ctx context.Context, name string, parentGroup []string) (*CGroup, error) {
	return store.NewGroup(ctx, store.DefaultRealm, name, parentGroup)
}

func (store *CredentaDB) NewGroup(ctx context.Context, realm, name string, parentGroup []string) (*CGroup, error) {
	if realm == "" || name == "" {
		return nil, errors.New("realm and name is required")
	}

	groupFileName := fmt.Sprintf("%s%s/%s_IN_%s.json", store.BaseFolder, store.GroupFolder, name, realm)
	_, err := os.Stat(groupFileName)
	if err == nil {
		return nil, errors.New("group already exists")
	}

	theGroup := &CGroup{
		FilePath:     groupFileName,
		Realm:        realm,
		Name:         name,
		ParentGroups: parentGroup,
		Attributes:   make([]*Attribute, 0),
		RoleMasks:    make([]uint64, RoleMaskCount),

		CreatedAt: time.Now(),
		CreatedBy: ctx.Value(ETX_USER).(string),
		UpdatedAt: time.Now(),
		UpdatedBy: ctx.Value(ETX_USER).(string),
	}

	return theGroup, nil
}

func (store *CredentaDB) ChangeUserPassword(ctx context.Context, realm, user, password string, vMethod VerificationMethod) error {
	theUser, err := store.GetUser(ctx, realm, user)
	if err != nil {
		return err
	}
	valid, err := store.PassPolicy.IsPasswordValid(password)
	if err != nil || !valid {
		return fmt.Errorf("invalid password format")
	}

	hash, err := MakeVerification(vMethod, password)
	if err != nil {
		return err
	}

	theUser.VerificationHash = hash
	theUser.VerificationMethod = vMethod

	return nil
}

func (store *CredentaDB) NewDefaultUser(ctx context.Context, id, password string, groups []string, idType IdType, vMethod VerificationMethod) (*CUser, error) {
	return store.NewUser(ctx, store.DefaultRealm, id, password, groups, idType, vMethod)
}

func (store *CredentaDB) NewUser(ctx context.Context, realm, id, password string, groups []string, idType IdType, vMethod VerificationMethod) (*CUser, error) {
	if realm == "" || id == "" || password == "" {
		return nil, fmt.Errorf("in NewUser function. realm, id and password is required")
	}

	valid, err := store.PassPolicy.IsPasswordValid(password)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid email password")
	}

	hash, err := MakeVerification(vMethod, password)
	if err != nil {
		return nil, err
	}

	userFileName := fmt.Sprintf("%s%s/%s_IN_%s.json", store.BaseFolder, store.UserFolder, id, realm)
	_, err = os.Stat(userFileName)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	theUser := &CUser{
		FilePath:           userFileName,
		Realm:              realm,
		Id:                 id,
		IDType:             idType,
		Groups:             groups,
		Attributes:         make(map[string]*Attribute),
		RoleMasks:          make([]uint64, RoleMaskCount),
		VerificationMethod: vMethod,
		VerificationHash:   hash,
		Enable:             true,
		Active:             false,

		CreatedAt: time.Time{},
		CreatedBy: ctx.Value(ETX_USER).(string),
		UpdatedAt: time.Time{},
		UpdatedBy: ctx.Value(ETX_USER).(string),
	}

	return theUser, nil
}

func (store *CredentaDB) GetDefaultGroup(ctx context.Context, name string) (*CGroup, error) {
	return store.GetGroup(ctx, store.DefaultRealm, name)
}

func (store *CredentaDB) GetGroup(ctx context.Context, realm, name string) (*CGroup, error) {
	if realm == "" || name == "" {
		return nil, fmt.Errorf("in GetGroup function. realm and name are required")
	}

	theGroup := &CGroup{
		FilePath: fmt.Sprintf("%s%s/%s_IN_%s.json", store.BaseFolder, store.GroupFolder, name, realm),
	}

	err := theGroup.ReloadFromFile(ctx)
	if err != nil {
		return nil, err
	}
	return theGroup, nil
}

func (store *CredentaDB) GetDefaultUser(ctx context.Context, id string) (*CUser, error) {
	return store.GetUser(ctx, store.DefaultRealm, id)
}

func (store *CredentaDB) GetUser(ctx context.Context, realm, id string) (*CUser, error) {
	if realm == "" || id == "" {
		return nil, fmt.Errorf("in GetUser function. realm and id are required")
	}

	theUser := &CUser{
		FilePath: fmt.Sprintf("%s%s/%s_IN_%s.json", store.BaseFolder, store.UserFolder, id, realm),
	}

	err := theUser.ReloadFromFile(ctx)
	if err != nil {
		return nil, err
	}
	return theUser, nil
}

func (store *CredentaDB) GetDefaultUserWithAuth(ctx context.Context, id, password string) (*CUser, []uint64, error) {
	return store.GetUserWithAuth(ctx, store.DefaultRealm, id, password)
}

func (store *CredentaDB) GetUserWithAuth(ctx context.Context, realm, id, password string) (*CUser, []uint64, error) {
	if realm == "" || id == "" || password == "" {
		return nil, nil, errors.New("in GetUserWithAuth function. realm and id and password are required")
	}
	user, err := store.GetUser(ctx, realm, id)
	if err != nil {
		return nil, nil, fmt.Errorf("in GetUserWithAuth function : %v", err)
	}
	if MatchVerification(user.VerificationMethod, password, user.VerificationHash) {
		if !user.Active {
			return nil, nil, errors.New("in GetUserWithAuth function. User is not activated")
		}
		if !user.Enable {
			return nil, nil, errors.New("in GetUserWithAuth function. User is disabled")
		}
		if user.Groups == nil || len(user.Groups) == 0 {
			return user, user.RoleMasks, nil
		} else {
			ret := user.RoleMasks
			for _, grp := range user.Groups {
				grpRoleMask := store.GetRoleMasksOfGroups(ctx, realm, grp)
				for i := 0; i < RoleMaskCount; i++ {
					ret[i] = ret[i] | grpRoleMask[i]
				}
			}
			return user, ret, nil
		}
	}
	return nil, nil, fmt.Errorf("invalid authentication")
}

/*
ListUserIDs will return a map of realm name to array of user id. The function will go to directory with format
`BaseFolder/UserFolder` and look for file with `USERID_IN_REALM.json` name. It will return an error if no folder with
that name is found. By default, the BaseFolder is "." which equals to the name of the project.
*/
func (store *CredentaDB) ListUserIDs(ctx context.Context) (map[string][]string, error) {
	if !pathExists(fmt.Sprintf("%s%s", store.BaseFolder, store.UserFolder)) {
		return nil, fmt.Errorf("in ListUserDataFiles function, folder %s%s not exists", store.BaseFolder, store.UserFolder)
	}
	entries, err := os.ReadDir(fmt.Sprintf("%s%s", store.BaseFolder, store.UserFolder))
	if err != nil {
		return nil, fmt.Errorf("in ListUserDataFiles function, error reading directory %s%s: %w", store.BaseFolder, store.UserFolder, err)
	}
	return store.listDataFiles(ctx, entries), nil
}

/*
ListGroupNames will return a map of realm name to array of group name. The function will go to directory with format
`BaseFolder/GroupFolder` and look for file with `NAME_IN_REALM.json` name. It will return an error if no folder with
that name is found. By default, the BaseFolder is "." which equals to the name of the project.
*/
func (store *CredentaDB) ListGroupNames(ctx context.Context) (map[string][]string, error) {
	if !pathExists(fmt.Sprintf("%s%s", store.BaseFolder, store.GroupFolder)) {
		return nil, fmt.Errorf("in ListGroupDataFiles function, folder %s%s not exists", store.BaseFolder, store.GroupFolder)
	}
	entries, err := os.ReadDir(fmt.Sprintf("%s%s", store.BaseFolder, store.GroupFolder))
	if err != nil {
		return nil, fmt.Errorf("in ListGroupDataFiles function, error reading directory %s%s: %w", store.BaseFolder, store.GroupFolder, err)
	}
	return store.listDataFiles(ctx, entries), nil
}

/*
listDataFiles return  list map of realm name to data string for each entries. This function will be called by
ListUserDataFiles or ListGroupDataFiles function.
*/
func (store *CredentaDB) listDataFiles(ctx context.Context, entries []os.DirEntry) map[string][]string {
	ret := make(map[string][]string)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			n := strings.Split(entry.Name(), ".")
			ne := strings.Split(n[0], "_IN_")
			id := ne[0]
			realm := ne[1]
			if _, ok := ret[realm]; ok {
				ret[realm] = append(ret[realm], id)
			} else {
				ret[realm] = make([]string, 0)
				ret[realm] = append(ret[realm], id)
			}
		}
	}
	return ret
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsRoleFlagOn(roles []uint64, roleSquence int) bool {
	seq, bit := toUint64ByBit(roleSquence)
	return isBitFlagOn(roles[seq], bit)
}
