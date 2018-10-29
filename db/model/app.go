// Package model contains the types for schema 'public'.
package model

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// App represents a row from 'public.apps'.
type App struct {
	ID        int64          `json:"id"`         // id
	Name      sql.NullString `json:"name"`       // name
	SecretKey sql.NullString `json:"secret_key"` // secret_key
	Status    sql.NullInt64  `json:"status"`     // status
	AccountID sql.NullInt64  `json:"account_id"` // account_id
	CreatedAt pq.NullTime    `json:"created_at"` // created_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the App exists in the database.
func (a *App) Exists() bool {
	return a._exists
}

// Deleted provides information if the App has been deleted from the database.
func (a *App) Deleted() bool {
	return a._deleted
}

// Insert inserts the App to the database.
func (a *App) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if a._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.apps (` +
		`name, secret_key, status, account_id, created_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, a.Name, a.SecretKey, a.Status, a.AccountID, a.CreatedAt)
	err = db.QueryRow(sqlstr, a.Name, a.SecretKey, a.Status, a.AccountID, a.CreatedAt).Scan(&a.ID)
	if err != nil {
		return err
	}

	// set existence
	a._exists = true

	return nil
}

// Update updates the App in the database.
func (a *App) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if a._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.apps SET (` +
		`name, secret_key, status, account_id, created_at` +
		`) = ( ` +
		`$1, $2, $3, $4, $5` +
		`) WHERE id = $6`

	// run query
	XOLog(sqlstr, a.Name, a.SecretKey, a.Status, a.AccountID, a.CreatedAt, a.ID)
	_, err = db.Exec(sqlstr, a.Name, a.SecretKey, a.Status, a.AccountID, a.CreatedAt, a.ID)
	return err
}

// Save saves the App to the database.
func (a *App) Save(db XODB) error {
	if a.Exists() {
		return a.Update(db)
	}

	return a.Insert(db)
}

// Upsert performs an upsert for App.
//
// NOTE: PostgreSQL 9.5+ only
func (a *App) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if a._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.apps (` +
		`id, name, secret_key, status, account_id, created_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, secret_key, status, account_id, created_at` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.secret_key, EXCLUDED.status, EXCLUDED.account_id, EXCLUDED.created_at` +
		`)`

	// run query
	XOLog(sqlstr, a.ID, a.Name, a.SecretKey, a.Status, a.AccountID, a.CreatedAt)
	_, err = db.Exec(sqlstr, a.ID, a.Name, a.SecretKey, a.Status, a.AccountID, a.CreatedAt)
	if err != nil {
		return err
	}

	// set existence
	a._exists = true

	return nil
}

// Delete deletes the App from the database.
func (a *App) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return nil
	}

	// if deleted, bail
	if a._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.apps WHERE id = $1`

	// run query
	XOLog(sqlstr, a.ID)
	_, err = db.Exec(sqlstr, a.ID)
	if err != nil {
		return err
	}

	// set deleted
	a._deleted = true

	return nil
}

// AppsQuery returns offset-limit rows from 'public.apps' filte by filter,
// ordered by "id" in descending order.
func AppFilter(db XODB, filter string, offset, limit int64) ([]*App, error) {
	sqlstr := `SELECT ` +
		`id, name, secret_key, status, account_id, created_at` +
		` FROM public.apps `

	if filter != "" {
		sqlstr = sqlstr + " WHERE " + filter
	}

	sqlstr = sqlstr + " order by id desc offset $1 limit $2"

	XOLog(sqlstr, offset, limit)
	q, err := db.Query(sqlstr, offset, limit)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*App
	for q.Next() {
		a := App{}

		// scan
		err = q.Scan(&a.ID, &a.Name, &a.SecretKey, &a.Status, &a.AccountID, &a.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}

	return res, nil
} // Account returns the Account associated with the App's AccountID (account_id).
//
// Generated from foreign key 'apps_account_id_fkey'.
func (a *App) Account(db XODB) (*Account, error) {
	return AccountByID(db, a.AccountID.Int64)
}

// AppsByAccountID retrieves a row from 'public.apps' as a App.
//
// Generated from index 'apps_account_id_idx'.
func AppsByAccountID(db XODB, accountID sql.NullInt64) ([]*App, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, secret_key, status, account_id, created_at ` +
		`FROM public.apps ` +
		`WHERE account_id = $1`

	// run query
	XOLog(sqlstr, accountID)
	q, err := db.Query(sqlstr, accountID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*App{}
	for q.Next() {
		a := App{
			_exists: true,
		}

		// scan
		err = q.Scan(&a.ID, &a.Name, &a.SecretKey, &a.Status, &a.AccountID, &a.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}

	return res, nil
}

// RetrieveAppByID retrieves a row from 'public.apps' as a App.
//
// Generated from index 'apps_pkey'.
func AppByID(db XODB, id int64) (*App, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, secret_key, status, account_id, created_at ` +
		`FROM public.apps ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	a := App{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&a.ID, &a.Name, &a.SecretKey, &a.Status, &a.AccountID, &a.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// AppBySecretKey retrieves a row from 'public.apps' as a App.
//
// Generated from index 'apps_secret_key_idx'.
func AppBySecretKey(db XODB, secretKey sql.NullString) (*App, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, secret_key, status, account_id, created_at ` +
		`FROM public.apps ` +
		`WHERE secret_key = $1`

	// run query
	XOLog(sqlstr, secretKey)
	a := App{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, secretKey).Scan(&a.ID, &a.Name, &a.SecretKey, &a.Status, &a.AccountID, &a.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
