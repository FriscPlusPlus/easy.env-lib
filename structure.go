package easyenv

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type EasyEnv struct {
	connections       []*Connection
	currentConnection *Connection
}

type Connection struct {
	dbName    string     // db path
	db        *sql.DB    // db instance
	projects  []Project  // projects and the associated env data
	templates []Template // templates of all the envs
}

type Project struct {
	projectID   int
	projectName string
	path        string
	values      []DataSet
	method      string
}

type Template struct {
	templateID   int
	templateName string
	values       []DataSet
	method       string
}

type DataSet struct {
	keyName string
	value   string
	method  string
}
