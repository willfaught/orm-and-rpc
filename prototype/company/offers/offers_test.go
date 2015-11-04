// Manual

package offers

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/company/gossie/src/mockgossie"
	"github.com/willfaught/company/prototype/company/offer"
)

func Test(t *testing.T) {
	suite.Run(t, &Suite{Interface: func() Interface {
		return Service{Configuration: Configuration{OfferRepository: offer.NewRepository(mockgossie.NewMockConnectionPool())}}
	}})
}
