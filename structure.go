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
	dbName    string     // db filename
	db        *sql.DB    // db instance
	projects  []Project  // projects and the associated env data
	templates []Template // templates of all the envs
}

type Project struct {
	projectID string
	path      string
	values    []DataSet
	new       bool
}

type Template struct {
	templateID string
	values     []DataSet
	new        bool
}

type DataSet struct {
	keyName   string
	value string
}
