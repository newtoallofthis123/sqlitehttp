package db

import "regexp"

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (d *Db) CreateFillableMap(tableNames []string, columns []string, all bool) map[string]interface{} {
	var fillableMap = make(map[string]interface{})
	for _, tableName := range tableNames {
		for _, rowInfo := range d.RowsInfo[tableName] {
			if !all {
				if !contains(columns, *rowInfo.Name) {
					continue
				}
			}

			switch *rowInfo.Type {
			case "TEXT":
				fillableMap[*rowInfo.Name] = ""
			case "INTEGER":
				fillableMap[*rowInfo.Name] = 0
			case "REAL":
				fillableMap[*rowInfo.Name] = 0.0
			case "BLOB":
				fillableMap[*rowInfo.Name] = []byte{}
			case "NULL":
				fillableMap[*rowInfo.Name] = nil
			}
		}
	}
	return fillableMap
}

func getKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ParseRows parses the rows and table name from a query.
// This information is used to create a map of fillable
// values for the table
func (d *Db) parseRows(query string) ([]string, string) {
	// TODO: Add Support for joins and join parsing

	var rows []string
	var rawRows string
	var tableName string
	re := regexp.MustCompile(`(?i)SELECT\s+(.*?)\s+FROM\s+(\w+);`)
	matches := re.FindStringSubmatch(query)
	if len(matches) == 3 {
		rawRows = matches[1]
		tableName = matches[2]
	}

	if rawRows == "*" {
		// return getKeys(d.GetColumnsInfo(tableName)), tableName
	}

	re = regexp.MustCompile(`(?i)(\w+),?`)
	matches = re.FindAllString(rawRows, -1)
	for _, match := range matches {
		rows = append(rows, match)
	}

	return rows, tableName
}

func ParseTableNames(query string) []string {
	var tableNames []string
	// re := regexp.MustCompile(`(?i)\b(?:FROM|JOIN)\s+([a-zA-Z0-9_]+)`)
	re := regexp.MustCompile(`(?i)\b(?:FROM|JOIN)\s+([a-zA-Z0-9_]+)\b`)
	for _, match := range re.FindAllStringSubmatch(query, -1) {
		if len(match) > 1 {
			tableNames = append(tableNames, match[1])
		}
	}

	return tableNames
}
