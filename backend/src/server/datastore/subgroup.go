package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/metadata"
	"time"
)

type Subgroup struct {
	Id          int64          `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Weight      int            `json:"weight" db:"weight"`
	Description sql.NullString `json:"description" db:"description"`
	ImageUrl    sql.NullString `json:"imageUrl" db:"image_url"`
	Created     *time.Time     `json:"created" db:"created"`
	meta        metadata.Metadata
}

func (s *Subgroup) SetExists() {
	s.meta.Exists = true
}

func (s *Subgroup) SetDeleted() {
	s.meta.Deleted = true
}

func (s *Subgroup) Exists() bool {
	return s.meta.Exists
}

func (s *Subgroup) Deleted() bool {
	return s.meta.Deleted
}

func (ds *Datastore) InsertSubgroup(s *Subgroup) error {
	var err error

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.subgroup (` +
		`name, weight, description, image_url, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, s.Name, s.Weight, s.Description, s.ImageUrl, s.Created).Scan(&s.Id)
	if err != nil {
		return err
	}

	s.SetExists()

	return nil
}

func (ds *Datastore) UpdateSubgroup(s *Subgroup) error {
	var err error

	if !s.Exists() {
		return errors.New("update failed: does not exist")
	}

	if s.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.subgroup SET (` +
		`name, weight, description, image_url, created` +
		`) = ( ` +
		`$1, $2, $3, $4, $5` +
		`) WHERE id = $6`

	_, err = ds.postgres.Exec(sql, s.Name, s.Weight, s.Description, s.ImageUrl, s.Created, s.Id)
	return err
}

func (ds *Datastore) SaveSubgroup(s *Subgroup) error {
	if s.Exists() {
		return ds.UpdateSubgroup(s)
	}

	return ds.InsertSubgroup(s)
}

func (ds *Datastore) UpsertSubgroup(s *Subgroup) error {
	var err error

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.subgroup (` +
		`id, name, weight, description, image_url, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, weight, description, image_url, created` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.weight, EXCLUDED.description, EXCLUDED.image_url, EXCLUDED.created` +
		`)`

	_, err = ds.postgres.Exec(sql, s.Id, s.Name, s.Weight, s.Description, s.ImageUrl, s.Created)
	if err != nil {
		return err
	}

	s.SetExists()

	return nil
}

func (ds *Datastore) DeleteSubgroup(s *Subgroup) error {
	var err error

	if !s.Exists() {
		return nil
	}

	if s.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.subgroup WHERE id = $1`

	_, err = ds.postgres.Exec(sql, s.Id)
	if err != nil {
		return err
	}

	s.SetDeleted()

	return nil
}

func (ds *Datastore) GetSubgroupById(id int64) (*Subgroup, error) {
	var err error

	const sql = `SELECT ` +
		`id, name, weight, description, image_url, created ` +
		`FROM places4all.subgroup ` +
		`WHERE id = $1`

	s := Subgroup{}
	s.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&s.Id, &s.Name, &s.Weight, &s.Description, &s.ImageUrl, &s.Created)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
