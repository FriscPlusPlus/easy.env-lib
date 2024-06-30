package easyenv

import (
	"fmt"

	"github.com/google/uuid"
)

func NewTemplate(templateName string) *Template {
	template := new(Template)

	template.templateID = uuid.NewString()
	template.SetTemplateName(templateName)

	return template
}

// getters

func (template *Template) GetTemplateID() string {
	return template.templateID
}

func (template *Template) GetTemplateName() string {
	return template.templateName
}

func (template *Template) GetEnvironments() []*DataSet {
	return template.values
}

func (template *Template) GetEnvironmentByKey(keyName string) (*DataSet, error) {

	for _, env := range template.values {
		if env.keyName == keyName {
			return env, nil
		}
	}

	return nil, fmt.Errorf("no enviorment found with the key %s", keyName)
}

// setters
func (template *Template) SetTemplateName(templateName string) {
	template.templateName = templateName
}

func (template *Template) AddEnvrioment(keyName, value string) (*DataSet, error) {

	_, ok := template.GetEnvironmentByKey(keyName)

	if ok == nil {
		return nil, fmt.Errorf("an enviorment with the key %s already exists", keyName)
	}

	env := NewDataSet(keyName, value)
	template.values = append(template.values, env)
	return env, nil
}

func (template *Template) Remove() {
	template.deleted = true

	for _, data := range template.values {
		data.Remove()
	}
}

func (template *Template) RemoveEnviorment(keyName string) {

	tmp := make([]*DataSet, 0)
	foundIndex := 0
	for index, env := range template.values {
		if env.keyName == keyName {
			foundIndex = index
			break
		}
	}
	tmp = append(tmp, template.values[:foundIndex]...)
	tmp = append(tmp, template.values[foundIndex+1:]...)

	template.values = tmp
}

// this method will remove the enviorment
func (template *Template) RemoveAllEnviorments() {
	template.values = []*DataSet{}
}
