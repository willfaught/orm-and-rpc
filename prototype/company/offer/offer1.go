// Generated

package offer

import "github.com/company/gossie/src/gossie"

type model1 struct {
	Created int64 `keyspace:"company" cf:"offer" key:"ID"`
	Deleted int64
	ID      string
	Updated int64
}

func (m1 *model1) migrateFrom(m2 model2) {
	m1.Created = m2.Created
	m1.Deleted = m2.Deleted
	m1.ID = m2.ID
	m1.Updated = m2.Updated
}

func (m1 model1) migrateTo(m2 *model2) {
	m2.Created = m1.Created
	m2.Deleted = m1.Deleted
	m2.ID = m1.ID
	m2.Updated = m1.Updated
}

type repository1 struct {
	connection gossie.ConnectionPool
	mapping    gossie.Mapping
}

func (r1 repository1) load(id string) (model1, bool) {
	var m1 model1
	if err := r1.connection.Query(r1.mapping).GetOne(id, &m1); err != nil {
		if err != gossie.Done {
			panic(err)
		}
		return model1{}, false
	}
	return m1, true
}

func (r1 repository1) loadAll() []model1 {
	var all []model1
	var id string
	for {
		var some, ok = r1.loadRange(id, 4096)
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

func (r1 repository1) loadRange(id string, limit int) ([]model1, bool) {
	var therange = &gossie.Range{Count: limit}
	if id != "" {
		therange.Start = []byte(id)
		therange.Count += 1
	}
	var result, err = r1.connection.Query(r1.mapping).RangeGet(therange)
	if err != nil {
		panic(err)
	}
	if id != "" {
		if err := result.Next(&model1{}); err != nil {
			panic(err)
		}
	}
	var some []model1
	for len(some) < limit {
		var m1 model1
		if err := result.Next(&m1); err != nil {
			if err != gossie.Done {
				panic(err)
			}
			break
		}
		some = append(some, m1)
	}
	return some, true
}

func (r1 repository1) store(m1 model1) {
	if err := r1.connection.Batch().Insert(r1.mapping, &m1).Run(); err != nil {
		panic(err)
	}
}
