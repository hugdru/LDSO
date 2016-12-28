package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
)

type Legislation struct {
	Id          int64       `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description zero.String `json:"description" db:"description"`
	Url         zero.String `json:"url" db:"url"`

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

func (l *Legislation) MustSet(name string) error {

	if name != "" {
		l.Name = name
	} else {
		return errors.New("name must be set")
	}

	return nil
}

func (l *Legislation) AllSetIfNotEmptyOrNil(name string, description string, url string) error {
	if name != "" {
		l.Name = name
	}

	return l.OptionalSetIfNotEmptyOrNil(description, url)
}

func (l *Legislation) OptionalSetIfNotEmptyOrNil(description, url string) error {
	if description != "" {
		l.Description = zero.StringFrom(description)
	}

	if url != "" {
		l.Url = zero.StringFrom(url)
	}

	return nil
}

func ALegislation(allocateObjects bool) Legislation {
	return Legislation{}
}

func NewLegislation(allocateObjects bool) *Legislation {
	legislation := ALegislation(allocateObjects)
	return &legislation
}

func (ds *Datastore) InsertLegislation(l *Legislation) error {

	if l == nil {
		return errors.New("legislation should not be nil")
	}

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.legislation (` +
		`name, description, url` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, l.Name, l.Description, l.Url).Scan(&l.Id)
	if err != nil {
		return err
	}

	l.SetExists()

	return err
}

func (ds *Datastore) UpdateLegislation(l *Legislation) error {

	if l == nil {
		return errors.New("legislation should not be nil")
	}

	if !l.Exists() {
		return errors.New("update failed: does not exist")
	}

	if l.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.legislation SET (` +
		`name, description, url` +
		`) = ( ` +
		`$1, $2, $3` +
		`) WHERE id = $4`

	_, err := ds.postgres.Exec(sql, l.Name, l.Description, l.Url, l.Id)
	return err
}

func (ds *Datastore) SaveLegislation(l *Legislation) error {
	if l.Exists() {
		return ds.UpdateLegislation(l)
	}

	return ds.InsertLegislation(l)
}

func (ds *Datastore) UpsertLegislation(l *Legislation) error {

	if l == nil {
		return errors.New("legislation should not be nil")
	}

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.legislation (` +
		`id, name, description, url` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, description, url` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.description, EXCLUDED.url` +
		`)`

	_, err := ds.postgres.Exec(sql, l.Id, l.Name, l.Description, l.Url)
	if err != nil {
		return err
	}

	l.SetExists()

	return err
}

func (ds *Datastore) DeleteLegislation(l *Legislation) error {

	if l == nil {
		return errors.New("legislation should not be nil")
	}

	if !l.Exists() {
		return nil
	}

	if l.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.legislation WHERE id = $1`

	_, err := ds.postgres.Exec(sql, l.Id)
	if err != nil {
		return err
	}

	l.SetDeleted()

	return err
}

func (ds *Datastore) DeleteLegislationById(id int64) error {

	const sql = `DELETE FROM places4all.legislation WHERE id = $1`

	_, err := ds.postgres.Exec(sql, id)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetLegislationById(id int64) (*Legislation, error) {

	const sql = `SELECT ` +
		`id, name, description, url ` +
		`FROM places4all.legislation ` +
		`WHERE id = $1`

	l := ALegislation(false)
	l.SetExists()

	err := ds.postgres.QueryRow(sql, id).Scan(&l.Id, &l.Name, &l.Description, &l.Url)
	if err != nil {
		return nil, err
	}

	return &l, err
}

func (ds *Datastore) GetLegislationByName(name string) (*Legislation, error) {

	const sql = `SELECT ` +
		`id, name, description, url ` +
		`FROM places4all.legislation ` +
		`WHERE name = $1`

	l := ALegislation(false)
	l.SetExists()

	err := ds.postgres.QueryRow(sql, name).Scan(&l.Id, &l.Name, &l.Description, &l.Url)
	if err != nil {
		return nil, err
	}

	return &l, err
}

func (ds *Datastore) GetLegislations(limit, offset int, filter map[string]interface{}) ([]*Legislation, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := `SELECT ` +
		`id, name, description, url ` +
		`FROM places4all.legislation ` +
		where +
		`ORDER BY legislation.id DESC LIMIT ` + strconv.Itoa(limit) +
		` OFFSET ` + strconv.Itoa(offset)
	sql = ds.postgres.Rebind(sql)

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	legislations := make([]*Legislation, 0)
	for rows.Next() {
		l := NewLegislation(true)
		l.SetExists()
		err = rows.StructScan(l)
		if err != nil {
			return nil, err
		}
		legislations = append(legislations, l)
	}

	return legislations, err
}
