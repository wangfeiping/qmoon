// Package model contains the types for schema 'public'.
package model

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// Genesi represents a row from 'public.genesis'.
type Genesi struct {
	ID          int64          `json:"id"`           // id
	ChainID     sql.NullString `json:"chain_id"`     // chain_id
	GenesisTime pq.NullTime    `json:"genesis_time"` // genesis_time
	Data        sql.NullString `json:"data"`         // data
	CreatedAt   pq.NullTime    `json:"created_at"`   // created_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Genesi exists in the database.
func (g *Genesi) Exists() bool {
	return g._exists
}

// Deleted provides information if the Genesi has been deleted from the database.
func (g *Genesi) Deleted() bool {
	return g._deleted
}

// Insert inserts the Genesi to the database.
func (g *Genesi) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if g._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.genesis (` +
		`chain_id, genesis_time, data, created_at` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, g.ChainID, g.GenesisTime, g.Data, g.CreatedAt)
	err = db.QueryRow(sqlstr, g.ChainID, g.GenesisTime, g.Data, g.CreatedAt).Scan(&g.ID)
	if err != nil {
		return err
	}

	// set existence
	g._exists = true

	return nil
}

// Update updates the Genesi in the database.
func (g *Genesi) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !g._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if g._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.genesis SET (` +
		`chain_id, genesis_time, data, created_at` +
		`) = ( ` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5`

	// run query
	XOLog(sqlstr, g.ChainID, g.GenesisTime, g.Data, g.CreatedAt, g.ID)
	_, err = db.Exec(sqlstr, g.ChainID, g.GenesisTime, g.Data, g.CreatedAt, g.ID)
	return err
}

// Save saves the Genesi to the database.
func (g *Genesi) Save(db XODB) error {
	if g.Exists() {
		return g.Update(db)
	}

	return g.Insert(db)
}

// Upsert performs an upsert for Genesi.
//
// NOTE: PostgreSQL 9.5+ only
func (g *Genesi) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if g._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.genesis (` +
		`id, chain_id, genesis_time, data, created_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, chain_id, genesis_time, data, created_at` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.chain_id, EXCLUDED.genesis_time, EXCLUDED.data, EXCLUDED.created_at` +
		`)`

	// run query
	XOLog(sqlstr, g.ID, g.ChainID, g.GenesisTime, g.Data, g.CreatedAt)
	_, err = db.Exec(sqlstr, g.ID, g.ChainID, g.GenesisTime, g.Data, g.CreatedAt)
	if err != nil {
		return err
	}

	// set existence
	g._exists = true

	return nil
}

// Delete deletes the Genesi from the database.
func (g *Genesi) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !g._exists {
		return nil
	}

	// if deleted, bail
	if g._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.genesis WHERE id = $1`

	// run query
	XOLog(sqlstr, g.ID)
	_, err = db.Exec(sqlstr, g.ID)
	if err != nil {
		return err
	}

	// set deleted
	g._deleted = true

	return nil
}

// GenesisQuery returns offset-limit rows from 'public.genesis' filte by filter,
// ordered by "id" in descending order.
func GenesiFilter(db XODB, filter string, offset, limit int64) ([]*Genesi, error) {
	sqlstr := `SELECT ` +
		`id, chain_id, genesis_time, data, created_at` +
		` FROM public.genesis `

	if filter != "" {
		sqlstr = sqlstr + " WHERE " + filter
	}

	if limit > 0 {
		sqlstr = sqlstr + fmt.Sprintf(" offset %d limit %d", offset, limit)
	}

	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Genesi
	for q.Next() {
		g := Genesi{}

		// scan
		err = q.Scan(&g.ID, &g.ChainID, &g.GenesisTime, &g.Data, &g.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &g)
	}

	return res, nil
}

// GenesiByChainID retrieves a row from 'public.genesis' as a Genesi.
//
// Generated from index 'genesis_chain_id_idx'.
func GenesiByChainID(db XODB, chainID sql.NullString) (*Genesi, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, chain_id, genesis_time, data, created_at ` +
		`FROM public.genesis ` +
		`WHERE chain_id = $1`

	// run query
	XOLog(sqlstr, chainID)
	g := Genesi{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, chainID).Scan(&g.ID, &g.ChainID, &g.GenesisTime, &g.Data, &g.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// GenesiByID retrieves a row from 'public.genesis' as a Genesi.
//
// Generated from index 'genesis_pkey'.
func GenesiByID(db XODB, id int64) (*Genesi, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, chain_id, genesis_time, data, created_at ` +
		`FROM public.genesis ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	g := Genesi{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&g.ID, &g.ChainID, &g.GenesisTime, &g.Data, &g.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &g, nil
}
