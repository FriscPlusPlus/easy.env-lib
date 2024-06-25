package easyenv

import (
	"errors"
	"fmt"
	"sync"
)

func createTables(connection *Connection) error {

	db := connection.db
	_, err := db.Exec("CREATE TABLE projects(projectID TEXT PRIMARY KEY, projectName TEXT, path TEXT)")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE templates(templateID TEXT PRIMARY KEY, templateName TEXT)")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE templateValues(keyName TEXT PRIMARY KEY, templateID stringEGER, value TEXT, FOREIGN KEY(templateID) REFERENCES templates(templateID))")

	if err != nil {
		return err
	}
	return nil
}

func saveDataInDB(connection *Connection) error {
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
		query := "INSERT INTO projects(projectID, projectName, path) VALUES(?, ?, ?) ON CONFLICT(projectID) DO UPDATE SET projectName = ?, path = ? WHERE projectID = ?"
		_, err := tx.Exec(query, project.projectID, project.projectName, project.path, project.projectName, project.path, project.projectID)
		if err != nil {
			tx.Rollback()
			*errorResult = err
			return
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

		query := "INSERT INTO templates(templateID, templateName) VALUES(?, ?) ON CONFLICT(templateID) DO UPDATE SET templateName = ? WHERE templateID = ?"
		_, err := tx.Exec(query, template.templateID, template.templateName, template.templateName, template.templateID)
		if err != nil {
			tx.Rollback()
			*errorResult = err
			return
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

			_, err := tx.Exec("INSERT INTO templateValues(keyName, templateID, value) VALUES(?, ?, ?) ON CONFLICT(keyName, templateID) DO UPDATE SET value = ? WHERE keyName = ? AND templateID = ?", templateEnv.keyName, templateEnv.templateID, templateEnv.value, templateEnv.value, templateEnv.keyName, template.templateID)
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

func removeData(connection *Connection, tableName string, parameterName string, id string) error {
	db := connection.db

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, parameterName)

	_, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func removeTemplateEnvData(connection *Connection, templateID string, keyName string) error {
	db := connection.db

	_, err := db.Exec("DELETE FROM templateValues WHERE templateID = ? AND keyName = ?", templateID, keyName)

	if err != nil {
		return err
	}

	return nil
}
