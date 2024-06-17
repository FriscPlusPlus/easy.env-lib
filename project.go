package easyenv

func (prj *Project) GetProjectName() string {
	return prj.projectName
}

func (prj *Project) GetProjectID() string {
	return prj.projectID
}

func (prj *Project) GetPath() string {
	return prj.path
}

func (prj *Project) GetEnvironments() []ProjectDataSet {
	return prj.values
}

func (prj *Project) SetProjectName(value string) {
	prj.projectName = value
}

func (prj *Project) SetEnvValue(keyName, value string) {
	for _, env := range prj.values {
		if env.keyName == keyName {
			env.value = value
			break
		}
	}
}

func (prj *Project) AddEnvrioment(keyName, value string) {
	prj.values = append(prj.values, ProjectDataSet{
		keyName: keyName,
		value:   value,
	})
}
