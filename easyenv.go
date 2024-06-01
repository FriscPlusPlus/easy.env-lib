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
	CreateNewDB(dbName string) (*Connection, error)
	AddProject(projectID, path string) error
	AddEnviormentToProject(projectID, key, value string) error
	AddEnviormentToTemplate(template, key, value string) error
	RemoveEnviormentFromProject(projectID, key, value string) error
	RemoveEnviormentFromTemplate(projectID, key, value string) error
	LoadTemplate(template string) error
	LoadAllTemplate() error
	SaveTemplate(template string) error
	SaveAllTemplate() error
	LoadProject(projectID string) error
	LoadAllProject() error
	SetEnviormentFromTemplate(template, projectID string) error
	SaveEnvForProject(projectID string) error
	SaveEnvForAllProjects() error
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
