package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/company/gossie/src/mockgossie"
	"github.com/willfaught/company/prototype/company/offer"
	"github.com/willfaught/company/prototype/company/offers"
)

func TestLocal(t *testing.T) {
	suite.Run(t, &offers.Suite{Interface: func() offers.Interface {
		return offers.Offers{Configuration: offers.Configuration{OfferRepository: offer.NewRepository(mockgossie.NewMockConnectionPool())}}
	}})
}

func TestRemote(t *testing.T) {
	// Service
	var service = offers.Offers{Configuration: offers.Configuration{OfferRepository: offer.NewRepository(mockgossie.NewMockConnectionPool())}}

	// Server
	var server = offers.MustNewServer(offers.Configuration{}, service, ":5000")
	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	// Client
	var client = offers.MustNewClient(offers.Configuration{}, ":5000")

	suite.Run(t, &offers.Suite{Interface: func() offers.Interface {
		return client
	}})

	// Client
	if err := client.Close(); err != nil {
		panic(err)
	}

	// Server
	if err := server.Stop(); err != nil {
		panic(err)
	}
}
