package easyenv

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

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

func (easy *EasyEnv) SaveDB() error {

	err := easy.isCurrentDBSet()

	if err != nil {
		return err
	}

	err = saveDataInDB(easy.currentConnection)

	if err != nil {
		return err
	}

	err = saveEnvInFile(easy.currentConnection)

	if err != nil {
		return err
	}

	return nil
}

func (easy *EasyEnv) AddProject(projectName, path string) (*Project, error) {
	err := easy.isCurrentDBSet()

	if err != nil {
		return nil, err
	}

	project := NewProject(projectName, path)

	easy.currentConnection.projects = append(easy.currentConnection.projects, project)

	return project, nil
}

func (easy *EasyEnv) AddTemplate(templateName string) (*Template, error) {
	err := easy.isCurrentDBSet()

	if err != nil {
		return nil, err
	}

	template := NewTemplate(templateName)

	easy.currentConnection.templates = append(easy.currentConnection.templates, template)

	return template, nil
}

func (easy *EasyEnv) RemoveProject(projectID string) error {

	err := easy.isCurrentDBSet()

	if err != nil {
		return err
	}

	err = removeData(easy.currentConnection, "projects", "projectID", projectID)

	if err != nil {
		return err
	}

	tmp := make([]*Project, 0)
	foundIndex, _, err := easy.GetProject(projectID)

	if err != nil {
		return err
	}

	tmp = append(tmp, easy.currentConnection.projects[:foundIndex]...)
	tmp = append(tmp, easy.currentConnection.projects[foundIndex+1:]...)
	easy.currentConnection.projects = tmp
	return nil
}

func (easy *EasyEnv) RemoveTemplate(templateID string) error {

	err := easy.isCurrentDBSet()

	if err != nil {
		return err
	}

	err = removeData(easy.currentConnection, "templates", "templateID", templateID)

	if err != nil {
		return err
	}

	tmp := make([]*Template, 0)
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


/*
 Getters
*/

func (easy *EasyEnv) GetProject(projectID string) (int, *Project, error) {

	err := easy.isCurrentDBSet()

	if err != nil {
		return 0, nil, err
	}

	foundIndex := 0
	var foundProject *Project
	for index, project := range easy.currentConnection.projects {
		if project.GetProjectID() == projectID {
			foundIndex = index
			foundProject = project
			return foundIndex, foundProject, nil
		}
	}
	return foundIndex, foundProject, fmt.Errorf("no project found with ID %s. Please check the ID and try again", projectID)
}

func (easy *EasyEnv) GetProjects() ([]*Project, error) {
	err := easy.isCurrentDBSet()

	if err != nil {
		return nil, err
	}

	return easy.currentConnection.projects, nil
}

func (easy *EasyEnv) GetTemplate(templateID string) (int, *Template, error) {
	err := easy.isCurrentDBSet()

	if err != nil {
		return 0, nil, err
	}

	foundIndex := 0
	var foundTemplate *Template
	for index, template := range easy.currentConnection.templates {
		if template.GetTemplateID() == templateID {
			foundIndex = index
			foundTemplate = template
			return foundIndex, foundTemplate, nil
		}
	}
	return 0, foundTemplate, fmt.Errorf("no template found with ID %s. Please verify the ID and try again", templateID)
}

func (easy *EasyEnv) GetTemplates() ([]*Template, error) {
	err := easy.isCurrentDBSet()

	if err != nil {
		return nil, err
	}

	return easy.currentConnection.templates, nil
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

func (easy *EasyEnv) isCurrentDBSet() error {

	if easy.currentConnection == nil {
		return fmt.Errorf("no database is currently open. Please open a database first using 'Open(path/to/sqlitefile)' before making any other calls")
	}

	return nil
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

func createEnvString(environments []*DataSet) string {
	var result string

	for _, env := range environments {
		result += fmt.Sprintf("%s=%s\n", env.keyName, env.value)
	}

	return result
}
