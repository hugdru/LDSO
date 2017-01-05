package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type CriterionAccessibility struct {
	IdCriterion     int64  `json:"-" db:"id_criterion"`
	IdAccessibility int64  `json:"id" db:"id_accessibility"`
	Name            string `json:"name" db:"name"`
	Weight          int    `json:"weight" db:"weight"`
	Closed          bool   `json:"closed" db:"closed"`

	meta metadata.Metadata
}

func (a *CriterionAccessibility) SetExists() {
	a.meta.Exists = true
}

func (a *CriterionAccessibility) SetDeleted() {
	a.meta.Deleted = true
}

func (a *CriterionAccessibility) Exists() bool {
	return a.meta.Exists
}

func (a *CriterionAccessibility) Deleted() bool {
	return a.meta.Deleted
}

func ACriterionAccessibility(allocateObjects bool) CriterionAccessibility {
	criterionAccessibility := CriterionAccessibility{}
	//if allocateObjects {
	//}
	return criterionAccessibility
}

func NewCriterionAccessibility(allocateObjects bool) *CriterionAccessibility {
	criterionAccessibility := ACriterionAccessibility(allocateObjects)
	return &criterionAccessibility
}

func (ds *Datastore) DeleteCriterionAccessibilitiesByCriterionId(idCriterion int64) error {
	const sql = `DELETE FROM places4all.criterion_accessibility ` +
		`WHERE criterion_accessibility.id_criterion = $1`

	_, err := ds.postgres.Exec(sql, idCriterion)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) DeleteCriterionAccessibilityByIds(idCriterion, idAccessibility int64) error {
	const sql = `DELETE FROM places4all.criterion_accessibility ` +
		`WHERE criterion_accessibility.id_criterion = $1 AND ` +
		`criterion_accessibility.id_accessibility = $2`

	_, err := ds.postgres.Exec(sql, idCriterion, idAccessibility)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetCriterionAccessibilitiesByCriterionId(idCriterion int64) ([]*CriterionAccessibility, error) {
	rows, err := ds.postgres.Queryx(
		`SELECT criterion_accessibility.id_criterion, `+
			`criterion_accessibility.id_accessibility, accessibility.name, criterion_accessibility.weight, criterion_accessibility.closed `+
			`FROM places4all.criterion_accessibility `+
			`JOIN places4all.accessibility ON accessibility.id = criterion_accessibility.id_accessibility `+
			`WHERE criterion_accessibility.id_criterion = $1`, idCriterion)
	if err != nil {
		return nil, err
	}

	criteriaAccessibilities := make([]*CriterionAccessibility, 0)
	for rows.Next() {
		criterionAccessibility := NewCriterionAccessibility(false)
		criterionAccessibility.SetExists()
		err := rows.StructScan(criterionAccessibility)
		if err != nil {
			return nil, err
		}
		criteriaAccessibilities = append(criteriaAccessibilities, criterionAccessibility)
	}

	return criteriaAccessibilities, err
}

func (ds *Datastore) GetCriterionAccessibilityByIds(idCriterion, idAccessibility int64) (*CriterionAccessibility, error) {
	row := ds.postgres.QueryRowx(
		`SELECT criterion_accessibility.id_criterion, `+
			`criterion_accessibility.id_accessibility, accessibility.name, criterion_accessibility.weight, criterion_accessibility.closed `+
			`FROM places4all.criterion_accessibility `+
			`JOIN places4all.accessibility ON accessibility.id = criterion_accessibility.id_accessibility `+
			`WHERE criterion_accessibility.id_criterion = $1 AND `+
			`criterion_accessibility.id_accessibility = $2`, idCriterion, idAccessibility)

	criterionAccessibility := NewCriterionAccessibility(false)
	criterionAccessibility.SetExists()
	err := row.StructScan(criterionAccessibility)
	if err != nil {
		return nil, err
	}

	return criterionAccessibility, err
}

func (ds *Datastore) SaveCriterionAccessibility(criterionAccessibility *CriterionAccessibility) error {

	if criterionAccessibility == nil {
		return errors.New("criterionAccessibility should not be nil")
	}

	if !criterionAccessibility.Exists() {
		return errors.New("update failed: does not exist")
	}

	if criterionAccessibility.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.criterion_accessibility SET weight = $1 WHERE id_criterion = $2 AND id_accessibility = $3`

	_, err := ds.postgres.Exec(sql, criterionAccessibility.Weight, criterionAccessibility.IdCriterion, criterionAccessibility.IdAccessibility)
	return err
}

func (ds *Datastore) InsertCriterionAccessibilityByIds(criterionId, accessibilityId int64, weight int) error {
	const sql = `INSERT INTO places4all.criterion_accessibility(id_criterion, id_accessibility, weight) VALUES($1, $2, $3)`

	_, err := ds.postgres.Exec(sql, criterionId, accessibilityId, weight)
	if err != nil {
		return err
	}

	return err

}
