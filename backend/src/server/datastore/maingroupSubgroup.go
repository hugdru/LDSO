package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type MaingroupSubgroup struct {
	IdMaingroup int64 `json:"idMaingroup" db:"id_maingroup"`
	IdSubgroup  int64 `json:"idSubgroup" db:"id_subgroup"`
	meta        metadata.Metadata
}

func (ms *MaingroupSubgroup) SetExists() {
	ms.meta.Exists = true
}

func (ms *MaingroupSubgroup) SetDeleted() {
	ms.meta.Deleted = true
}

func (ms *MaingroupSubgroup) Exists() bool {
	return ms.meta.Exists
}

func (ms *MaingroupSubgroup) Deleted() bool {
	return ms.meta.Deleted
}

func (ds *Datastore) InsertMaingroupSubgroup(ms *MaingroupSubgroup) error {
	var err error

	if ms.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.maingroup_subgroup (` +
		`id_maingroup` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id_subgroup`

	err = ds.postgres.QueryRow(sql, ms.IdMaingroup).Scan(&ms.IdSubgroup)
	if err != nil {
		return err
	}

	ms.SetExists()

	return nil
}

func (ds *Datastore) UpdateMaingroupSubgroup(ms *MaingroupSubgroup) error {
	var err error

	if !ms.Exists() {
		return errors.New("update failed: does not exist")
	}

	if ms.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.maingroup_subgroup SET (` +
		`id_maingroup` +
		`) = ( ` +
		`$1` +
		`) WHERE id_subgroup = $2`

	_, err = ds.postgres.Exec(sql, ms.IdMaingroup, ms.IdSubgroup)
	return err
}

func (ds *Datastore) SaveMaingroupSubgroup(ms *MaingroupSubgroup) error {
	if ms.Exists() {
		return ds.UpdateMaingroupSubgroup(ms)
	}

	return ds.InsertMaingroupSubgroup(ms)
}

func (ds *Datastore) UpsertMaingroupSubgroup(ms *MaingroupSubgroup) error {
	var err error

	if ms.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.maingroup_subgroup (` +
		`id_maingroup, id_subgroup` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id_subgroup) DO UPDATE SET (` +
		`id_maingroup, id_subgroup` +
		`) = (` +
		`EXCLUDED.id_maingroup, EXCLUDED.id_subgroup` +
		`)`

	_, err = ds.postgres.Exec(sql, ms.IdMaingroup, ms.IdSubgroup)
	if err != nil {
		return err
	}

	ms.SetExists()

	return nil
}

func (ds *Datastore) DeleteMaingroupSubgroup(ms *MaingroupSubgroup) error {
	var err error

	if !ms.Exists() {
		return nil
	}

	if ms.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.maingroup_subgroup WHERE id_subgroup = $1`

	_, err = ds.postgres.Exec(sql, ms.IdSubgroup)
	if err != nil {
		return err
	}

	ms.SetDeleted()

	return nil
}

func (ds *Datastore) GetMaingroupSubgroupMaingroup(ms *MaingroupSubgroup) (*Maingroup, error) {
	return ds.GetMaingroupById(ms.IdMaingroup)
}

func (ds *Datastore) GetMaingroupSubgroupSubgroup(ms *MaingroupSubgroup) (*Subgroup, error) {
	return ds.GetSubgroupById(ms.IdSubgroup)
}

func (ds *Datastore) GetMaingroupSubgroupByIds(idMaingroup, idSubgroup int64) (*MaingroupSubgroup, error) {
	var err error

	const sql = `SELECT ` +
		`id_maingroup, id_subgroup ` +
		`FROM places4all.maingroup_subgroup ` +
		`WHERE id_maingroup = $1 AND id_subgroup = $2`

	ms := MaingroupSubgroup{}
	ms.SetExists()

	err = ds.postgres.QueryRow(sql, idMaingroup, idSubgroup).Scan(&ms.IdMaingroup, &ms.IdSubgroup)
	if err != nil {
		return nil, err
	}

	return &ms, nil
}
