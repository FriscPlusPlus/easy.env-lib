package easyenv

import "fmt"

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

func (prj *Project) GetEnvironments() []*ProjectDataSet {
	return prj.values
}

func (prj *Project) GetEnvironmentByKey(keyName string) (*ProjectDataSet, error) {

	for _, env := range prj.values {
		if env.keyName == keyName {
			return env, nil
		}
	}

	return nil, fmt.Errorf("No enviorment found with the key %s", keyName)
}

// setters

func (prj *Project) SetProjectName(value string) {
	prj.projectName = value
}

func (prj *Project) SetPath(value string) {
	prj.path = value
}

func (prj *Project) SetEnvValue(keyName, value string) (*ProjectDataSet, error) {
	for _, env := range prj.values {
		if env.keyName == keyName {
			env.value = value
			return env, nil

		}
	}

	return nil, fmt.Errorf("No enviorment found with the key %s", keyName)
}

func (prj *Project) AddEnvrioment(keyName, value string) (*ProjectDataSet, error) {

	_, ok := prj.GetEnvironmentByKey(keyName)

	if ok != nil {
		return nil, fmt.Errorf("An enviorment with the key %s already exists!", keyName)
	}

	env := new(ProjectDataSet)
	env.keyName = keyName
	env.value = value
	prj.values = append(prj.values, env)
	return env, nil
}
