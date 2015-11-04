// Generated

package offers

import (
	"net/rpc"

	"github.com/willfaught/company/prototype/company/offer"
)

var _ Interface = Client{}

// Client is the RPC client.
type Client struct {
	// Client is the wrapped RPC client.
	*rpc.Client
	// Configuration is the configuration.
	Configuration Configuration
	// Context is the call context.
	Context Context
	// Name is the receiver name.
	Name string
}

// MustNewClient calls NewClient and panics if there is an error.
func MustNewClient(c Configuration, address string) Client {
	var client, err = NewClient(c, address)
	if err != nil {
		panic(err)
	}
	return client
}

// NewClient makes a Client that uses h and dials address.
func NewClient(c Configuration, address string) (Client, error) {
	var client, err = rpc.Dial("tcp", address)
	if err != nil {
		return Client{}, err
	}
	return Client{Client: client, Configuration: c, Name: "Marketing"}, nil // Default receiver name is "Marketing".
}

// Delete deletes m.
func (c Client) Delete(m offer.Offer) (offer.Offer, error) {
	m.Interface = nil // Avoid encoding the Service or Client
	var argument = struct {
		Arguments struct { // Method arguments
			M offer.Offer // m
		}
		Context Context // Call context
	}{
		Arguments: struct { // Method arguments
			M offer.Offer // m
		}{
			M: m, // m
		},
		Context: c.Context, // Call context
	}
	var result struct { // Method result
		Offer offer.Offer // Only one non-error result
	}
	c.Configuration.OnCall("Delete", argument.Context, argument.Arguments) // Called because Configuration implements it
	var err = c.Client.Call(c.Name+".Delete", argument, &result)
	c.Configuration.OnReturn("Delete", argument.Context, result, err) // Called because Configuration implements it
	result.Offer.Interface = c                                        // Enable method calls like m.Delete(). See offer.Offer.Delete.
	return result.Offer, err
}

// ForContext uses c.
func (c Client) ForContext(co Context) Interface { // Generated because ForContext(Context) Interface was found in the Service interface.
	return Client{Client: c.Client, Configuration: c.Configuration, Context: co, Name: c.Name}
}

// Get gets the Offer for id.
func (c Client) Get(id string) (offer.Offer, error) {
	var argument = struct {
		Arguments struct {
			Id string // id
		}
		Context Context
	}{
		Arguments: struct {
			Id string // id
		}{
			Id: id, // id
		},
		Context: c.Context,
	}
	var result struct {
		Offer offer.Offer
	}
	c.Configuration.OnCall("Get", argument.Context, argument.Arguments)
	var err = c.Client.Call(c.Name+".Get", argument, &result)
	c.Configuration.OnReturn("Get", argument.Context, result, err)
	result.Offer.Interface = c
	return result.Offer, err
}

// New creates m.
func (c Client) New(m offer.Offer) (offer.Offer, error) {
	m.Interface = nil
	var argument = struct {
		Arguments struct {
			M offer.Offer
		}
		Context Context
	}{
		Arguments: struct {
			M offer.Offer
		}{
			M: m,
		},
		Context: c.Context,
	}
	var result struct {
		Offer offer.Offer
	}
	c.Configuration.OnCall("New", argument.Context, argument.Arguments)
	var err = c.Client.Call(c.Name+".New", argument, &result)
	c.Configuration.OnReturn("New", argument.Context, result, err)
	result.Offer.Interface = c
	return result.Offer, err
}

// Set updates m.
func (c Client) Set(m offer.Offer) (offer.Offer, error) {
	m.Interface = nil
	var argument = struct {
		Arguments struct {
			M offer.Offer
		}
		Context Context
	}{
		Arguments: struct {
			M offer.Offer
		}{
			M: m,
		},
		Context: c.Context,
	}
	var result struct {
		Offer offer.Offer
	}
	c.Configuration.OnCall("Set", argument.Context, argument.Arguments)
	var err = c.Client.Call(c.Name+".Set", argument, &result)
	c.Configuration.OnReturn("Set", argument.Context, result, err)
	result.Offer.Interface = c
	return result.Offer, err
}
