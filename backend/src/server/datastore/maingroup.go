package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/metadata"
	"time"
)

type Maingroup struct {
	Id          int64          `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Weight      int            `json:"weight" db:"weight"`
	Description sql.NullString `json:"description" db:"description"`
	ImageUrl    sql.NullString `json:"imageUrl" db:"image_url"`
	Created     *time.Time     `json:"created" db:"created"`
	meta        metadata.Metadata
}

func (m *Maingroup) SetExists() {
	m.meta.Exists = true
}

func (m *Maingroup) SetDeleted() {
	m.meta.Deleted = true
}

func (m *Maingroup) Exists() bool {
	return m.meta.Exists
}

func (m *Maingroup) Deleted() bool {
	return m.meta.Deleted
}

func (ds *Datastore) InsertMaingroup(m *Maingroup) error {
	var err error

	if m.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.maingroup (` +
		`name, weight, description, image_url, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, m.Name, m.Weight, m.Description, m.ImageUrl, m.Created).Scan(&m.Id)
	if err != nil {
		return err
	}

	m.SetExists()

	return nil
}

func (ds *Datastore) UpdateMaingroup(m *Maingroup) error {
	var err error

	if !m.Exists() {
		return errors.New("update failed: does not exist")
	}

	if m.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.maingroup SET (` +
		`name, weight, description, image_url, created` +
		`) = ( ` +
		`$1, $2, $3, $4, $5` +
		`) WHERE id = $6`

	_, err = ds.postgres.Exec(sql, m.Name, m.Weight, m.Description, m.ImageUrl, m.Created, m.Id)
	return err
}

func (ds *Datastore) SaveMaingroup(m *Maingroup) error {
	if m.Exists() {
		return ds.UpdateMaingroup(m)
	}

	return ds.InsertMaingroup(m)
}

func (ds *Datastore) UpsertMaingroup(m *Maingroup) error {
	var err error

	if m.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.maingroup (` +
		`id, name, weight, description, image_url, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, weight, description, image_url, created` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.weight, EXCLUDED.description, EXCLUDED.image_url, EXCLUDED.created` +
		`)`

	_, err = ds.postgres.Exec(sql, m.Id, m.Name, m.Weight, m.Description, m.ImageUrl, m.Created)
	if err != nil {
		return err
	}

	m.SetExists()

	return nil
}

func (ds *Datastore) DeleteMaingroup(m *Maingroup) error {
	var err error

	if !m.Exists() {
		return nil
	}

	if m.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.maingroup WHERE id = $1`

	_, err = ds.postgres.Exec(sql, m.Id)
	if err != nil {
		return err
	}

	m.SetDeleted()

	return nil
}

func (ds *Datastore) GetMaingroupById(id int64) (*Maingroup, error) {
	var err error

	const sql = `SELECT ` +
		`id, name, weight, description, image_url, created ` +
		`FROM places4all.maingroup ` +
		`WHERE id = $1`

	m := Maingroup{}
	m.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&m.Id, &m.Name, &m.Weight, &m.Description, &m.ImageUrl, &m.Created)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
