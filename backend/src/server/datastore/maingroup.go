package datastore

import (
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
	"time"
)

type Maingroup struct {
	Id          int64     `json:"id" db:"id"`
	IdTemplate  int64     `json:"idTemplate" db:"id_template"`
	Name        string    `json:"name" db:"name"`
	Weight      int       `json:"weight" db:"weight"`
	CreatedDate time.Time `json:"createdDate" db:"created_date"`

	// Objects
	Subgroups []*Subgroup `json:"subgroups,omitempty"`

	meta metadata.Metadata
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

func (m *Maingroup) MustSet(idTemplate int64, name string, weight int) error {

	if idTemplate != 0 {
		m.IdTemplate = idTemplate
	} else {
		return errors.New("idTemplate must be set")
	}
	if name != "" {
		m.Name = name
	} else {
		return errors.New("name must be set")
	}
	if weight != -1 {
		m.Weight = weight
	} else {
		return errors.New("weight must be set")
	}

	return nil
}

func (m *Maingroup) AllSetIfNotEmptyOrNil(idTemplate int64, name string, weight int) error {
	if idTemplate != 0 {
		m.IdTemplate = idTemplate
	}
	if name != "" {
		m.Name = name
	}
	if weight != -1 {
		m.Weight = weight
	}

	return nil
}

func (m *Maingroup) UpdateSetIfNotEmptyOrNil(name string, weight int) error {
	if name != "" {
		m.Name = name
	}
	if weight != -1 {
		m.Weight = weight
	}

	return nil
}

func AMaingroup(allocateObjects bool) Maingroup {
	maingroup := Maingroup{}
	if allocateObjects {
		maingroup.Subgroups = make([]*Subgroup, 0)
	}
	return maingroup
}

func NewMaingroup(allocateObjects bool) *Maingroup {
	maingroup := AMaingroup(allocateObjects)
	return &maingroup
}

func (ds *Datastore) InsertMaingroup(m *Maingroup) error {

	if m == nil {
		return errors.New("maingroup should not be nil")
	}

	if m.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.maingroup (` +
		`id_template, name, weight, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, m.IdTemplate, m.Name, m.Weight, m.CreatedDate).Scan(&m.Id)
	if err != nil {
		return err
	}

	m.SetExists()

	return err
}

func (ds *Datastore) UpdateMaingroup(m *Maingroup) error {

	if m == nil {
		return errors.New("maingroup should not be nil")
	}

	if !m.Exists() {
		return errors.New("update failed: does not exist")
	}

	if m.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.maingroup SET (` +
		`id_template, name, weight, created_date` +
		`) = ( ` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5`

	_, err := ds.postgres.Exec(sql, m.IdTemplate, m.Name, m.Weight, m.CreatedDate, m.Id)
	return err
}

func (ds *Datastore) SaveMaingroup(m *Maingroup) error {
	if m.Exists() {
		return ds.UpdateMaingroup(m)
	}

	return ds.InsertMaingroup(m)
}

func (ds *Datastore) UpsertMaingroup(m *Maingroup) error {

	if m == nil {
		return errors.New("maingroup should not be nil")
	}

	if m.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.maingroup (` +
		`id, id_template, name, weight, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_template, name, weight, created_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_template, EXCLUDED.name, EXCLUDED.weight, EXCLUDED.created_date` +
		`)`

	_, err := ds.postgres.Exec(sql, m.Id, m.IdTemplate, m.Name, m.Weight, m.CreatedDate)
	if err != nil {
		return err
	}

	m.SetExists()

	return err
}

func (ds *Datastore) DeleteMaingroup(m *Maingroup) error {

	if m == nil {
		return errors.New("maingroup should not be nil")
	}

	if !m.Exists() {
		return nil
	}

	if m.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.maingroup WHERE id = $1`

	_, err := ds.postgres.Exec(sql, m.Id)
	if err != nil {
		return err
	}

	m.SetDeleted()

	return err
}

func (ds *Datastore) DeleteMaingroupById(id int64) error {

	const sql = `DELETE FROM places4all.maingroup WHERE id = $1`

	_, err := ds.postgres.Exec(sql, id)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetMaingroupById(id int64) (*Maingroup, error) {

	const sql = `SELECT ` +
		`id, id_template, name, weight, created_date ` +
		`FROM places4all.maingroup ` +
		`WHERE id = $1`

	m := AMaingroup(false)
	m.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(&m)
	if err != nil {
		return nil, err
	}

	return &m, err
}

func (ds *Datastore) GetMaingroupsByTemplateIdWithSubgroups(idTemplate int64) ([]*Maingroup, error) {
	maingroups := make([]*Maingroup, 0)
	rows, err := ds.postgres.Queryx(
		`SELECT maingroup.id, maingroup.id_template, maingroup.name, maingroup.weight, maingroup.created_date `+
			`FROM places4all.maingroup `+
			`WHERE maingroup.id_template = $1`, idTemplate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		maingroup := NewMaingroup(false)
		err := rows.StructScan(maingroup)
		if err != nil {
			return nil, err
		}
		maingroups = append(maingroups, maingroup)
		maingroup.Subgroups, err = ds.GetSubgroupsByMaingroupId(maingroup.Id)
		if err != nil {
			return nil, err
		}
	}

	return maingroups, err
}

func (ds *Datastore) GetMaingroups(limit, offset int, filter map[string]interface{}) ([]*Maingroup, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := `SELECT maingroup.id, maingroup.id_template, maingroup.name, maingroup.weight, maingroup.created_date ` +
		`FROM places4all.maingroup ` +
		where +
		`ORDER BY maingroup.id DESC LIMIT ` + strconv.Itoa(limit) +
		` OFFSET ` + strconv.Itoa(offset)
	sql = ds.postgres.Rebind(sql)

	maingroups := make([]*Maingroup, 0)
	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		maingroup := NewMaingroup(false)
		err := rows.StructScan(maingroup)
		if err != nil {
			return nil, err
		}
		maingroups = append(maingroups, maingroup)
	}

	return maingroups, err
}
