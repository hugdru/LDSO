package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type CriterionAccessibility struct {
	IdCriterion     int64 `json:"idCriterion" db:"id_criterion"`
	IdAccessibility int64 `json:"idAccessibility" db:"id_accessibility"`
	Weight          int   `json:"weight" db:"weight"`
	meta            metadata.Metadata
}

func (ca *CriterionAccessibility) SetExists() {
	ca.meta.Exists = true
}

func (ca *CriterionAccessibility) SetDeleted() {
	ca.meta.Deleted = true
}

func (ca *CriterionAccessibility) Exists() bool {
	return ca.meta.Exists
}

func (ca *CriterionAccessibility) Deleted() bool {
	return ca.meta.Deleted
}

func (ds *Datastore) InsertCriterionAccessibility(ca *CriterionAccessibility) error {
	var err error

	if ca.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.criterion_accessibility (` +
		`id_criterion, weight` +
		`) VALUES (` +
		`$1, $2` +
		`) RETURNING id_accessibility`

	err = ds.postgres.QueryRow(sql, ca.IdCriterion, ca.Weight).Scan(&ca.IdAccessibility)
	if err != nil {
		return err
	}

	ca.SetExists()

	return nil
}

func (ds *Datastore) UpdateCriterionAccessibility(ca *CriterionAccessibility) error {
	var err error

	if !ca.Exists() {
		return errors.New("update failed: does not exist")
	}

	if ca.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.criterion_accessibility SET (` +
		`id_criterion, weight` +
		`) = ( ` +
		`$1, $2` +
		`) WHERE id_accessibility = $3`

	_, err = ds.postgres.Exec(sql, ca.IdCriterion, ca.Weight, ca.IdAccessibility)
	return err
}

func (ds *Datastore) SaveCriterionAccessibility(ca *CriterionAccessibility) error {
	if ca.Exists() {
		return ds.UpdateCriterionAccessibility(ca)
	}

	return ds.InsertCriterionAccessibility(ca)
}

func (ds *Datastore) UpsertCriterionAccessibility(ca *CriterionAccessibility) error {
	var err error

	if ca.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.criterion_accessibility (` +
		`id_criterion, id_accessibility, weight` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) ON CONFLICT (id_accessibility) DO UPDATE SET (` +
		`id_criterion, id_accessibility, weight` +
		`) = (` +
		`EXCLUDED.id_criterion, EXCLUDED.id_accessibility, EXCLUDED.weight` +
		`)`

	_, err = ds.postgres.Exec(sql, ca.IdCriterion, ca.IdAccessibility, ca.Weight)
	if err != nil {
		return err
	}

	ca.SetExists()

	return nil
}

func (ds *Datastore) DeleteCriterionAccessibility(ca *CriterionAccessibility) error {
	var err error

	if !ca.Exists() {
		return nil
	}

	if ca.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.criterion_accessibility WHERE id_accessibility = $1`

	_, err = ds.postgres.Exec(sql, ca.IdAccessibility)
	if err != nil {
		return err
	}

	ca.SetDeleted()

	return nil
}

func (ds *Datastore) GetCriterionAccessibilityAccessibility(ca *CriterionAccessibility) (*Accessibility, error) {
	return ds.GetAccessibilityById(ca.IdAccessibility)
}

func (ds *Datastore) GetCriterionAccessibilityCriterion(ca *CriterionAccessibility) (*Criterion, error) {
	return ds.GetCriterionById(ca.IdCriterion)
}

func (ds *Datastore) GetCriterionAccessibilityByIds(idCriterion, idAccessibility int64) (*CriterionAccessibility, error) {
	var err error

	const sql = `SELECT ` +
		`id_criterion, id_accessibility, weight ` +
		`FROM places4all.criterion_accessibility ` +
		`WHERE id_criterion = $1 AND id_accessibility = $2`

	ca := CriterionAccessibility{}
	ca.SetExists()

	err = ds.postgres.QueryRow(sql, idCriterion, idAccessibility).Scan(&ca.IdCriterion, &ca.IdAccessibility, &ca.Weight)
	if err != nil {
		return nil, err
	}

	return &ca, nil
}
