// Generated

package offer

import (
	"time"

	"github.com/company/gossie/src/gossie"
)

type model5 struct { // Added field: Name
	Created int64 `keyspace:"company" cf:"offer5" key:"ID"`
	Deleted int64
	ID      string
	Name    string
	Updated int64
}

func (m5 *model5) migrateFrom(m Offer) {
	if m.Created.IsZero() {
		m5.Created = 0
	} else {
		m5.Created = m.Created.UnixNano()
	}
	if m.Deleted.IsZero() {
		m5.Deleted = 0
	} else {
		m5.Deleted = m.Deleted.UnixNano()
	}
	m5.ID = m.ID
	m5.Name = m.Name
	if m.Updated.IsZero() {
		m5.Updated = 0
	} else {
		m5.Updated = m.Updated.UnixNano()
	}
}

func (m5 model5) migrateTo(m *Offer) {
	if m5.Created == 0 {
		m.Created = time.Time{}
	} else {
		m.Created = time.Unix(0, m5.Created).UTC()
	}
	if m5.Deleted == 0 {
		m.Deleted = time.Time{}
	} else {
		m.Deleted = time.Unix(0, m5.Deleted).UTC()
	}
	m.ID = m5.ID
	m.Name = m5.Name
	if m5.Updated == 0 {
		m.Updated = time.Time{}
	} else {
		m.Updated = time.Unix(0, m5.Updated).UTC()
	}
}

type repository5 struct {
	connection  gossie.ConnectionPool
	mapping     gossie.Mapping
	repository4 repository4
}

func (r5 repository5) load(id string) (model5, bool) {
	var m4, ok = r5.repository4.load(id)
	if !ok {
		return model5{}, false
	}
	var m5 model5
	if err := r5.connection.Query(r5.mapping).GetOne(id, &m5); err != nil {
		if err != gossie.Done {
			panic(err)
		}
	}
	m4.migrateTo(&m5)
	return m5, true
}

func (r5 repository5) loadAll() []model5 {
	var all []model5
	var id string
	for {
		var some, ok = r5.loadRange(id, 4096)
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

func (r5 repository5) loadRange(id string, limit int) ([]model5, bool) {
	var m4s, ok = r5.repository4.loadRange(id, limit)
	if !ok {
		return nil, false
	}
	var m5s []model5
	for _, m4 := range m4s {
		var m5, ok = r5.load(m4.ID)
		if !ok {
			m4.migrateTo(&m5)
		}
		m5s = append(m5s, m5)
	}
	return m5s, true
}

func (r5 repository5) store(m5 model5) {
	if err := r5.connection.Batch().Insert(r5.mapping, &m5).Run(); err != nil {
		panic(err)
	}
	var m4, _ = r5.repository4.load(m5.ID)
	m4.migrateFrom(m5)
	r5.repository4.store(m4)
}
