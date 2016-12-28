package datastore

import (
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
	"time"
)

type Subgroup struct {
	Id          int64     `json:"id" db:"id"`
	IdMaingroup int64     `json:"idMaingroup" db:"id_maingroup"`
	Name        string    `json:"name" db:"name"`
	Weight      int       `json:"weight" db:"weight"`
	CreatedDate time.Time `json:"createdDate" db:"created_date"`

	// Objects
	Criteria []*Criterion `json:"criteria,omitempty"`

	meta metadata.Metadata
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

func (s *Subgroup) MustSet(idMaingroup int64, name string, weight int) error {

	if idMaingroup != 0 {
		s.IdMaingroup = idMaingroup
	} else {
		return errors.New("idMaingroup must be set")
	}
	if name != "" {
		s.Name = name
	} else {
		return errors.New("name must be set")
	}
	if weight != -1 {
		s.Weight = weight
	} else {
		return errors.New("weight must be set")
	}

	return nil
}

func (s *Subgroup) AllSetIfNotEmptyOrNil(idMaingroup int64, name string, weight int) error {
	if idMaingroup != 0 {
		s.IdMaingroup = idMaingroup
	}
	if name != "" {
		s.Name = name
	}
	if weight != -1 {
		s.Weight = weight
	}

	return nil
}

func (s *Subgroup) UpdateSetIfNotEmptyOrNil(name string, weight int) error {
	if name != "" {
		s.Name = name
	}
	if weight != -1 {
		s.Weight = weight
	}

	return nil
}

func ASubgroup(allocateObjects bool) Subgroup {
	subgroup := Subgroup{}
	if allocateObjects {
		subgroup.Criteria = make([]*Criterion, 0)
	}
	return subgroup
}

func NewSubgroup(allocateObjects bool) *Subgroup {
	subgroup := ASubgroup(allocateObjects)
	return &subgroup
}

func (ds *Datastore) InsertSubgroup(s *Subgroup) error {

	if s == nil {
		return errors.New("subgroup should not be nil")
	}

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.subgroup (` +
		`id_maingroup, name, weight, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, s.IdMaingroup, s.Name, s.Weight, s.CreatedDate).Scan(&s.Id)
	if err != nil {
		return err
	}

	s.SetExists()

	return err
}

func (ds *Datastore) UpdateSubgroup(s *Subgroup) error {

	if s == nil {
		return errors.New("subgroup should not be nil")
	}

	if !s.Exists() {
		return errors.New("update failed: does not exist")
	}

	if s.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.subgroup SET (` +
		`id_maingroup, name, weight, created_date` +
		`) = ( ` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5`

	_, err := ds.postgres.Exec(sql, s.IdMaingroup, s.Name, s.Weight, s.CreatedDate, s.Id)
	return err
}

func (ds *Datastore) SaveSubgroup(s *Subgroup) error {
	if s.Exists() {
		return ds.UpdateSubgroup(s)
	}

	return ds.InsertSubgroup(s)
}

func (ds *Datastore) UpsertSubgroup(s *Subgroup) error {

	if s == nil {
		return errors.New("subgroup should not be nil")
	}

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.subgroup (` +
		`id, id_maingroup, name, weight, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_maingroup, name, weight, created_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_maingroup, EXCLUDED.name, EXCLUDED.weight, EXCLUDED.created_date` +
		`)`

	_, err := ds.postgres.Exec(sql, s.Id, s.IdMaingroup, s.Name, s.Weight, s.CreatedDate)
	if err != nil {
		return err
	}

	s.SetExists()

	return err
}

func (ds *Datastore) DeleteSubgroup(s *Subgroup) error {

	if s == nil {
		return errors.New("subgroup should not be nil")
	}

	if !s.Exists() {
		return nil
	}

	if s.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.subgroup WHERE id = $1`

	_, err := ds.postgres.Exec(sql, s.Id)
	if err != nil {
		return err
	}

	s.SetDeleted()

	return err
}

func (ds *Datastore) DeleteSubgroupById(id int64) error {

	const sql = `DELETE FROM places4all.subgroup WHERE id = $1`

	_, err := ds.postgres.Exec(sql, id)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetSubgroupById(id int64) (*Subgroup, error) {

	const sql = `SELECT ` +
		`id, id_maingroup, name, weight, created_date ` +
		`FROM places4all.subgroup ` +
		`WHERE id = $1`

	s := ASubgroup(false)
	s.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(&s)
	if err != nil {
		return nil, err
	}

	return &s, err
}

func (ds *Datastore) GetSubgroupsByMaingroupIdWithCriteria(idMaingroup int64) ([]*Subgroup, error) {
	subgroups := make([]*Subgroup, 0)
	rows, err := ds.postgres.Queryx(
		`SELECT subgroup.id, subgroup.id_maingroup, subgroup.name, subgroup.weight, subgroup.created_date `+
			`FROM places4all.subgroup `+
			`WHERE subgroup.id_maingroup = $1`, idMaingroup)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		subgroup := NewSubgroup(false)
		err := rows.StructScan(subgroup)
		if err != nil {
			return nil, err
		}
		subgroups = append(subgroups, subgroup)
		subgroup.Criteria, err = ds.GetCriteriaBySubgroupId(subgroup.Id)
		if err != nil {
			return nil, err
		}
	}

	return subgroups, err
}

func (ds *Datastore) GetSubgroups(limit, offset int, filter map[string]interface{}) ([]*Subgroup, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := `SELECT subgroup.id, subgroup.id_maingroup, subgroup.name, subgroup.weight, subgroup.created_date ` +
		`FROM places4all.subgroup ` +
		where +
		`ORDER BY subgroup.id DESC LIMIT ` + strconv.Itoa(limit) +
		` OFFSET ` + strconv.Itoa(offset)
	sql = ds.postgres.Rebind(sql)

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	subgroups := make([]*Subgroup, 0)
	for rows.Next() {
		subgroup := NewSubgroup(false)
		err := rows.StructScan(subgroup)
		if err != nil {
			return nil, err
		}
		subgroups = append(subgroups, subgroup)
	}

	return subgroups, err
}
