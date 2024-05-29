package test

import (
	"database/sql"
	"fmt"
)

type dbIntf interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

// DropAll is a testing utility that resets a database
func DropAll(db dbIntf) error {
	rows, err := db.Query(`select "name" from "sqlite_master" where "type" = 'table'`)
	if err != nil {
		return err
	}
	var dropStatements []string
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		if name == "sqlite_sequence" {
			continue
		}
		dropStatements = append(dropStatements, fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, name))
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, statement := range dropStatements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}
