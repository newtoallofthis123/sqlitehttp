package db

type RowInfo struct {
	Cid        int     `json:"cid"`
	Name       *string `json:"name"`
	Type       *string `json:"type"`
	Notnull    int     `json:"notnull"`
	DFlt_value *string `json:"dflt_value"`
	Pk         int     `json:"pk"`
}

// Some Utility Functions

func (d *Db) GetTableNames() []string {
	var tableNames []string
	for tableName := range d.RowsInfo {
		tableNames = append(tableNames, tableName)
	}
	return tableNames
}

func (d *Db) GetColumnsInfo(tableName string) map[string]string {
	var columnsInfo = make(map[string]string)
	for _, rowInfo := range d.RowsInfo[tableName] {
		columnsInfo[*rowInfo.Name] = *rowInfo.Type
	}
	return columnsInfo
}
