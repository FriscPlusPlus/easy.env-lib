package easyenv

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

// constructor

func NewProject(projectName, path string) *Project {
	project := new(Project)

	project.projectID = uuid.NewString()
	project.SetProjectName(projectName)
	project.SetPath(path)

	return project
}

//  getters

func (prj *Project) GetProjectName() string {
	return prj.projectName
}

func (prj *Project) GetProjectID() string {
	return prj.projectID
}

func (prj *Project) GetPath() string {
	return prj.path
}

func (prj *Project) GetEnvironments() []*DataSet {
	return prj.values
}

func (prj *Project) GetEnvironmentByKey(keyName string) (*DataSet, error) {

	for _, env := range prj.values {
		if env.keyName == keyName {
			return env, nil
		}
	}

	return nil, fmt.Errorf("no enviorment found with the key %s", keyName)
}

// setters

func (prj *Project) SetProjectName(value string) {
	prj.projectName = value
}

func (prj *Project) SetPath(value string) {
	lastChar := string(value[len(value)-1])
	if lastChar != "\\" {
		value = fmt.Sprintf("%s\\", value)
	}
	prj.path = value
}

func (prj *Project) AddEnvrioment(keyName, value string) (*DataSet, error) {

	_, ok := prj.GetEnvironmentByKey(keyName)

	if ok != nil {
		return nil, fmt.Errorf("an enviorment with the key %s already exists", keyName)
	}

	env := NewProjectDtaSet(keyName, value)
	prj.values = append(prj.values, env)
	return env, nil
}

func (prj *Project) RemoveEnviorment(keyName string) {

	tmp := make([]*DataSet, 0)
	foundIndex := 0
	for index, env := range prj.values {
		if env.keyName == keyName {
			foundIndex = index
			break
		}
	}
	tmp = append(tmp, prj.values[:foundIndex]...)
	tmp = append(tmp, prj.values[foundIndex+1:]...)

	prj.values = tmp
}

// this method will remove the enviorment and the related .env file
func (prj *Project) RemoveAllEnviorments() error {
	err := os.Remove(fmt.Sprintf("%s.env", prj.path))

	if err != nil {
		return err
	}

	prj.values = []*DataSet{}

	return nil
}
