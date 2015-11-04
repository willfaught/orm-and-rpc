# Better services

## Goals

- Adaptive
- Faster
- Generateable
- Safer
- Simpler

## Features

- General
    - Smaller APIs from simpler error handling: panics for unrecoverable errors, including unrecoverable network errors
    - Value types used wherever possible to avoid GC load
    - Services and models are independent, can mix and match
- Services
    - No complicated and vague HTTP semantics: direct RPC for all service methods
    - Service interface versions and migration
    - Faster RPC
        - Gob encoding is more efficient than JSON encoding
        - [JSON is not a good fit for internal communication](https://www.hakkalabs.co/articles/distributed-systems-go-good-bad-otherwise)
    - Standalone services that can be mixed arbitrarily with any RPC server
    - Support multiple, simultaneous encodings, like gob and HTTP JSON-RPC
    - Service interfaces enable:
         - Two services with high-bandwidth communication to be co-located in the same process as an optimization, only violating deployment isolation as a tradeoff for better performance
         - Easy proxying
         - Tests that double as unit and integration tests
- Models
    - Separate packages for separate namespaces
    - [Clean architecture](https://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html): models use repositories that know nothing about implementation details
    - Models are independent of services with methods that plug into compatible services
    - Model versions and schema migration
    - Rich Go types for model fields
    - Zero configuration: default column and JSON field names

## To do

- Fix the JSON-RPC example

## Notes

The shell scripts have not been tested.
