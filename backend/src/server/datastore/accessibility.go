package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/metadata"
)

type Accessibility struct {
	Id          int64          `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description" db:"description"`
	ImageUrl    sql.NullString `json:"imageUrl" db:"image_url"`
	meta        metadata.Metadata
}

func (a *Accessibility) SetExists() {
	a.meta.Exists = true
}

func (a *Accessibility) SetDeleted() {
	a.meta.Deleted = true
}

func (a *Accessibility) Exists() bool {
	return a.meta.Exists
}

func (a *Accessibility) Deleted() bool {
	return a.meta.Deleted
}

func (ds *Datastore) InsertAccessibility(a *Accessibility) error {
	var err error

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.accessibility (` +
		`name, description, image_url` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, a.Name, a.Description, a.ImageUrl).Scan(&a.Id)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) UpdateAccessibility(a *Accessibility) error {
	var err error

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.accessibility SET (` +
		`name, description, image_url` +
		`) = ( ` +
		`$1, $2, $3` +
		`) WHERE id = $4`

	_, err = ds.postgres.Exec(sql, a.Name, a.Description, a.ImageUrl, a.Id)
	return err
}

func (ds *Datastore) SaveAccessibility(a *Accessibility) error {
	if a.Exists() {
		return ds.UpdateAccessibility(a)
	}
	return ds.InsertAccessibility(a)
}

func (ds *Datastore) UpsertAccessibility(a *Accessibility) error {
	var err error

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.accessibility (` +
		`id, name, description, image_url` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, description, image_url` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.description, EXCLUDED.image_url` +
		`)`

	_, err = ds.postgres.Exec(sql, a.Id, a.Name, a.Description, a.ImageUrl)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) DeleteAccessibility(a *Accessibility) error {
	var err error

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.accessibility WHERE id = $1`

	_, err = ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return nil
}

func (ds *Datastore) GetAccessibilityById(id int64) (*Accessibility, error) {
	var err error

	const sql = `SELECT ` +
		`id, name, description, image_url ` +
		`FROM places4all.accessibility ` +
		`WHERE id = $1`

	a := Accessibility{}
	a.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&a.Id, &a.Name, &a.Description, &a.ImageUrl)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
