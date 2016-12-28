package datastore

import (
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
)

type Accessibility struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`

	meta metadata.Metadata
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

func (a *Accessibility) MustSet(name string) error {
	if name != "" {
		a.Name = name
	} else {
		return errors.New("name must be set")
	}
	return nil
}

func (a *Accessibility) AllSetIfNotEmptyOrNil(name string) error {
	if name != "" {
		a.Name = name
	}
	return nil
}

func AAccessibility(allocateObjects bool) Accessibility {
	accessibility := Accessibility{}
	//if allocateObjects {
	//}
	return accessibility
}

func NewAccessibility(allocateObjects bool) *Accessibility {
	accessibility := AAccessibility(allocateObjects)
	return &accessibility
}

func (ds *Datastore) InsertAccessibility(a *Accessibility) error {

	if a == nil {
		return errors.New("accessibility should not be nil")
	}

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.accessibility(name) VALUES ($1) RETURNING id`
	err := ds.postgres.QueryRow(sql, a.Name).Scan(&a.Id)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) UpdateAccessibility(a *Accessibility) error {

	if a == nil {
		return errors.New("accessibility should not be nil")
	}

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.accessibility SET name = $1 WHERE id = $2`

	_, err := ds.postgres.Exec(sql, a.Name, a.Id)
	return err
}

func (ds *Datastore) SaveAccessibility(a *Accessibility) error {
	if a.Exists() {
		return ds.UpdateAccessibility(a)
	}
	return ds.InsertAccessibility(a)
}

func (ds *Datastore) UpsertAccessibility(a *Accessibility) error {

	if a == nil {
		return errors.New("accessibility should not be nil")
	}

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.accessibility (id, name) ` +
		`VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET ` +
		`(id, name) = ` +
		`(EXCLUDED.id, EXCLUDED.name)`

	_, err := ds.postgres.Exec(sql, a.Id, a.Name)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) DeleteAccessibility(a *Accessibility) error {

	if a == nil {
		return errors.New("accessibility should not be nil")
	}

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.accessibility WHERE id = $1`

	_, err := ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return err
}

func (ds *Datastore) DeleteAccessibilityById(idAccessibility int64) error {

	const sql = `DELETE FROM places4all.accessibility WHERE id = $1`

	_, err := ds.postgres.Exec(sql, idAccessibility)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetAccessibilityById(id int64) (*Accessibility, error) {

	const sql = `SELECT id, name FROM places4all.accessibility WHERE id = $1`

	a := AAccessibility(true)
	a.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(&a)
	if err != nil {
		return nil, err
	}

	return &a, err
}

func (ds *Datastore) GetAccessibilities(limit, offset int, filter map[string]interface{}) ([]*Accessibility, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := ds.postgres.Rebind(`SELECT ` +
		`id, name ` +
		`FROM places4all.accessibility ` +
		where)

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	accessibilities := make([]*Accessibility, 0)
	for rows.Next() {
		a := NewAccessibility(true)
		a.SetExists()
		err = rows.StructScan(a)
		if err != nil {
			return nil, err
		}
		accessibilities = append(accessibilities, a)
	}

	return accessibilities, err
}
