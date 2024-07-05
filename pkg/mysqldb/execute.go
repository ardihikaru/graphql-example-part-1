package mysqldb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/developersismedika/sqlx"
)

// Queryx executes query
func (d *Storage) Queryx(query string, args interface{}) (*sqlx.Rows, error) {
	var err error
	var rows *sqlx.Rows
	start := time.Now()

	if args == nil {
		rows, err = d.Db.Queryx(query)
	} else {
		rows, err = d.Db.Queryx(query, args)
	}

	// TODO: publish if exists

	// logs query execution times
	d.Log.Debug(fmt.Sprintf("Scan Execution took %s", time.Since(start)))

	return rows, err
}

// Exec executes a single query with a transaction
func (d *Storage) Exec(query string, args interface{}) (*sql.Result, error) {
	var err error
	start := time.Now()

	tx := d.Db.MustBegin()

	rslt, err := tx.Exec(query, args)
	if err != nil {
		d.Log.Warn(fmt.Sprintf("query execution failed. rolling back the changes: %s", err.Error()))
		_ = tx.Rollback()

		return nil, err
	}

	// debug only: print results
	d.insertOnePrint(rslt)

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// TODO: publish if exists

	// logs query execution times
	d.Log.Debug(fmt.Sprintf("Scan Execution took %s", time.Since(start)))

	return &rslt, nil
}

// NamedExec executes a single query with a transaction
func (d *Storage) NamedExec(qArgs QueryArgs) (*sql.Result, error) {
	var err error
	start := time.Now()

	tx := d.Db.MustBegin()

	rslt, err := tx.NamedExec(qArgs.Query, qArgs.Args)
	if err != nil {
		d.Log.Warn(fmt.Sprintf("query execution failed. rolling back the changes: %s", err.Error()))
		_ = tx.Rollback()

		return nil, err
	}

	// debug only: print results
	d.insertOnePrint(rslt)

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// TODO: publish if exists

	// logs query execution times
	d.Log.Debug(fmt.Sprintf("Scan Execution took %s", time.Since(start)))

	return &rslt, nil
}

// ExecMany executes multiple queries with a transaction
func (d *Storage) ExecMany(qArgsList []QueryArgs) error {
	var err error
	start := time.Now()

	tx := d.Db.MustBegin()

	for k, _ := range qArgsList {
		exec, err := tx.NamedExec(qArgsList[k].Query, qArgsList[k].Args)
		if err != nil {
			d.Log.Warn(fmt.Sprintf("query execution failed. rolling back the changes"))
			_ = tx.Rollback()

			return err
		}

		// debug only: print results
		d.insertOnePrint(exec)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// TODO: publish if exists

	// logs query execution times
	d.Log.Debug(fmt.Sprintf("Scan Execution took %s", time.Since(start)))

	return nil
}

// insertOnePrint prints out the query result
func (d *Storage) insertOnePrint(exec sql.Result) {
	// prints log on debug mode only
	lastInsertedId, err := exec.LastInsertId()
	if err != nil {
		d.Log.Debug(fmt.Sprintf("last inserted ID: %v", lastInsertedId))
	}

	// prints log on debug mode only
	rowsAffected, err := exec.RowsAffected()
	if err != nil {
		d.Log.Debug(fmt.Sprintf("Number of row affected: %v", rowsAffected))
	}
}
