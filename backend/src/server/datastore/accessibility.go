package datastore

import (
	"errors"
	"server/datastore/generators"
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
	err := ds.postgres.QueryRow(sql, a.Name).Scan(&a.Id)
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

func (ds *Datastore) DeleteAccessibilityById(idAccessibility int64) error {

	const sql = `DELETE FROM places4all.accessibility WHERE id = $1`

	_, err := ds.postgres.Exec(sql, idAccessibility)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) DeleteAccessibilitiesByCriterionId(idCriterion int64) error {
	const sql = `DELETE FROM places4all.criterion_accessibility ` +
		`WHERE criterion_accessibility.id_criterion = $1`

	_, err := ds.postgres.Exec(sql, idCriterion)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) DeleteAccessibilitiesByIds(idCriterion, idAccessibility int64) error {
	const sql = `DELETE FROM places4all.criterion_accessibility ` +
		`WHERE criterion_accessibility.id_criterion = $1 AND ` +
		`criterion_accessibility.id_accessibility = $2`

	_, err := ds.postgres.Exec(sql, idCriterion, idAccessibility)
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

func (ds *Datastore) GetAccessibilitiesByCriterionId(idCriterion int64) ([]*Accessibility, error) {
	rows, err := ds.postgres.Queryx(
		`SELECT accessibility.id, accessibility.name, criterion_accessibility.weight `+
			`FROM places4all.criterion_accessibility `+
			`JOIN places4all.accessibility ON accessibility.id = criterion_accessibility.id_accessibility `+
			`WHERE criterion_accessibility.id_criterion = $1`, idCriterion)
	if err != nil {
		return nil, err
	}

	accessibilities := make([]*Accessibility, 0)
	for rows.Next() {
		accessibility := NewAccessibility(false)
		accessibility.SetExists()
		err := rows.StructScan(accessibility)
		if err != nil {
			return nil, err
		}
		accessibilities = append(accessibilities, accessibility)
	}

	return accessibilities, err
}

func (ds *Datastore) GetAccessibilityByIds(idCriterion, idAccessibility int64) (*Accessibility, error) {
	row := ds.postgres.QueryRowx(
		`SELECT accessibility.id, accessibility.name, criterion_accessibility.weight `+
			`FROM places4all.criterion_accessibility `+
			`JOIN places4all.accessibility ON accessibility.id = criterion_accessibility.id_accessibility `+
			`WHERE criterion_accessibility.id_criterion = $1 AND `+
			`criterion_accessibility.id_accessibility = $2`, idCriterion, idAccessibility)

	accessibility := NewAccessibility(false)
	accessibility.SetExists()
	err := row.StructScan(accessibility)
	if err != nil {
		return nil, err
	}

	return accessibility, err
}

func (ds *Datastore) GetAccessibilities(limit, offset int, filter map[string]string) ([]*Accessibility, error) {

	where, values := generators.GenerateSearchClause(filter)

	sql := `SELECT ` +
		`id, name ` +
		`FROM places4all.accessibility ` +
		where
	sql = ds.postgres.Rebind(sql)

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

func (ds *Datastore) SaveAccessibilityByCriterionIdAccessibility(criterionId int64, accessibility *Accessibility) error {
	if !accessibility.Exists() {
		return errors.New("update failed: does not exist")
	}

	if accessibility.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.criterion_accessibility SET weight = $1 WHERE id_criterion = $2 AND id_accessibility = $3`

	_, err := ds.postgres.Exec(sql, accessibility.Weight, criterionId, accessibility.Id)
	return err
}

func (ds *Datastore) InsertAccessibilityByCriterionId(criterionId, accessibilityId int64, weight int) error {
	const sql = `INSERT INTO places4all.criterion_accessibility(id_criterion, id_accessibility, weight) VALUES($1, $2, $3)`

	_, err := ds.postgres.Exec(sql, criterionId, accessibilityId, weight)
	if err != nil {
		return err
	}

	return err

}
