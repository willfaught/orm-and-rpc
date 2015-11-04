// Generated

package offer

import "github.com/company/gossie/src/gossie"

type model2 struct { // Added field: Added
	Added   string `keyspace:"company" cf:"offer2" key:"ID"`
	Created int64
	Deleted int64
	ID      string
	Updated int64
}

func (m2 *model2) migrateFrom(m3 model3) {
	m2.Added = m3.Renamed
	m2.Created = m3.Created
	m2.Deleted = m3.Deleted
	m2.ID = m3.ID
	m2.Updated = m3.Updated
}

func (m2 model2) migrateTo(m3 *model3) {
	m3.Created = m2.Created
	m3.Deleted = m2.Deleted
	m3.ID = m2.ID
	m3.Renamed = m2.Added
	m3.Updated = m2.Updated
}

type repository2 struct {
	connection  gossie.ConnectionPool
	mapping     gossie.Mapping
	repository1 repository1
}

func (r2 repository2) load(id string) (model2, bool) {
	var m1, ok = r2.repository1.load(id)
	if !ok {
		return model2{}, false
	}
	var m2 model2
	if err := r2.connection.Query(r2.mapping).GetOne(id, &m2); err != nil {
		if err != gossie.Done {
			panic(err)
		}
	}
	m1.migrateTo(&m2)
	return m2, true
}

func (r2 repository2) loadAll() []model2 {
	var all []model2
	var id string
	for {
		var some, ok = r2.loadRange(id, 4096)
		if !ok {
			panic("invalid model id: " + id)
		}
		if len(some) == 0 {
			break
		}
		all = append(all, some...)
		id = some[len(some)-1].ID
	}
	return all
}

func (r2 repository2) loadRange(id string, limit int) ([]model2, bool) {
	var m1s, ok = r2.repository1.loadRange(id, limit)
	if !ok {
		return nil, false
	}
	var m2s []model2
	for _, m1 := range m1s {
		var m2, ok = r2.load(m1.ID)
		if !ok {
			m1.migrateTo(&m2)
		}
		m2s = append(m2s, m2)
	}
	return m2s, true
}

func (r2 repository2) store(m2 model2) {
	if err := r2.connection.Batch().Insert(r2.mapping, &m2).Run(); err != nil {
		panic(err)
	}
	var m1, _ = r2.repository1.load(m2.ID)
	m1.migrateFrom(m2)
	r2.repository1.store(m1)
}
