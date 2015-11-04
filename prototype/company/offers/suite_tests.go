// Manual

package offers

import (
	"time"

	"github.com/willfaught/company/prototype/company/offer"
)

func (s *Suite) TestDelete() {
	var i = s.Interface()

	var m, err = i.Delete(offer.Offer{})
	s.Error(err)

	m, err = i.Delete(offer.Offer{ID: "nerddomo"})
	s.Error(err)

	m, err = i.New(offer.Offer{Name: "nerddomo"})
	s.NoError(err)
	n, err := i.Delete(m)
	s.NoError(err)
	s.False(n.Created.IsZero())
	s.NotEmpty(n.ID)
	s.Equal("nerddomo", n.Name)
	s.True(n.Updated.IsZero())
	n.Deleted = time.Time{}
	s.Equal(m, n)

	m, err = i.New(offer.Offer{Name: "nerddomo"})
	s.NoError(err)
	n, err = i.Delete(m)
	s.NoError(err)
	n, err = i.Delete(m)
	s.Error(err)
}

func (s *Suite) TestGet() {
	var i = s.Interface()

	var m, err = i.Get("")
	s.Error(err)

	m, err = i.Get("nerddomo")
	s.Error(err)

	m, err = i.New(offer.Offer{Name: "nerddomo"})
	s.NoError(err)
	n, err := i.Get(m.ID)
	s.NoError(err)
	s.Equal(m, n)

	m, err = i.New(offer.Offer{Name: "nerddomo"})
	s.NoError(err)
	n, err = i.Delete(m)
	s.NoError(err)
	n, err = i.Get(n.ID)
	s.Error(err)
}

func (s *Suite) TestNew() {
	var i = s.Interface()

	var m, err = i.New(offer.Offer{})
	s.Error(err)

	m, err = i.New(offer.Offer{ID: "nerddomo"})
	s.Error(err)

	var now = time.Now()
	m, err = i.New(offer.Offer{Created: now, Deleted: now, ID: "nerddomo", Name: "nerddomo", Updated: now})
	s.NoError(err)
	s.False(m.Created.IsZero())
	s.NotEqual(now, m.Created)
	s.True(m.Deleted.IsZero())
	s.NotEmpty(m.ID)
	s.NotEqual("nerddomo", m.ID)
	s.Equal("nerddomo", m.Name)
	s.True(m.Updated.IsZero())
}

func (s *Suite) TestSet() {
	var i = s.Interface()

	var m, err = i.Set(offer.Offer{})
	s.Error(err)

	m, err = i.Set(offer.Offer{ID: "nerddomo"})
	s.Error(err)

	m, err = i.New(offer.Offer{Name: "nerddomo"})
	s.NoError(err)
	m.Name = ""
	m, err = i.Set(m)
	s.Error(err)

	m, err = i.New(offer.Offer{Name: "nerddomo"})
	s.NoError(err)
	var n = m
	var now = time.Now()
	n.Created = now
	n.Deleted = time.Now()
	n.Updated = now
	n, err = i.Set(n)
	s.NoError(err)
	s.Equal(m.Created, n.Created)
	s.Equal(m.Deleted, n.Deleted)
	s.Equal(m.ID, n.ID)
	s.Equal(m.Name, n.Name)
	s.False(n.Updated.IsZero())
	s.NotEqual(now, n.Updated)

	m, err = i.New(offer.Offer{Name: "nerddomo"})
	s.NoError(err)
	n, err = i.Delete(m)
	s.NoError(err)
	n, err = i.Set(n)
	s.Error(err)
}
