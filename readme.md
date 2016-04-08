[![Build][Build-Status-Image]][Build-Status-Url]
[![License][License-Image]][License-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]
[![Release][Release-Image]][Release-Url]


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

* https://github.com/go-kit/kit ?
  
  Could be used as an underlying implementation

## Limitations

Only a small subset of Swagger functionality is supported by the tool. 

[License-Url]: https://raw.githubusercontent.com/viktorasm/gontractor/master/LICENSE
[License-Image]: https://img.shields.io/:license-mit-blue.svg
[ReportCard-Url]: http://goreportcard.com/report/viktorasm/gontractor
[ReportCard-Image]: http://goreportcard.com/badge/viktorasm/gontractor
[Build-Status-Url]: http://travis-ci.org/viktorasm/gontractor
[Build-Status-Image]: https://img.shields.io/travis/viktorasm/gontractor.svg
[Release-Url]: https://github.com/viktorasm/gontractor/releases/tag/v0.1
[Release-image]: http://img.shields.io/badge/gontractor-v0.1-1eb0fc.svg
