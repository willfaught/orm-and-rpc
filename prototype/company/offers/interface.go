// Generated

package offers

import "github.com/willfaught/company/prototype/company/offer"

// Interface is the service interface.
type Interface interface {
	// Delete deletes m.
	Delete(m offer.Offer) (offer.Offer, error)
	// ForContext uses c.
	ForContext(c Context) Interface
	// Get gets the Offer for id.
	Get(id string) (offer.Offer, error)
	// New creates m.
	New(m offer.Offer) (offer.Offer, error)
	// Set updates m.
	Set(m offer.Offer) (offer.Offer, error)
}
