# gontractor (proof-of-concept)

Enables contract-first creation of Golang services by generating boilerplate code.

## Why another microservice lib?
* Write service contract as Swagger specification
* Generate boilerplate (request/response types, service interface, server setup/startup code);
* Add service implementation manually
* Update contract
* Regenerate boilerplate
* Update service implementation as needed;

## Why not...

* https://github.com/go-swagger/go-swagger ?
  Output too bloated. 
* https://github.com/go-kit/kit
  Could be used as an underlying implementation

## Current status


## Limitations

Only a small subset of Swagger functionality is supported by the tool. 