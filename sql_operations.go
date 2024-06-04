package easyenv

func createTables(connection *Connection) error {

	db := connection.db
	_, err := db.Exec("CREATE TABLE projects(projectID TEXT, path TEXT, PRIMARY KEY(projectID))")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE templates(templateID INTEGER PRIMARY KEY, templateName TEXT)")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE templateValues(keyName TEXT PRIMARY KEY, templateID INTEGER, value TEXT, FOREIGN KEY(templateID) REFERENCES templates(templateID))")

	if err != nil {
		return err
	}

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
	for _, project := range connection.projects {

		if project.needSave {

			_, err := db.Exec("INSERT INTO projects(projectID, path) VALUES(?, ?) ON CONFLICT(projectID) DO UPDATE SET projectID = ?, path = ?", project.projectID, project.path, project.projectID, project.path)

			if err != nil {
				return err
			}
		}

	}
	return nil
}

func saveTemplates(connection *Connection) error {
	db := connection.db
	for _, template := range connection.templates {

		if template.needSave {

			_, err := db.Exec("INSERT INTO templates(templateID, templateName) VALUES(?, ?) ON CONFLICT(templateID) DO UPDATE SET templateID = ?, templateName = ?", template.templateID, template.templateName, template.templateID, template.templateName)

			if err != nil {
				return err
			}
		}

	}
	return nil
}
