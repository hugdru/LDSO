package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Accessibility struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`

	// Objects
	Weight int `json:"weight,omitempty" db:"weight"` // Only used when in relation to a certain criterion

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

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.accessibility(name) VALUES ($1) RETURNING id`
	result, err := ds.postgres.Exec(sql, a.Name)
	if err != nil {
		return err
	}
	a.Id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	a.SetExists()

	return err
}

func (ds *Datastore) UpdateAccessibility(a *Accessibility) error {

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

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.accessibility (id, name) ` +
		`VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET ` +
		`(id, name, description, image_url) = ` +
		`(EXCLUDED.id, EXCLUDED.name)`

	_, err := ds.postgres.Exec(sql, a.Id, a.Name)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) DeleteAccessibility(a *Accessibility) error {

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

func (ds *Datastore) GetAccessibilitiesByCriterionId(idCriterion int64) ([]*Accessibility, error) {
	accessibilities := make([]*Accessibility, 0)
	rows, err := ds.postgres.Queryx(
		`SELECT accessibility.id, accessibility.name, criterion_accessibility.weight `+
			`FROM places4all.criterion_accessibility `+
			`JOIN places4all.accessibility ON accessibility.id = criterion_accessibility.id_accessibility `+
			`WHERE criterion_accessibility.id_criterion = $1`, idCriterion)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		accessibility := NewAccessibility(false)
		err := rows.StructScan(accessibility)
		if err != nil {
			return nil, err
		}
		accessibilities = append(accessibilities, accessibility)
	}

	return accessibilities, err
}
