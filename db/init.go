package db

import (
	"database/sql"
	"fmt"

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
