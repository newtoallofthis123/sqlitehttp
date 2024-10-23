package db

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	DbPath   string
	db       *sql.DB
	RowsInfo map[string][]RowInfo
}

const tableQuery = `SELECT name FROM sqlite_master WHERE type='table';`
const columnQuery = `PRAGMA table_info(%s);`

func NewDb(dbPath string) (*Db, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &Db{DbPath: dbPath, db: db}, nil
}

func (d *Db) Discover() error {
	tablesRows, err := d.db.Query(tableQuery)
	var rowsInfo = make(map[string][]RowInfo)
	for tablesRows.Next() {
		var table string
		err = tablesRows.Scan(&table)
		if err != nil {
			return err
		}
		rowsInfo[table] = make([]RowInfo, 0)
		columnsRows, err := d.db.Query(fmt.Sprintf(columnQuery, table))
		if err != nil {
			return err
		}
		for columnsRows.Next() {
			var dbInfo RowInfo
			err = columnsRows.Scan(&dbInfo.Cid, &dbInfo.Name, &dbInfo.Type, &dbInfo.Notnull, &dbInfo.DFlt_value, &dbInfo.Pk)
			if err != nil {
				return err
			}
			rowsInfo[table] = append(rowsInfo[table], dbInfo)
		}
	}
	d.RowsInfo = rowsInfo
	return nil
}

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

func (d *Db) RunQuery(query string) ([]map[string]interface{}, error) {
	var res = make([]map[string]interface{}, 0)

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		var placeholder interface{}
		values[i] = &placeholder
	}

	for rows.Next() {
		var rowsMap = make(map[string]interface{})
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		for i, value := range values {
			actualValue := *(value.(*interface{}))
			rowsMap[columns[i]] = actualValue
		}
		res = append(res, rowsMap)
	}

	return res, nil
}

func (d *Db) RunExec(query string) (sql.Result, error) {
	return d.db.Exec(query)
}
