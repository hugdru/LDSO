package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/metadata"
	"time"
)

type Criterion struct {
	Id            int64          `json:"id" db:"id"`
	IdLegislation sql.NullInt64  `json:"idLegislation" db:"id_legislation"`
	Name          string         `json:"name" db:"name"`
	Weight        int            `json:"weight" db:"weight"`
	Description   sql.NullString `json:"description" db:"description"`
	ImageUrl      sql.NullString `json:"imageUrl" db:"image_url"`
	Created       *time.Time     `json:"created" db:"created"`
	meta          metadata.Metadata
}

func (c *Criterion) SetExists() {
	c.meta.Exists = true
}

func (c *Criterion) SetDeleted() {
	c.meta.Deleted = true
}

func (c *Criterion) Exists() bool {
	return c.meta.Exists
}

func (c *Criterion) Deleted() bool {
	return c.meta.Deleted
}

func (ds *Datastore) InsertCriterion(c *Criterion) error {
	var err error

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.criterion (` +
		`id_legislation, name, weight, description, image_url, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, c.IdLegislation, c.Name, c.Weight, c.Description, c.ImageUrl, c.Created).Scan(&c.Id)
	if err != nil {
		return err
	}

	c.SetExists()
	return nil
}

func (ds *Datastore) UpdateCriterion(c *Criterion) error {
	var err error

	if !c.Exists() {
		return errors.New("update failed: does not exist")
	}

	if c.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.criterion SET (` +
		`id_legislation, name, weight, description, image_url, created` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6` +
		`) WHERE id = $7`

	_, err = ds.postgres.Exec(sql, c.IdLegislation, c.Name, c.Weight, c.Description, c.ImageUrl, c.Created, c.Id)
	return err
}

func (ds *Datastore) SaveCriterion(c *Criterion) error {
	if c.Exists() {
		return ds.UpdateCriterion(c)
	}

	return ds.InsertCriterion(c)
}

func (ds *Datastore) UpsertCriterion(c *Criterion) error {
	var err error

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.criterion (` +
		`id, id_legislation, name, weight, description, image_url, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_legislation, name, weight, description, image_url, created` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_legislation, EXCLUDED.name, EXCLUDED.weight, EXCLUDED.description, EXCLUDED.image_url, EXCLUDED.created` +
		`)`

	_, err = ds.postgres.Exec(sql, c.Id, c.IdLegislation, c.Name, c.Weight, c.Description, c.ImageUrl, c.Created)
	if err != nil {
		return err
	}

	c.SetExists()

	return nil
}

func (ds *Datastore) DeleteCriterion(c *Criterion) error {
	var err error

	if !c.Exists() {
		return nil
	}

	if c.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.criterion WHERE id = $1`

	_, err = ds.postgres.Exec(sql, c.Id)
	if err != nil {
		return err
	}

	c.SetDeleted()

	return nil
}

func (ds *Datastore) GetCriterionLegislation(c *Criterion) (*Legislation, error) {
	return ds.GetLegislationById(c.IdLegislation.Int64)
}

func (ds *Datastore) GetCriterionById(id int64) (*Criterion, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_legislation, name, weight, description, image_url, created ` +
		`FROM places4all.criterion ` +
		`WHERE id = $1`

	c := Criterion{}
	c.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&c.Id, &c.IdLegislation, &c.Name, &c.Weight, &c.Description, &c.ImageUrl, &c.Created)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
