#!/bin/sh

# Generated

set -e -x

gore <<EOF
:import github.com/willfaught/company/prototype/company/offer
var offers = offer.LoadAll()
for _, offer := range offers {
    if err := offer.Validate(); err != nil {
        log.Printf("offer %q is invalid: %v", offer.ID, err)
    }
}
EOF
