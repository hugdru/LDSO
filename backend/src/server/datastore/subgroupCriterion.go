package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type SubgroupCriterion struct {
	IdSubgroup  int64 `json:"idSubgroup" db:"id_subgroup"`
	IdCriterion int64 `json:"idCriterion" db:"id_criterion"`
	meta        metadata.Metadata
}

func (sc *SubgroupCriterion) SetExists() {
	sc.meta.Exists = true
}

func (sc *SubgroupCriterion) SetDeleted() {
	sc.meta.Deleted = true
}

func (sc *SubgroupCriterion) Exists() bool {
	return sc.meta.Exists
}

func (sc *SubgroupCriterion) Deleted() bool {
	return sc.meta.Deleted
}

func (ds *Datastore) InsertSubgroupCriterion(sc *SubgroupCriterion) error {
	var err error

	if sc.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.subgroup_criterion (` +
		`id_subgroup` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id_criterion`

	err = ds.postgres.QueryRow(sql, sc.IdSubgroup).Scan(&sc.IdCriterion)
	if err != nil {
		return err
	}

	sc.SetExists()

	return nil
}

func (ds *Datastore) UpdateSubgroupCriterion(sc *SubgroupCriterion) error {
	var err error

	if !sc.Exists() {
		return errors.New("update failed: does not exist")
	}

	if sc.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.subgroup_criterion SET (` +
		`id_subgroup` +
		`) = ( ` +
		`$1` +
		`) WHERE id_criterion = $2`

	_, err = ds.postgres.Exec(sql, sc.IdSubgroup, sc.IdCriterion)
	return err
}

func (ds *Datastore) SaveSubgroupCriterion(sc *SubgroupCriterion) error {
	if sc.Exists() {
		return ds.UpdateSubgroupCriterion(sc)
	}

	return ds.InsertSubgroupCriterion(sc)
}

func (ds *Datastore) UpsertSubgroupCriterion(sc *SubgroupCriterion) error {
	var err error

	if sc.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.subgroup_criterion (` +
		`id_subgroup, id_criterion` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id_criterion) DO UPDATE SET (` +
		`id_subgroup, id_criterion` +
		`) = (` +
		`EXCLUDED.id_subgroup, EXCLUDED.id_criterion` +
		`)`

	_, err = ds.postgres.Exec(sql, sc.IdSubgroup, sc.IdCriterion)
	if err != nil {
		return err
	}

	sc.SetExists()

	return nil
}

func (ds *Datastore) DeleteSubgroupCriterion(sc *SubgroupCriterion) error {
	var err error

	if !sc.Exists() {
		return nil
	}

	if sc.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.subgroup_criterion WHERE id_criterion = $1`

	_, err = ds.postgres.Exec(sql, sc.IdCriterion)
	if err != nil {
		return err
	}

	sc.SetDeleted()

	return nil
}

func (ds *Datastore) GetSubgroupCriterionCriterion(sc *SubgroupCriterion) (*Criterion, error) {
	return ds.GetCriterionById(sc.IdCriterion)
}

func (ds *Datastore) GetSubgroupCriterionSubgroup(sc *SubgroupCriterion) (*Subgroup, error) {
	return ds.GetSubgroupById(sc.IdSubgroup)
}

func (ds *Datastore) GetSubgroupCriterionByIds(idSubgroup, idCriterion int64) (*SubgroupCriterion, error) {
	var err error

	const sql = `SELECT ` +
		`id_subgroup, id_criterion ` +
		`FROM places4all.subgroup_criterion ` +
		`WHERE id_subgroup = $1 AND id_criterion = $2`

	sc := SubgroupCriterion{}
	sc.SetExists()

	err = ds.postgres.QueryRow(sql, idSubgroup, idCriterion).Scan(&sc.IdSubgroup, &sc.IdCriterion)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}
