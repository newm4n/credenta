package credenta

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func userFileName(baseFolder, userFolder, realm, id string) (string, error) {
	if baseFolder == "" || userFolder == "" || realm == "" || id == "" {
		return "", fmt.Errorf("in userFileName function, baseFolder, userFolder, realm and id are required")
	}
	return fmt.Sprintf("%s%s/%s_IN_%s.json", baseFolder, userFolder, id, realm), nil
}
func groupFileName(baseFolder, groupFolder, realm, name string) (string, error) {
	if baseFolder == "" || groupFolder == "" || realm == "" || name == "" {
		return "", fmt.Errorf("in groupFileName function, baseFolder, groupFolder,  realm and id are required")
	}
	return fmt.Sprintf("%s%s/%s_IN_%s.json", baseFolder, groupFolder, name, realm), nil
}

/*
ListUserDataFiles will return a map of realm name to array of user id. The function will go to directory with format
`BaseFolder/UserFolder` and look for file with `USERID_IN_REALM.json` name. It will return an error if no folder with
that name is found. By default, the BaseFolder is "." which equals to the name of the project.
*/
func ListUserDataFiles(baseFolder, userFolder string) (map[string][]string, error) {
	if baseFolder == "" || userFolder == "" {
		return nil, fmt.Errorf("in ListUserDataFiles function, baseFolder and userFolder are required")
	}
	entries, err := os.ReadDir(fmt.Sprintf("%s%s", baseFolder, userFolder))
	if err != nil {
		return nil, fmt.Errorf("in ListUserDataFiles function, error reading directory %s%s: %w", baseFolder, userFolder, err)
	}
	return listDataFiles(entries), nil
}

/*
ListGroupDataFiles will return a map of realm name to array of group name. The function will go to directory with format
`BaseFolder/GroupFolder` and look for file with `NAME_IN_REALM.json` name. It will return an error if no folder with
that name is found. By default, the BaseFolder is "." which equals to the name of the project.
*/
func ListGroupDataFiles(baseFolder, groupFolder string) (map[string][]string, error) {
	if baseFolder == "" || groupFolder == "" {
		return nil, fmt.Errorf("in ListGroupDataFiles function, baseFolder and groupFolder are required")
	}
	entries, err := os.ReadDir(fmt.Sprintf("%s%s", baseFolder, groupFolder))
	if err != nil {
		return nil, fmt.Errorf("in ListGroupDataFiles function, error reading directory %s%s: %w", baseFolder, groupFolder, err)
	}
	return listDataFiles(entries), nil
}

/*
listDataFiles return  list map of realm name to data string for each entries. This function will be called by
ListUserDataFiles or ListGroupDataFiles function.
*/
func listDataFiles(entries []os.DirEntry) map[string][]string {
	ret := make(map[string][]string)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			n := strings.Split(entry.Name(), ".")
			ne := strings.Split(n[0], "_IN_")
			id := ne[0]
			realm := ne[1]
			if ids, ok := ret[realm]; ok {
				ids = append(ids, id)
			} else {
				ret[realm] = make([]string, 0)
				ret[realm] = append(ret[realm], id)
			}
		}
	}
	return ret
}

/*
LoadAllUserData will return map of realm to list of CUser objects.
It will list all user files in the BaseFolder + userFolder and read all of them.
*/
func LoadAllUserData(baseFolder, userFolder string) (map[string][]*CUser, error) {
	if baseFolder == "" || userFolder == "" {
		return nil, fmt.Errorf("in LoadAllUserData function, baseFolder and userFolder are required")
	}
	userFiles, err := ListUserDataFiles(baseFolder, userFolder)
	if err != nil {
		return nil, fmt.Errorf("in LoadAllUserData function: %w", err)
	}
	ret := make(map[string][]*CUser)
	for realm, ids := range userFiles {
		for _, id := range ids {
			theUser, err := LoadUserData(baseFolder, userFolder, realm, id)
			if err != nil {
				return nil, fmt.Errorf("in LoadAllUserData function. error on LoadUserData for realm %s and user id %s: %w", realm, id, err)
			}
			if userDataArray, ok := ret[realm]; ok {
				userDataArray = append(userDataArray, theUser)
			} else {
				ret[realm] = make([]*CUser, 0)
				ret[realm] = append(ret[realm], theUser)
			}
		}
	}
	return ret, nil
}

/*
LoadAllGroupData will return map of realm to list of CGroup objects.
It will list all group files in the BaseFolder + groupFolder and read all of them.
*/
func LoadAllGroupData(baseFolder, groupFolder string) (map[string][]*CGroup, error) {
	if baseFolder == "" || groupFolder == "" {
		return nil, fmt.Errorf("in LoadAllGroupData function, baseFolder and groupFolder are required")
	}
	groupFiles, err := ListGroupDataFiles(baseFolder, groupFolder)
	if err != nil {
		return nil, fmt.Errorf("in LoadAllGroupData function: %w", err)
	}
	ret := make(map[string][]*CGroup)
	for realm, names := range groupFiles {
		for _, name := range names {
			theGroup, err := LoadGroupData(baseFolder, groupFolder, realm, name)
			if err != nil {
				return nil, fmt.Errorf("in LoadAllGroupData function. error on LoadGroupData for realm %s and name %s: %w", realm, name, err)
			}
			if groupDataArray, ok := ret[realm]; ok {
				groupDataArray = append(groupDataArray, theGroup)
			} else {
				ret[realm] = make([]*CGroup, 0)
				ret[realm] = append(ret[realm], theGroup)
			}
		}
	}
	return ret, nil
}

/*
LoadUserData will try to locate the user's file based on its realm and id. If the file name si correct, it will
try to open and read the user information from the file.
*/
func LoadUserData(baseFolder, userFolder, realm, id string) (*CUser, error) {
	if baseFolder == "" || userFolder == "" || realm == "" || id == "" {
		return nil, fmt.Errorf("in LoadUserData function, baseFolder, userFolder, realm and id are required")
	}
	fileName, err := userFileName(baseFolder, userFolder, realm, id)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("in LoadUserData function, error opening file %s: %w", fileName, err)
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
		return nil, fmt.Errorf("in LoadUserData function, error unmarshaling data into CUser: %w", err)
	}
	return nUser, nil
}

/*
LoadGroupData will try to locate the group's file based on its realm and name. If the file name si correct, it will
try to open and read the group information from the file.
*/
func LoadGroupData(baseFolder, groupFolder, realm, name string) (*CGroup, error) {
	if baseFolder == "" || groupFolder == "" || realm == "" || name == "" {
		return nil, fmt.Errorf("in LoadGroupData function, baseFolder, groupFolder, realm and name are required")
	}
	fileName, err := userFileName(baseFolder, groupFolder, realm, name)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("in LoadGroupData function, error opening file %s: %w", fileName, err)
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
		return nil, fmt.Errorf("in LoadGroupData function, error unmarshaling data into CGroup: %w", err)
	}
	return nGroup, nil
}

/*
StoreUserData will marshall the content of usr and write it json content into it appropriate json file in
its respective folder location.
*/
func StoreUserData(baseFolder, userFolder string, usr *CUser) error {
	if baseFolder == "" || userFolder == "" {
		return fmt.Errorf("in StoreUserData function, baseFolder and userFolder are required")
	}
	if usr == nil {
		return fmt.Errorf("in StoreUserData function. nil user")
	}

	dataBytes, err := json.Marshal(usr)
	if err != nil {
		return fmt.Errorf("in StoreUserData function. error json marshal: %w", err)
	}

	fileName, err := userFileName(baseFolder, userFolder, usr.Realm, usr.Id)
	if err != nil {
		return fmt.Errorf("in StoreUserData function. error calling userFileName function: %w", err)
	}

	return writeDataToFile(dataBytes, fileName)
}

/*
StoreGroupData will marshall the content of grp and write it json content into it appropriate json file in
its respective folder location.
*/
func StoreGroupData(baseFolder, groupFolder string, grp *CGroup) error {
	if baseFolder == "" || groupFolder == "" {
		return fmt.Errorf("in StoreGroupData function, baseFolder and groupFolder are required")
	}
	if grp == nil {
		return fmt.Errorf("in StoreGroupData function. nil grp")
	}

	dataBytes, err := json.Marshal(grp)
	if err != nil {
		return fmt.Errorf("in StoreGroupData function. error json marshal: %w", err)
	}

	fileName, err := groupFileName(baseFolder, groupFolder, grp.Realm, grp.Name)
	if err != nil {
		return fmt.Errorf("in StoreGroupData function. error calling groupFileName function: %w", err)
	}

	return writeDataToFile(dataBytes, fileName)
}

/*
writeDataToFile will write the byte array data into file specified by fileName.
If the file is already exist, this function will truncate it and replace it content.
If the file is not exist, this function will first try to create the file with the specified content.
*/
func writeDataToFile(data []byte, fileName string) error {
	if data == nil {
		return fmt.Errorf("in writeDataToFile function. nil data")
	}
	if fileName == "" {
		return fmt.Errorf("in writeDataToFile function. empty filename")
	}
	if _, err := os.Stat(fileName); err == nil {
		f, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("in writeDataToFile function. error opening file %s: %w", fileName, err)
		}
		defer f.Close()
		err = f.Truncate(0)
		if err != nil {
			return fmt.Errorf("in writeDataToFile function. error truncate file %s: %w", fileName, err)
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("in writeDataToFile function. error seek in file %s: %w", fileName, err)
		}
		_, err = f.Write(data)
		if err != nil {
			return fmt.Errorf("in writeDataToFile function. error writing into file %s: %w", fileName, err)
		}
	} else if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("in writeDataToFile function. error creating file %s: %w", fileName, err)
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf("in writeDataToFile function. error writing data to file %s: %w", fileName, err)
		}
	} else {
		return fmt.Errorf("in writeDataToFile function. error obtaining stat of file %s: %w", fileName, err)
	}
	return nil
}
