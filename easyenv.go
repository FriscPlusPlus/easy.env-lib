package easyenv

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

type EasyEnvDefinition interface {
	NewEasyEnv() *EasyEnv
	Load(dbName string) (*Connection, error)
	Open(dbName string) (*Connection, error)
	CloseDB(dbName string) error
	SaveCurrentDB() error // this will save the data from the buffer to the current db and write for each project env  file
	SaveAllDB() error     // this will save the data from the buffer to the all the open db and write for each project env  file

	CreateNewDB(dbName string) (*Connection, error)

	AddProject(projectName, path string) error
	AddTemplate(templateName string) error

	RemoveProject(project Project) error
	RemoveTemplate(template Template) error

	AddEnvToProject(projectID, keyName, value string) error
	AddEnvToTemplate(template, keyName, value string) error

	RemoveEnvFromProject(projectID, keyName, value string) error
	RemoveEnvFromTemplate(projectID, keyName, value string) error

	LoadTemplates() error

	LoadProjects() error

	LoadTemplateIntoProject(template, projectID string) error
}

func NewEasyEnv() *EasyEnv {
	return new(EasyEnv)
}

func (easy *EasyEnv) Load(dbName string) (*Connection, error) {
	db, err := sql.Open("sqlite3", dbName)

	connection := new(Connection)

	if err != nil {
		return nil, err
	}

	connection.dbName = dbName
	connection.db = db

	easy.connections = append(easy.connections, connection)
	easy.currentConnection = connection
	return easy.currentConnection, nil
}

func (easy *EasyEnv) Open(dbName string) (*Connection, error) {
	connection, err := easy.getConnectionByDBname(dbName)

	if err != nil {
		return nil, err
	}

	easy.currentConnection = connection
	return connection, nil
}

func (easy *EasyEnv) CloseDB(dbName string) error {
	connection, err := easy.getConnectionByDBname(dbName)

	if err != nil {
		return err
	}

	err = connection.db.Close()

	if err != nil {
		return err
	}

	easy.removeConnection(dbName)

	if easy.currentConnection.dbName == dbName {
		easy.currentConnection = nil
	}

	return nil
}

func (easy *EasyEnv) CreateNewDB(dbName string) (*Connection, error) {
	connection, err := easy.Load(dbName)

	if err != nil {
		return nil, err
	}

	err = createTables(connection)

	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (easy *EasyEnv) SaveCurrentDB() error {

	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	err = saveDataInDB(easy.currentConnection)

	if err != nil {
		return err
	}

	err = saveEnvInFile(easy.currentConnection)

	easy.resetMethodState()

	return nil
}

func (easy *EasyEnv) AddProject(projectName, path string) error {

	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	var project Project
	project.projectName = projectName
	project.path = path
	project.method = "INSERT"
	easy.currentConnection.projects = append(easy.currentConnection.projects, project)

	return nil
}

func (easy *EasyEnv) AddTemplate(templateName string) error {

	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	var template Template

	template.templateName = templateName
	template.method = "INSERT"
	easy.currentConnection.templates = append(easy.currentConnection.templates, template)

	return nil
}

func (easy *EasyEnv) RemoveProject(projectID int) error {

	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	err = removeData(easy.currentConnection, "projects", "projectID", projectID)

	if err != nil {
		return err
	}

	tmp := make([]Project, 0)
	foundIndex, _ := easy.getProjectByID(projectID)

	tmp = append(tmp, easy.currentConnection.projects[:foundIndex]...)
	tmp = append(tmp, easy.currentConnection.projects[foundIndex+1:]...)
	easy.currentConnection.projects = tmp
	return nil
}

func (easy *EasyEnv) RemoveTemplate(templateID int) error {

	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	err = removeData(easy.currentConnection, "templates", "templateID", templateID)

	if err != nil {
		return err
	}

	tmp := make([]Template, 0)
	foundIndex := 0

	for index, project := range easy.currentConnection.templates {
		if project.templateID == templateID {
			foundIndex = index
			break
		}
	}

	tmp = append(tmp, easy.currentConnection.templates[:foundIndex]...)
	tmp = append(tmp, easy.currentConnection.templates[foundIndex+1:]...)
	easy.currentConnection.templates = tmp

	return nil
}

func (easy *EasyEnv) AddEnvToProject(projectID int, keyName, value string) error {

	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	env := ProjectDataSet{
		keyName: keyName,
		value:   value,
	}

	projects := easy.currentConnection.projects
	index, _ := easy.getProjectByID(projectID)

	projects[index].values = append(projects[index].values, env)

	return nil
}

func (easy *EasyEnv) AddEnvToTemplate(templateID int, keyName, value string) error {
	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	env := TemplateDataSet{
		templateID: templateID,
		keyName:    keyName,
		value:      value,
		method:     "INSERT",
	}

	templates := easy.currentConnection.templates
	index, _ := easy.getTemplateByID(templateID)

	templates[index].values = append(templates[index].values, env)

	return nil
}

/*
	Unexported methods
*/

func (easy *EasyEnv) getConnectionByDBname(dbName string) (*Connection, error) {
	for _, connection := range easy.connections {
		if connection.dbName == dbName {
			return connection, nil
		}
	}
	return nil, fmt.Errorf("no connection found for the database with the name: %s", dbName)
}

func (easy *EasyEnv) removeConnection(dbName string) {
	tmp := make([]*Connection, 0)
	foundIndex := 0
	for index, connection := range easy.connections {
		if connection.dbName == dbName {
			foundIndex = index
			break
		}
	}
	tmp = append(tmp, easy.connections[:foundIndex]...)
	tmp = append(tmp, easy.connections[foundIndex+1:]...)
	easy.connections = tmp
}

func (easy *EasyEnv) checkIfcurrentDBisSet() error {

	if easy.currentConnection == nil {
		return fmt.Errorf("no database is currently open. Please open a database first using 'Open(path/to/sqlitefile)' before making any other calls")
	}

	return nil
}

func (easy *EasyEnv) resetMethodState() {
	projects := easy.currentConnection.projects
	templates := easy.currentConnection.templates

	for i := range projects {
		if len(projects[i].method) > 0 {
			projects[i].method = ""
		}
	}

	for i := range templates {
		if len(templates[i].method) > 0 {
			templates[i].method = ""
		}
	}
}

func (easy *EasyEnv) getProjectByID(projectID int) (int, Project) {
	foundIndex := 0
	var foundProject Project
	for index, project := range easy.currentConnection.projects {
		if project.projectID == projectID {
			foundIndex = index
			foundProject = project
			break
		}
	}
	return foundIndex, foundProject
}

func (easy *EasyEnv) getTemplateByID(templateID int) (int, Template) {
	foundIndex := 0
	var foundTemplate Template
	for index, template := range easy.currentConnection.templates {
		if template.templateID == templateID {
			foundIndex = index
			foundTemplate = template
			break
		}
	}
	return foundIndex, foundTemplate
}

func saveEnvInFile(connection *Connection) error {
	var err error
	for _, project := range connection.projects {
		envPath := path.Join(project.path, ".env")
		os.Remove(envPath)

		envString := createEnvString(project.values)

		err = os.WriteFile(envPath, []byte(envString), 0644)
		
		if err != nil {
			return err
		}
	}
	return err
}

func createEnvString(enviorments []ProjectDataSet) string {
	var result string

	for _, env := range enviorments {
		result += fmt.Sprintf("%s=%s\n", env.keyName, env.value)
	}

	return result
}
