package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
	"time"
)

type Criterion struct {
	Id            int64     `json:"id" db:"id"`
	IdSubgroup    int64     `json:"idSubgroup" db:"id_subgroup"`
	IdLegislation zero.Int  `json:"idLegislation" db:"id_legislation"`
	Name          string    `json:"name" db:"name"`
	Weight        int       `json:"weight" db:"weight"`
	Closed        zero.Bool `json:"closed" db:"closed"`
	CreatedDate   time.Time `json:"createdDate" db:"created_date"`

	// Objects
	Legislation string `json:"legislation, omitempty"`
	//Legislation            *Legislation              `json:"legislation, omitempty"`
	CriterionAccessibility []*CriterionAccessibility `json:"accessibilities,omitempty"`

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

func (c *Criterion) MustSet(idSubgroup int64, name string, weight int) error {

	if idSubgroup != 0 {
		c.IdSubgroup = idSubgroup
	} else {
		return errors.New("idSubgroup must be set")
	}
	if name != "" {
		c.Name = name
	} else {
		return errors.New("name must be set")
	}
	if weight != -1 {
		c.Weight = weight
	} else {
		return errors.New("weight must be set")
	}

	return nil
}

func (c *Criterion) AllSetIfNotEmptyOrNil(idSubgroup int64, name string, weight int, idLegislation int64) error {
	if idSubgroup != 0 {
		c.IdSubgroup = idSubgroup
	}
	if name != "" {
		c.Name = name
	}
	if weight != -1 {
		c.Weight = weight
	}

	return c.OptionalSetIfNotEmptyOrNil(idLegislation)
}

func (c *Criterion) OptionalSetIfNotEmptyOrNil(idLegislation int64) error {

	if idLegislation != 0 {
		c.IdLegislation = zero.IntFrom(idLegislation)
	}

	return nil
}

func (c *Criterion) UpdateSetIfNotEmptyOrNil(name string, weight int, idLegislation int64) error {
	if name != "" {
		c.Name = name
	}
	if weight != -1 {
		c.Weight = weight
	}
	if idLegislation != 0 {
		c.IdLegislation = zero.IntFrom(idLegislation)
	}

	return nil
}

func ACriterion(allocateObjects bool) Criterion {
	criterion := Criterion{}
	if allocateObjects {
		criterion.CriterionAccessibility = make([]*CriterionAccessibility, 0)
	}
	return criterion
}

func NewCriterion(allocateObjects bool) *Criterion {
	criterion := ACriterion(allocateObjects)
	return &criterion
}

func (ds *Datastore) InsertCriterion(c *Criterion) error {

	if c == nil {
		return errors.New("criterion should not be nil")
	}

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.criterion (` +
		`id_subgroup, id_legislation, name, weight, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, c.IdSubgroup, c.IdLegislation, c.Name, c.Weight, c.CreatedDate).Scan(&c.Id)
	if err != nil {
		return err
	}

	c.SetExists()
	return err
}

func (ds *Datastore) UpdateCriterion(c *Criterion) error {

	if c == nil {
		return errors.New("criterion should not be nil")
	}

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

	if c == nil {
		return errors.New("criterion should not be nil")
	}

	if c.Exists() {
		return ds.UpdateCriterion(c)
	}

	return ds.InsertCriterion(c)
}

func (ds *Datastore) UpsertCriterion(c *Criterion) error {

	if c == nil {
		return errors.New("criterion should not be nil")
	}

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

	if c == nil {
		return errors.New("criterion should not be nil")
	}

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

func (ds *Datastore) DeleteCriterionById(id int64) error {

	const sql = `DELETE FROM places4all.criterion WHERE id = $1`

	_, err := ds.postgres.Exec(sql, id)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetCriterionLegislation(c *Criterion) (*Legislation, error) {
	return ds.GetLegislationById(c.IdLegislation.Int64)
}

func (ds *Datastore) GetCriterionByIdWithLegislation(id int64) (*Criterion, error) {

	const sql = `SELECT ` +
		`id, id_subgroup, id_legislation, name, weight, closed, created_date ` +
		`FROM places4all.criterion ` +
		`WHERE id = $1`

	c := ACriterion(false)
	c.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(&c)
	if err != nil {
		return nil, err
	}
	if c.IdLegislation.Valid {
		legislation, err := ds.GetLegislationById(c.IdLegislation.Int64)
		if err != nil {
			return nil, err
		}
		c.Legislation = legislation.Name
	}
	/*
		if c.IdLegislation.Valid {
			c.Legislation, err = ds.GetLegislationById(c.IdLegislation.Int64)
			if err != nil {
				return nil, err
			}
		}
	*/
	return &c, err
}

func (ds *Datastore) GetCriteriaBySubgroupIdWithLegislationAndCriterionAccessibility(idSubgroup int64) ([]*Criterion, error) {
	criteria := make([]*Criterion, 0)
	rows, err := ds.postgres.Queryx(
		`SELECT criterion.id, criterion.id_subgroup, criterion.id_legislation, criterion.name, criterion.weight, criterion.closed, criterion.created_date `+
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
			legislation, err := ds.GetLegislationById(criterion.IdLegislation.Int64)
			if err != nil {
				return nil, err
			}
			criterion.Legislation = legislation.Name
		}
		/*
			if criterion.IdLegislation.Valid {
				criterion.Legislation, err = ds.GetLegislationById(criterion.IdLegislation.Int64)
				if err != nil {
					return nil, err
				}
			}
		*/
		criterion.CriterionAccessibility, err = ds.GetCriterionAccessibilitiesByCriterionId(criterion.Id)
		if err != nil {
			return nil, err
		}
		criteria = append(criteria, criterion)
	}

	return criteria, err
}

func (ds *Datastore) GetCriteriaWithLegislation(limit, offset int, filter map[string]interface{}) ([]*Criterion, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := `SELECT criterion.id, criterion.id_subgroup, criterion.id_legislation, criterion.name, criterion.weight, criterion.closed, criterion.created_date ` +
		`FROM places4all.criterion ` +
		where +
		`ORDER BY criterion.id DESC LIMIT ` + strconv.Itoa(limit) +
		` OFFSET ` + strconv.Itoa(offset)
	sql = ds.postgres.Rebind(sql)

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	criteria := make([]*Criterion, 0)
	for rows.Next() {
		criterion := NewCriterion(false)
		err := rows.StructScan(criterion)
		if err != nil {
			return nil, err
		}
		if criterion.IdLegislation.Valid {
			legislation, err := ds.GetLegislationById(criterion.IdLegislation.Int64)
			if err != nil {
				return nil, err
			}
			criterion.Legislation = legislation.Name
		}
		/*
			if criterion.IdLegislation.Valid {
				criterion.Legislation, err = ds.GetLegislationById(criterion.IdLegislation.Int64)
				if err != nil {
					return nil, err
				}
			}
		*/
		criteria = append(criteria, criterion)
	}

	return criteria, err
}

func (ds *Datastore) GetCriteriaByTemplateIdWithCriterionAccessibility(idTemplate int64) ([]*Criterion, error) {
	const sql = `SELECT criterion.id, criterion.id_subgroup, criterion.id_legislation, criterion.name, criterion.weight, criterion.closed, criterion.created_date ` +
		`FROM places4all.criterion ` +
		`JOIN places4all.subgroup ON subgroup.id = criterion.id_subgroup ` +
		`JOIN places4all.maingroup ON maingroup.id = subgroup.id_maingroup AND maingroup.id_template = $1`

	criteria := make([]*Criterion, 0)
	rows, err := ds.postgres.Queryx(sql, idTemplate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		criterion := NewCriterion(false)
		criterion.SetExists()
		err := rows.StructScan(criterion)
		if err != nil {
			return nil, err
		}
		if criterion.IdLegislation.Valid {
			legislation, err := ds.GetLegislationById(criterion.IdLegislation.Int64)
			if err != nil {
				return nil, err
			}
			criterion.Legislation = legislation.Name
		}
		/*
			if c.IdLegislation.Valid {
				c.Legislation, err = ds.GetLegislationById(c.IdLegislation.Int64)
				if err != nil {
					return nil, err
				}
			}
		*/
		criterion.CriterionAccessibility, err = ds.GetCriterionAccessibilitiesByCriterionId(criterion.Id)
		if err != nil {
			return nil, err
		}
		criteria = append(criteria, criterion)
	}

	return criteria, err
}
