package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/metadata"
	"time"
)

type Template struct {
	Id          int64          `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description" db:"description"`
	Created     *time.Time     `json:"created" db:"created"`
	meta        metadata.Metadata
}

func (t *Template) SetExists() {
	t.meta.Exists = true
}

func (t *Template) SetDeleted() {
	t.meta.Deleted = true
}

func (t *Template) Exists() bool {
	return t.meta.Exists
}

func (t *Template) Deleted() bool {
	return t.meta.Deleted
}

func (ds *Datastore) InsertTemplate(t *Template) error {
	var err error

	if t.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.template (` +
		`name, description, created` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, t.Name, t.Description, t.Created).Scan(&t.Id)
	if err != nil {
		return err
	}

	t.SetExists()

	return nil
}

func (ds *Datastore) UpdateTemplate(t *Template) error {
	var err error

	if !t.Exists() {
		return errors.New("update failed: does not exist")
	}

	if t.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.template SET (` +
		`name, description, created` +
		`) = ( ` +
		`$1, $2, $3` +
		`) WHERE id = $4`

	_, err = ds.postgres.Exec(sql, t.Name, t.Description, t.Created, t.Id)
	return err
}

func (ds *Datastore) SaveTemplate(t *Template) error {
	if t.Exists() {
		return ds.UpdateTemplate(t)
	}

	return ds.InsertTemplate(t)
}

func (ds *Datastore) UpsertTemplate(t *Template) error {
	var err error

	if t.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.template (` +
		`id, name, description, created` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, description, created` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.description, EXCLUDED.created` +
		`)`

	_, err = ds.postgres.Exec(sql, t.Id, t.Name, t.Description, t.Created)
	if err != nil {
		return err
	}

	t.SetExists()

	return nil
}

func (ds *Datastore) DeleteTemplate(t *Template) error {
	var err error

	if !t.Exists() {
		return nil
	}

	if t.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.template WHERE id = $1`

	_, err = ds.postgres.Exec(sql, t.Id)
	if err != nil {
		return err
	}

	t.SetDeleted()

	return nil
}

func (ds *Datastore) GetTemplateById(id int64) (*Template, error) {
	var err error

	const sql = `SELECT ` +
		`id, name, description, created ` +
		`FROM places4all.template ` +
		`WHERE id = $1`

	t := Template{}
	t.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&t.Id, &t.Name, &t.Description, &t.Created)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
