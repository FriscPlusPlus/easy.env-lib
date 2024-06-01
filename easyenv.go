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

	CreateNewDB(dbName string) (*Connection, error)

	AddProject(projectID, path string)
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

	if easy.currentConnection.dbName == dbName {
		easy.currentConnection = nil
	}

	if err != nil {
		return err
	}

	err = connection.db.Close()

	if err != nil {
		return err
	}

	easy.removeConnection(dbName)

	return nil
}

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

func (easy *EasyEnv) CreateNewDB(dbName string) (*Connection, error) {
	connection, err := easy.Load(dbName)

	if err != nil {
		return nil, err
	}

	db := connection.db
	_, err = db.Exec("CREATE TABLE projects(projectID TEXT, path TEXT, PRIMARY KEY(projectID))")

	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE templates(templateID TEXT, keyName TEXT, value TEXT, PRIMARY KEY(templateID, keyName))")

	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (easy *EasyEnv) AddProject(projectID, path string) {
	var project Project
	project.projectID = projectID
	project.path = path
	project.new = true
	easy.currentConnection.projects = append(easy.currentConnection.projects, project)
}

/*func (easy *EasyEnv) SaveDB(dbName string) error {
	db := easy.currentConnection.db
	for _, project := range easy.currentConnection.projects {
	}
}*/

// TODO: add check everywhere for currentdb, at least one db needs to be loaded
// TODO: remove the connection return, the instance is already in the easyenv
