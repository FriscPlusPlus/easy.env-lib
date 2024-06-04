package easyenv

import (
	"fmt"
)

func createTables(connection *Connection) error {

	db := connection.db
	_, err := db.Exec("CREATE TABLE projects(projectID INTEGER PRIMARY KEY AUTOINCREMENT, projectName TEXT, path TEXT)")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE templates(templateID INTEGER PRIMARY KEY AUTOINCREMENT, templateName TEXT)")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE templateValues(keyName TEXT PRIMARY KEY, templateID INTEGER, value TEXT, FOREIGN KEY(templateID) REFERENCES templates(templateID))")

	if err != nil {
		return err
	}
	return nil
}

func save(connection *Connection) error {

	err := saveProjects(connection)

	if err != nil {
		return err
	}

	err = saveTemplates(connection)

	if err != nil {
		return err
	}

	return nil
}

func saveProjects(connection *Connection) error {
	db := connection.db

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	for _, project := range connection.projects {

		switch project.method {
		case "INSERT":
			_, err := tx.Exec("INSERT INTO projects(projectName, path) VALUES(?, ?)", project.projectName, project.path)
			if err != nil {
				tx.Rollback()
				return err
			}
		case "UPDATE":
			_, err := tx.Exec("UPDATE projects SET projectName = ?, path = ? WHERE projectID = ?", project.projectName, project.path, project.projectID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func saveTemplates(connection *Connection) error {
	db := connection.db

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	for _, template := range connection.templates {

		switch template.method {
		case "INSERT":
			_, err := tx.Exec("INSERT INTO templates(templateName) VALUES(?)", template.templateName)
			if err != nil {
				tx.Rollback()
				return err
			}
		case "UPDATE":
			_, err := tx.Exec("UPDATE templates SET templateName = ? WHERE templateID = ?", template.templateName, template.templateID)
			if err != nil {
				tx.Rollback()
				return err
			}

		}

	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func removeData(connection *Connection, tableName string, parameterName string, id int) error {
	db := connection.db

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, parameterName)

	_, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
