package easyenv

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type EasyEnvDefinition interface {
	NewEasyEnv() *EasyEnv
	Load(dbName string) (*Connection, error)
	Open(dbName string) (*Connection, error)
	CloseDB(dbName string) error
	SaveDB(dbName string) error
	SaveCurrentDB() error
	SaveAllDB() error

	CreateNewDB(dbName string) (*Connection, error)

	AddProject(projectID, path string) error
	AddTemplate(template string) error

	RemoveProject(projectID, path string) error
	RemoveTemplate(template string) error

	AddEnvToProject(projectID, keyName, value string) error
	AddEnvToTemplate(template, keyName, value string) error

	RemoveEnvFromProject(projectID, keyName, value string) error
	RemoveEnvFromTemplate(projectID, keyName, value string) error

	LoadTemplate(template string) error
	LoadAllTemplate() error

	LoadProject(projectID string) error
	LoadAllProject() error

	LoadTemplateIntoProject(template, projectID string) error

	WriteEnvInProjectFile(projectID string) error
	WriteEnvInAllProjectsFile() error
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

func (easy *EasyEnv) SaveDB(dbName string) error {
	connection, err := easy.getConnectionByDBname(dbName)

	if err != nil {
		return err
	}

	err = save(connection)

	if err != nil {
		return err
	}

	return nil
}

func (easy *EasyEnv) SaveCurrentDB() error {

	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	err = save(easy.currentConnection)

	if err != nil {
		return err
	}

	return nil
}

func (easy *EasyEnv) AddProject(projectID, path string) error {

	err := easy.checkIfcurrentDBisSet()

	if err != nil {
		return err
	}

	var project Project
	project.projectID = projectID
	project.path = path
	project.needSave = true
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
	template.needSave = true
	easy.currentConnection.templates = append(easy.currentConnection.templates, template)

	return nil
}

// TODO: add check everywhere for currentdb, at least one db needs to be loaded

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
	tmpConnections := make([]*Connection, 0)
	foundIndex := 0
	for index, connection := range easy.connections {
		if connection.dbName == dbName {
			foundIndex = index
			break
		}
	}
	tmpConnections = append(tmpConnections, easy.connections[:foundIndex]...)
	tmpConnections = append(tmpConnections, easy.connections[foundIndex+1:]...)
	easy.connections = tmpConnections
}

func (easy *EasyEnv) checkIfcurrentDBisSet() error {

	if easy.currentConnection == nil {
		return fmt.Errorf("No open database was found, please open a databse first `Open(path/to/sqlitefile)` before doing any operation")
	}

	return nil
}
