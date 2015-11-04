// Generated

package offer

import "github.com/company/gossie/src/gossie"

type model4 struct { // Removed field: Removed
	Created int64 `keyspace:"company" cf:"offer4" key:"ID"`
	Deleted int64
	ID      string
	Updated int64
}

func (m4 *model4) migrateFrom(m5 model5) {
	m4.Created = m5.Created
	m4.Deleted = m5.Deleted
	m4.ID = m5.ID
	m4.Updated = m5.Updated
}

func (m4 model4) migrateTo(m5 *model5) {
	m5.Created = m4.Created
	m5.Deleted = m4.Deleted
	m5.ID = m4.ID
	m5.Updated = m4.Updated
}

type repository4 struct {
	connection  gossie.ConnectionPool
	mapping     gossie.Mapping
	repository3 repository3
}

func (r4 repository4) load(id string) (model4, bool) {
	var m3, ok = r4.repository3.load(id)
	if !ok {
		return model4{}, false
	}
	var m4 model4
	if err := r4.connection.Query(r4.mapping).GetOne(id, &m4); err != nil {
		if err != gossie.Done {
			panic(err)
		}
	}
	m3.migrateTo(&m4)
	return m4, true
}

func (r4 repository4) loadAll() []model4 {
	var all []model4
	var id string
	for {
		var some, ok = r4.loadRange(id, 4096)
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

func (r4 repository4) loadRange(id string, limit int) ([]model4, bool) {
	var m3s, ok = r4.repository3.loadRange(id, limit)
	if !ok {
		return nil, false
	}
	var m4s []model4
	for _, m3 := range m3s {
		var m4, ok = r4.load(m3.ID)
		if !ok {
			m3.migrateTo(&m4)
		}
		m4s = append(m4s, m4)
	}
	return m4s, true
}

func (r4 repository4) store(m4 model4) {
	if err := r4.connection.Batch().Insert(r4.mapping, &m4).Run(); err != nil {
		panic(err)
	}
	var m3, _ = r4.repository3.load(m4.ID)
	m3.migrateFrom(m4)
	r4.repository3.store(m3)
}
