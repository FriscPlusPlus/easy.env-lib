package easyenv

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type EasyEnvDefinition interface {
	NewEasyEnv() *EasyEnv
	Load(dbName string) (*Connection, error)
	Open(dbName string) (*Connection, error)
	CloseDB(dbName string) error
	SaveDB() error // this will save the data from the buffer to the current db and write for each project env  file

	CreateNewDB(dbName string) (*Connection, error)

	AddProject(projectName, path string) error
	AddTemplate(templateName string) error

	RemoveProject(project Project) error
	RemoveTemplate(template Template) error

	AddEnvToProject(projectID, keyName, value string) error
	AddEnvToTemplate(template, keyName, value string) error

	RemoveEnvFromProject(projectID, keyName string) error
	RemoveEnvFromTemplate(projectID, keyName, value string) error

	GetTemplateByID(templateID string) (Template, error)
	GetProjectByID(projectID string) (Project, error)

	GetAllTemplates() ([]Template, error)
	GetAllProjects() ([]Project, error)

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
	project := new(Project)
	err := easy.isCurrentDBSet()

	if err != nil {
		return project, err
	}

	project.projectID = uuid.NewString()
	project.projectName = projectName
	project.path = path
	project.method = "INSERT"
	easy.currentConnection.projects = append(easy.currentConnection.projects, project)

	return project, nil
}

func (easy *EasyEnv) AddTemplate(templateName string) (*Template, error) {
	template := new(Template)
	err := easy.isCurrentDBSet()

	if err != nil {
		return template, err
	}

	template.templateID = uuid.NewString()
	template.templateName = templateName
	template.method = "INSERT"
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
	foundIndex, _, err := easy.GetProjectByID(projectID)

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

func (easy *EasyEnv) AddEnvToProject(projectID string, keyName, value string) error {

	err := easy.isCurrentDBSet()

	if err != nil {
		return err
	}

	env := ProjectDataSet{
		keyName: keyName,
		value:   value,
	}

	projects := easy.currentConnection.projects
	index, _, err := easy.GetProjectByID(projectID)

	if err != nil {
		return err
	}

	projects[index].values = append(projects[index].values, env)

	return nil
}

func (easy *EasyEnv) AddEnvToTemplate(templateID string, keyName, value string) error {
	err := easy.isCurrentDBSet()

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
	index, _, err := easy.GetTemplateByID(templateID)

	if err != nil {
		return err
	}

	templates[index].values = append(templates[index].values, env)

	return nil
}

func (easy *EasyEnv) RemoveEnvFromProject(projectID string, keyName string) error {
	err := easy.isCurrentDBSet()

	if err != nil {
		return err
	}

	projectIndex, project, err := easy.GetProjectByID(projectID)

	if err != nil {
		return err
	}

	tmp := make([]ProjectDataSet, 0)
	foundIndex := 0
	for index, env := range project.values {
		if env.keyName == keyName {
			foundIndex = index
			break
		}
	}
	tmp = append(tmp, project.values[:foundIndex]...)
	tmp = append(tmp, project.values[foundIndex+1:]...)

	easy.currentConnection.projects[projectIndex].values = tmp

	return nil
}

func (easy *EasyEnv) RemoveEnvFromTemplate(templateID string, keyName string) error {

	err := easy.isCurrentDBSet()

	if err != nil {
		return err
	}

	err = removeTemplateEnvData(easy.currentConnection, templateID, keyName)

	if err != nil {
		return err
	}

	templateIndex, template, err := easy.GetTemplateByID(templateID)

	if err != nil {
		return err
	}

	tmp := make([]TemplateDataSet, 0)
	foundIndex := 0
	for index, env := range template.values {
		if env.keyName == keyName {
			foundIndex = index
			break
		}
	}
	tmp = append(tmp, template.values[:foundIndex]...)
	tmp = append(tmp, template.values[foundIndex+1:]...)

	easy.currentConnection.templates[templateIndex].values = tmp

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

func (easy *EasyEnv) isCurrentDBSet() error {

	if easy.currentConnection == nil {
		return fmt.Errorf("no database is currently open. Please open a database first using 'Open(path/to/sqlitefile)' before making any other calls")
	}

	return nil
}

func (easy *EasyEnv) GetProjectByID(projectID string) (int, *Project, error) {
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

func (easy *EasyEnv) GetTemplateByID(templateID string) (int, *Template, error) {
	foundIndex := 0
	var foundTemplate *Template
	for index, template := range easy.currentConnection.templates {
		if template.GetTemplateByID() == templateID {
			foundIndex = index
			foundTemplate = template
			return foundIndex, foundTemplate, nil
		}
	}
	return 0, foundTemplate, fmt.Errorf("no template found with ID %s. Please verify the ID and try again", templateID)
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

func createEnvString(environments []ProjectDataSet) string {
	var result string

	for _, env := range environments {
		result += fmt.Sprintf("%s=%s\n", env.keyName, env.value)
	}

	return result
}
