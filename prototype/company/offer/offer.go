// Manual

package offer

import (
	"fmt"
	"time"
)

//go:generate hatch -file $GOFILE -model Offer -interface Interface -version model4 -keyspace company -table offer -key ID

// Interface is the Offer interface.
type Interface interface {
	// Delete deletes m.
	Delete(m Offer) (Offer, error)
	// New creates m.
	New(m Offer) (Offer, error)
	// Set updates m.
	Set(m Offer) (Offer, error)
}

// Offer is an example model.
type Offer struct {
	// Interface is used by methods.
	Interface Interface
	// Created is the created time.
	Created time.Time
	// Deleted is the deleted time.
	Deleted time.Time
	// ID is the ID.
	ID string
	// Name is the name.
	Name string
	// Updated is the updated time.
	Updated time.Time
}

// BetterName is an example of a non-service method.
func (m Offer) BetterName() string {
	return m.Name + m.Name
}

// Validate validates m.
func (m Offer) Validate() error {
	if m.Created.Before(time.Time{}) {
		return fmt.Errorf("model: invalid created: %v", m.Created)
	}
	if m.Deleted.Before(time.Time{}) || (m.Deleted.After(time.Time{}) && (m.Deleted.Before(m.Created) || (m.Updated.After(time.Time{}) && m.Deleted.Before(m.Updated)))) {
		return fmt.Errorf("model: invalid deleted: %v", m.Deleted)
	}
	if m.ID == "" {
		return fmt.Errorf("model: invalid id: %v", m.ID)
	}
	if m.Name == "" {
		return fmt.Errorf("model: invalid name: %v", m.Name)
	}
	if m.Updated.Before(time.Time{}) || (m.Updated.After(time.Time{}) && (m.Updated.Before(m.Created) || (m.Deleted.After(time.Time{}) && m.Deleted.Before(m.Updated)))) {
		return fmt.Errorf("model: invalid deleted: %v", m.Deleted)
	}
	return nil
}
