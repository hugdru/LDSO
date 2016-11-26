package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/metadata"
	"time"
)

type Criterion struct {
	Id            int64     `json:"id" db:"id"`
	IdSubgroup    int64     `json:"idSubgroup" db:"id_subgroup"`
	IdLegislation zero.Int  `json:"idLegislation" db:"id_legislation"`
	Name          string    `json:"name" db:"name"`
	Weight        int       `json:"weight" db:"weight"`
	CreatedDate   time.Time `json:"createdDate" db:"created_date"`

	// Objects
	Legislation     *Legislation     `json:"legislation,omitempty"`
	Accessibilities []*Accessibility `json:"accessibilities,omitempty"`

	meta metadata.Metadata
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

func ACriterion(allocateObjects bool) Criterion {
	criterion := Criterion{}
	if allocateObjects {
		criterion.Legislation = NewLegislation(true)
		criterion.Accessibilities = make([]*Accessibility, 0)
	}
	return criterion
}

func NewCriterion(allocateObjects bool) *Criterion {
	criterion := ACriterion(allocateObjects)
	return &criterion
}

func (ds *Datastore) InsertCriterion(c *Criterion) error {

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.criterion (` +
		`id_subgroup, id_legislation, name, weight, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING id`

	res, err := ds.postgres.Exec(sql, c.IdSubgroup, c.IdLegislation, c.Name, c.Weight, c.CreatedDate)
	if err != nil {
		return err
	}
	c.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}

	c.SetExists()
	return err
}

func (ds *Datastore) UpdateCriterion(c *Criterion) error {

	if !c.Exists() {
		return errors.New("update failed: does not exist")
	}

	if c.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.criterion SET (` +
		`id_subgroup, id_legislation, name, weight, created_date` +
		`) = ( ` +
		`$1, $2, $3, $4, $5` +
		`) WHERE id = $6`

	_, err := ds.postgres.Exec(sql, c.IdSubgroup, c.IdLegislation, c.Name, c.Weight, c.CreatedDate, c.Id)
	return err
}

func (ds *Datastore) SaveCriterion(c *Criterion) error {
	if c.Exists() {
		return ds.UpdateCriterion(c)
	}

	return ds.InsertCriterion(c)
}

func (ds *Datastore) UpsertCriterion(c *Criterion) error {

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.criterion (` +
		`id, id_subgroup, id_legislation, name, weight, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_subgroup, id_legislation, name, weight, created_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_subgroup, EXCLUDED.id_legislation, EXCLUDED.name, EXCLUDED.weight, EXCLUDED.created_date` +
		`)`

	_, err := ds.postgres.Exec(sql, c.Id, c.IdSubgroup, c.IdLegislation, c.Name, c.Weight, c.CreatedDate)
	if err != nil {
		return err
	}

	c.SetExists()

	return err
}

func (ds *Datastore) DeleteCriterion(c *Criterion) error {

	if !c.Exists() {
		return nil
	}

	if c.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.criterion WHERE id = $1`

	_, err := ds.postgres.Exec(sql, c.Id)
	if err != nil {
		return err
	}

	c.SetDeleted()

	return err
}

func (ds *Datastore) GetCriterionLegislation(c *Criterion) (*Legislation, error) {
	return ds.GetLegislationById(c.IdLegislation.Int64)
}

func (ds *Datastore) GetCriterionById(id int64) (*Criterion, error) {

	const sql = `SELECT ` +
		`id, id_subgroup, id_legislation, name, weight, created_date ` +
		`FROM places4all.criterion ` +
		`WHERE id = $1`

	c := ACriterion(false)
	c.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(&c)
	if err != nil {
		return nil, err
	}

	return &c, err
}

func (ds *Datastore) GetCriteriaBySubgroupId(idSubgroup int64) ([]*Criterion, error) {
	criteria := make([]*Criterion, 0)
	rows, err := ds.postgres.Queryx(
		`SELECT criterion.id, criterion.id_subgroup, criterion.id_legislation, criterion.name, criterion.weight,criterion.created_date `+
			`FROM places4all.criterion `+
			`WHERE criterion.id_subgroup = $1`, idSubgroup)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		criterion := NewCriterion(false)
		err := rows.StructScan(criterion)
		if err != nil {
			return nil, err
		}
		if criterion.IdLegislation.Valid {
			criterion.Legislation = NewLegislation(true)
			criterion.Legislation, err = ds.GetLegislationById(criterion.IdLegislation.Int64)
			if err != nil {
				return nil, err
			}
		}
		criteria = append(criteria, criterion)
		criterion.Accessibilities, err = ds.GetAccessibilitiesByCriterionId(criterion.Id)
		if err != nil {
			return nil, err
		}
	}

	return criteria, err
}
