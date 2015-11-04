// Manual

package offers

//go:generate hatch -file $GOFILE -service Offers -configuration Configuration -context Context

import (
	"fmt"
	"log"
	"time"

	"github.com/kr/pretty"
	tigertonic "github.com/rcrowley/go-tigertonic"
	"github.com/willfaught/company/prototype/company/offer"
)

var _ Interface = Offers{}

// Configuration is the Offers configuration.
type Configuration struct {
	// Name is the Offers name.
	Name string
	// OfferRepository is a model repository.
	OfferRepository offer.Repository
}

// OnCall is called before calling a service method. OnCall is a special method.
func (conf Configuration) OnCall(method string, cont Context, argument interface{}) {
	log.Printf("%v %v: Called\nContext=%# v\nArgument=%# v\n\n", conf.Name, method, pretty.Formatter(cont), pretty.Formatter(argument))
}

// OnPanic is called when recovering from a panic in a service method. OnPanic is a special method.
func (conf Configuration) OnPanic(method string, cont Context, value interface{}) {
	log.Printf("%v %v: Panicked\nContext=%# v\nValue=%# v\n\n", conf.Name, method, pretty.Formatter(cont), pretty.Formatter(value))
}

// OnReturn is called after a service method returns. OnReturn is a special method.
func (conf Configuration) OnReturn(method string, cont Context, result interface{}, err error) {
	log.Printf("%v %v: Returned\nContext=%# v\nResult=%# v\nError=%# v\n\n", conf.Name, method, pretty.Formatter(cont), pretty.Formatter(result), pretty.Formatter(err))
}

// Context is the Offers context.
type Context struct {
	// ID identifies the call.
	ID string
}

// Offers is an example service.
type Offers struct {
	// Configuration is the configuration.
	Configuration Configuration
	// Context is the call context.
	Context Context
}

// Delete deletes o.
func (os Offers) Delete(o offer.Offer) (offer.Offer, error) {
	o, ok := os.Configuration.OfferRepository.Load(o.ID)
	if !ok {
		return offer.Offer{}, fmt.Errorf("invalid id: %q", o.ID)
	}
	if !o.Deleted.IsZero() {
		return offer.Offer{}, fmt.Errorf("invalid id: %q", o.ID)
	}
	o.Deleted = time.Now().UTC()
	os.Configuration.OfferRepository.Store(o)
	return o, nil
}

// ForContext uses c.
func (os Offers) ForContext(c Context) Interface {
	return Offers{Configuration: os.Configuration, Context: c}
}

// Get gets the Offer for id.
func (os Offers) Get(id string) (offer.Offer, error) {
	var o, ok = os.Configuration.OfferRepository.Load(id)
	if !ok {
		return offer.Offer{}, fmt.Errorf("invalid id: %q", id)
	}
	if !o.Deleted.IsZero() {
		return offer.Offer{}, fmt.Errorf("invalid id: %q", id)
	}
	return o, nil
}

// New creates o.
func (os Offers) New(o offer.Offer) (offer.Offer, error) {
	o.ID = tigertonic.RandomBase62String(8)
	o.Created = time.Now().UTC()
	o.Deleted = time.Time{}
	o.Updated = time.Time{}
	if err := o.Validate(); err != nil {
		return offer.Offer{}, err
	}
	os.Configuration.OfferRepository.Store(o)
	return o, nil
}

// Set updates o.
func (os Offers) Set(o offer.Offer) (offer.Offer, error) {
	var old, ok = os.Configuration.OfferRepository.Load(o.ID)
	if !ok {
		return offer.Offer{}, fmt.Errorf("invalid id: %q", o.ID)
	}
	if !old.Deleted.IsZero() {
		return offer.Offer{}, fmt.Errorf("invalid id: %q", o.ID)
	}
	o.Created = old.Created
	o.Deleted = old.Deleted
	o.Updated = time.Now().UTC()
	if err := o.Validate(); err != nil {
		return offer.Offer{}, err
	}
	os.Configuration.OfferRepository.Store(o)
	return o, nil
}
