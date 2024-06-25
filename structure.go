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
	dbName    string      // db absolute path
	db        *sql.DB     // db instance
	projects  []*Project  // projects and the associated env data
	templates []*Template // templates of all the envs
}

type Project struct {
	projectID   string
	projectName string
	path        string // abosule path of the project containing the .env file
	values      []*DataSet
}

type Template struct {
	templateID   string
	templateName string
	values       []*DataSet
}

type DataSet struct {
	keyName string
	value   string
}
