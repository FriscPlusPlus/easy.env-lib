package easyenv

import (
	"database/sql"
	"errors"
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
	for _, connection := range easy.connections {
		if connection.dbName == dbName {
			return connection, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No connection found for the database with the name: %s", dbName))
}
