package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/metadata"
)

type Legislation struct {
	Id   int64          `json:"id" db:"id"`
	Name sql.NullString `json:"name" db:"name"`
	Url  sql.NullString `json:"url" db:"url"`
	meta metadata.Metadata
}

func (l *Legislation) SetExists() {
	l.meta.Exists = true
}

func (l *Legislation) SetDeleted() {
	l.meta.Deleted = true
}

func (l *Legislation) Exists() bool {
	return l.meta.Exists
}

func (l *Legislation) Deleted() bool {
	return l.meta.Deleted
}

func (ds *Datastore) InsertLegislation(l *Legislation) error {
	var err error

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.legislation (` +
		`name, url` +
		`) VALUES (` +
		`$1, $2` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, l.Name, l.Url).Scan(&l.Id)
	if err != nil {
		return err
	}

	l.SetExists()

	return nil
}

func (ds *Datastore) UpdateLegislation(l *Legislation) error {
	var err error

	if !l.Exists() {
		return errors.New("update failed: does not exist")
	}

	if l.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.legislation SET (` +
		`name, url` +
		`) = ( ` +
		`$1, $2` +
		`) WHERE id = $3`

	_, err = ds.postgres.Exec(sql, l.Name, l.Url, l.Id)
	return err
}

func (ds *Datastore) SaveLegislation(l *Legislation) error {
	if l.Exists() {
		return ds.UpdateLegislation(l)
	}

	return ds.InsertLegislation(l)
}

func (ds *Datastore) UpsertLegislation(l *Legislation) error {
	var err error

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.legislation (` +
		`id, name, url` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, url` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.url` +
		`)`

	_, err = ds.postgres.Exec(sql, l.Id, l.Name, l.Url)
	if err != nil {
		return err
	}

	l.SetExists()

	return nil
}

func (ds *Datastore) DeleteLegislation(l *Legislation) error {
	var err error

	if !l.Exists() {
		return nil
	}

	if l.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.legislation WHERE id = $1`

	_, err = ds.postgres.Exec(sql, l.Id)
	if err != nil {
		return err
	}

	l.SetDeleted()

	return nil
}

func (ds *Datastore) GetLegislationById(id int64) (*Legislation, error) {
	var err error

	const sql = `SELECT ` +
		`id, name, url ` +
		`FROM places4all.legislation ` +
		`WHERE id = $1`

	l := Legislation{}
	l.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&l.Id, &l.Name, &l.Url)
	if err != nil {
		return nil, err
	}

	return &l, nil
}
