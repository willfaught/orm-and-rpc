// Generated

package offer

import "github.com/company/gossie/src/gossie"

type model3 struct { // Renamed field: from Added to Renamed
	Created int64 `keyspace:"company" cf:"offer3" key:"ID"`
	Deleted int64
	ID      string
	Renamed string
	Updated int64
}

func (m3 *model3) migrateFrom(m4 model4) {
	m3.Created = m4.Created
	m3.Deleted = m4.Deleted
	m3.ID = m4.ID
	m3.Updated = m4.Updated
}

func (m3 model3) migrateTo(m4 *model4) {
	m4.Created = m3.Created
	m4.Deleted = m3.Deleted
	m4.ID = m3.ID
	m4.Updated = m3.Updated
}

type repository3 struct {
	connection  gossie.ConnectionPool
	mapping     gossie.Mapping
	repository2 repository2
}

func (r3 repository3) load(id string) (model3, bool) {
	var m2, ok = r3.repository2.load(id)
	if !ok {
		return model3{}, false
	}
	var m3 model3
	if err := r3.connection.Query(r3.mapping).GetOne(id, &m3); err != nil {
		if err != gossie.Done {
			panic(err)
		}
	}
	m2.migrateTo(&m3)
	return m3, true
}

func (r3 repository3) loadAll() []model3 {
	var all []model3
	var id string
	for {
		var some, ok = r3.loadRange(id, 4096)
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

func (r3 repository3) loadRange(id string, limit int) ([]model3, bool) {
	var m2s, ok = r3.repository2.loadRange(id, limit)
	if !ok {
		return nil, false
	}
	var m3s []model3
	for _, m2 := range m2s {
		var m3, ok = r3.load(m2.ID)
		if !ok {
			m2.migrateTo(&m3)
		}
		m3s = append(m3s, m3)
	}
	return m3s, true
}

func (r3 repository3) store(m3 model3) {
	if err := r3.connection.Batch().Insert(r3.mapping, &m3).Run(); err != nil {
		panic(err)
	}
	var m2, _ = r3.repository2.load(m3.ID)
	m2.migrateFrom(m3)
	r3.repository2.store(m2)
}
