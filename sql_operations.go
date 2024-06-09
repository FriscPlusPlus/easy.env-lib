package easyenv

import (
	"errors"
	"fmt"
	"sync"
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
	var err error
	var errorText string
	wg := new(sync.WaitGroup)

	var projectError error
	var templateError error

	wg.Add(1)
	wg.Add(1)

	go saveProjects(connection, &projectError, wg)
	go saveTemplates(connection, &templateError, wg)

	wg.Wait()

	var templateEnvError error

	wg.Add(1)

	go saveEnvTemplates(connection, &templateEnvError, wg)

	wg.Wait()

	if projectError != nil {
		errorText = fmt.Sprintf("An error occurred while saving the project. Details: %s\n", projectError.Error())
	}

	if templateError != nil {
		errorText = fmt.Sprintf("%sAn error occurred while saving the templates. Details: %s\n", errorText, templateError.Error())
	}

	if templateEnvError != nil {
		errorText = fmt.Sprintf("%sAn error occurred while saving the env in templates. details: %s\n", errorText, templateEnvError.Error())
	}

	if len(errorText) > 0 {
		err = errors.New(errorText)
	}

	return err
}

func saveProjects(connection *Connection, errorResult *error, wg *sync.WaitGroup) {

	defer wg.Done()
	db := connection.db

	tx, err := db.Begin()

	if err != nil {
		*errorResult = err
		return
	}
	for _, project := range connection.projects {

		switch project.method {
		case "INSERT":
			_, err := tx.Exec("INSERT INTO projects(projectName, path) VALUES(?, ?)", project.projectName, project.path)
			if err != nil {
				tx.Rollback()
				*errorResult = err
				return
			}
		case "UPDATE":
			_, err := tx.Exec("UPDATE projects SET projectName = ?, path = ? WHERE projectID = ?", project.projectName, project.path, project.projectID)
			if err != nil {
				tx.Rollback()
				*errorResult = err
				return
			}

		}
	}

	err = tx.Commit()
	if err != nil {
		*errorResult = err
		return
	}
}

func saveTemplates(connection *Connection, errorResult *error, wg *sync.WaitGroup) {
	defer wg.Done()

	db := connection.db

	tx, err := db.Begin()

	if err != nil {
		*errorResult = err
		return
	}
	for _, template := range connection.templates {

		switch template.method {
		case "INSERT":
			_, err := tx.Exec("INSERT INTO templates(templateName) VALUES(?)", template.templateName)
			if err != nil {
				tx.Rollback()
				*errorResult = err
				return
			}
		case "UPDATE":
			_, err := tx.Exec("UPDATE templates SET templateName = ? WHERE templateID = ?", template.templateName, template.templateID)
			if err != nil {
				tx.Rollback()
				*errorResult = err
				return
			}

		}
	}

	err = tx.Commit()

	if err != nil {
		*errorResult = err
	}
}

func saveEnvTemplates(connection *Connection, errorResult *error, wg *sync.WaitGroup) {
	defer wg.Done()

	db := connection.db

	tx, err := db.Begin()

	if err != nil {
		*errorResult = err
		return
	}

	for _, template := range connection.templates {
		for _, templateEnv := range template.values {

			switch templateEnv.method {
			case "INSERT":
				_, err := tx.Exec("INSERT INTO templateValues(keyName, templateID, value) VALUES(?, ?, ?)", templateEnv.keyName, templateEnv.templateID, templateEnv.value)
				if err != nil {
					tx.Rollback()
					*errorResult = err
					return
				}
			case "UPDATE":
				_, err := tx.Exec("UPDATE templateValues SET value = ? WHERE templateID = ? AND keyName = ?", templateEnv.value, templateEnv.templateID, templateEnv.keyName)
				if err != nil {
					tx.Rollback()
					*errorResult = err
					return
				}

			}

		}
	}

	err = tx.Commit()

	if err != nil {
		*errorResult = err
	}
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
