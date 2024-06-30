package easyenv

// constructor

func NewDataSet(keyName, value string) *DataSet {
	DataSet := new(DataSet)
	DataSet.keyName = keyName
	DataSet.SetValue(value)

	return DataSet
}

// getters

func (prjdta *DataSet) GetKey() string {
	return prjdta.keyName
}

func (prjdta *DataSet) GetValue() string {
	return prjdta.value
}

// setters
func (prjdta *DataSet) SetValue(value string) {
	prjdta.value = value
}

func (prjdta *DataSet) Remove() {
	prjdta.deleted = true
}
