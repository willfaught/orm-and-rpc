// Generated

package offer

import "github.com/company/gossie/src/gossie"

var _ Repository = repository{}

type Repository interface { // Implementation unknown to all models
	Load(id string) (Offer, bool)
	LoadAll() []Offer
	LoadRange(id string, limit int) ([]Offer, bool)
	Store(m Offer)
}

func NewRepository(c gossie.ConnectionPool) Repository {
	var r1 = repository1{connection: c, mapping: gossie.MustNewMapping(&model1{})}                  // v1
	var r2 = repository2{connection: c, mapping: gossie.MustNewMapping(&model2{}), repository1: r1} // v2, chained to v1
	var r3 = repository3{connection: c, mapping: gossie.MustNewMapping(&model3{}), repository2: r2} // v3, chained to v2
	var r4 = repository4{connection: c, mapping: gossie.MustNewMapping(&model4{}), repository3: r3} // v4, chained to v3
	var r5 = repository5{connection: c, mapping: gossie.MustNewMapping(&model5{}), repository4: r4} // v5, chained to v4
	return repository{repository5: r5}                                                              // The current version, chained to v5
}

type repository struct { // The current version, chained to v5, unknown to all models
	repository5 repository5
}

func (r repository) Load(id string) (Offer, bool) {
	var m5, ok = r.repository5.load(id) // Load the v5 version of id, recurse down to v1
	if !ok {
		return Offer{}, false
	}
	var m Offer
	m5.migrateTo(&m) // Migrate from v5 to the current version
	return m, true
}

func (r repository) LoadAll() []Offer {
	var all []Offer
	var id string
	for {
		var some, ok = r.LoadRange(id, 4096)
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

func (r repository) LoadRange(id string, limit int) ([]Offer, bool) {
	var m5s, ok = r.repository5.loadRange(id, limit) // Load the v5 version of the range, recurse down to v1
	if !ok {
		return nil, false
	}
	var ms []Offer
	for _, m5 := range m5s {
		var m Offer
		m5.migrateTo(&m) // Migrate from v5 to the current version
		ms = append(ms, m)
	}
	return ms, true
}

func (r repository) Store(m Offer) {
	var m5, _ = r.repository5.load(m.ID)
	m5.migrateFrom(m)       // Migrate from the current version to v5
	r.repository5.store(m5) // Recurse down to v1
}
