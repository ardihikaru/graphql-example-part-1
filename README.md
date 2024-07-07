<p align="center">
<a href="https://saweria.co/ardihikaru">
<img src="https://mfardiansyah.id/assets/images/saweria.png" width="30%" />
</a>
<br>
Love my work? Drop me a coffee here. :)
</p>

<p align="center">
<a href="#">
<img src="https://img.shields.io/badge/%20Platforms-Windows%20/%20Linux-blue.svg?style=flat-square" alt="Platforms" />
</a>
<img src="https://img.shields.io/badge/%20Licence-MIT-green.svg?style=flat-square" alt="license" />
</p>
<p align="center">
<a href="https://github.com/ardihikaru/graphql-example-part-1/blob/master/CODE_OF_CONDUCT.md">
<img src="https://img.shields.io/badge/Community-Code%20of%20Conduct-orange.svg?style=flat-squre" alt="Code of Conduct" />
</a>
<a href="https://github.com/ardihikaru/graphql-example-part-1/blob/master/SUPPORT.md">
<img src="https://img.shields.io/badge/Community-Support-red.svg?style=flat-square" alt="Support" />
</a>
<a href="https://github.com/ardihikaru/graphql-example-part-1/blob/master/CONTRIBUTING.md">
<img src="https://img.shields.io/badge/%20Community-Contribution-yellow.svg?style=flat-square" alt="Contribution" />
</a>
</p>
<hr>

# GraphQL-based Service with Go-Chi

Global Template Repository for Development and Operations Of Your Projects.

| Key               | Values                                                                                  |
|-------------------|-----------------------------------------------------------------------------------------|
| Author            | Muhammad Febrian Ardiansyah                                                             |
| Email             | mfardiansyah.id@gmail.com                                                               |
| LinkedIn          | [Muhammad Febrian Ardiansyah](https://www.linkedin.com/in/muhammad-febrian-ardiansyah/) |
| Personal Homepage | [https://mfardiansyah.id](https://mfardiansyah.id/)                                     |

## Table of Contents

* [Dependencies](#dependencies)
* [Prerequisites](#prerequisites)
* [Installation](#installation)
* [Development](#development)
* [Usage](#usage)
* [Generate private and public key](#generate-private-and-public-key)
* [Using graphQL generator](#using-graphql-generator)
* [Contributing](#contributing)
* [License](#license)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

```shell
apt-get -y install git
```

Or

```shell
yum -y install git
```

### Installation

A step by step series of examples that tell you how to get a development env running

Say what the step will be clone this repository.

```shell
git clone git@github.com:ardihikaru/graphql-example-part-1.git
```

* Golang linter
* `golang-ci` is one of the **IMPORTANT** packages. Any developer who will maintain this project should install it. The
  installation command is as follows:
  ```shell
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
  ```
* Everytime the code is updated, please run following command:
  ```shell
  golangci-lint run ./...
  
  # or this command:
  golangci-lint run ./... --fast -v
  ```

## Development

- N/A

## Usage

Reference and programming instructional materials.

## Generate private and public key

* [referenced article](https://www.scottbrady91.com/openssl/creating-rsa-keys-using-openssl)
  ```shell
  # generate a private key with the correct length
  openssl genrsa -out private-key.pem 3072

  # generate corresponding public key
  openssl rsa -in private-key.pem -pubout -out public-key.pem

  # optional: create a self-signed certificate
  openssl req -new -x509 -key private-key.pem -out cert.pem -days 360

  # optional: convert pem to pfx
  openssl pkcs12 -export -inkey private-key.pem -in cert.pem -out cert.pfx
  ```

## Using GraphQL Generator

* Cli command to use (one by one)
```shell
go get github.com/99designs/gqlgen@v0.17.45

go run github.com/99designs/gqlgen generate
```

* Cli command to use (at once)
```shell
go get github.com/99designs/gqlgen@v0.17.45 && go run github.com/99designs/gqlgen generate
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

Looking to contribute to our code but need some help? There's a few ways to get information:

* Connect with me on [Twitter](https://twitter.com/ardikucing)
* Connect with me on [Facebook](https://facebook.com/ardihikaru)
* Connect with me on [LinkedIn](https://linkedin.com/in/muhammad-febrian-ardiansyah)
* Log an issue here on github

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/ardihikaru/graphql-example-part-1/tags).

## Authors

* **[Muhammad Febrian Ardiansyah](https://github.com/ardihikaru)** - *Initial work*

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

<p> Copyright &copy; 2024 Public Use. All Rights Reserved.

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc

## MISC

* Validates CORS
```shell
curl -v --request OPTIONS 'http://localhost:8080/public/service-id' -H 'Origin: http://other-domain.com' -H 'Access-Control-Request-Method: GET'
```
```shell
curl -v --request OPTIONS 'http://localhost:8080/auth/login' -H 'Origin: http://other-domain.com' -H 'Access-Control-Request-Method: POST'
```
```shell
curl -v -X OPTIONS \
  http://localhost:8080/public/service-id \
  -H 'cache-control: no-cache' \
  -F Origin=http://www.google.com
```
* **Allowed** CORS result (please set `cors.Debug: true`)
  ```shell
  [cors] 2024/06/16 23:53:13 Handler: Preflight request
  [cors] 2024/06/16 23:53:13 Preflight response headers: map[Access-Control-Allow-Methods:[GET] Access-Control-Allow-Origin:[http://other-domain.com] Access-Control-Max-Age:[6000] Vary:[Origin Access-Control-Request-Method Access-Control-Request-Headers]]
  ```
* **NOT Allowed** CORS result (please set `cors.Debug: true`)
  ```shell
  [cors] 2024/06/16 23:52:13 Handler: Preflight request
  [cors] 2024/06/16 23:52:13 Preflight aborted: origin 'http://other-domain.com' not allowed
  ```

## Create new module:
kalau ada yg kurang
```
    go get github.com/99designs/gqlgen@v0.17.45
```

1. Create graphqls file for the module: edit ```graph/graphqls/{module_name}.graphqls```

3. Generate model & resolvers:
    ```
    go run github.com/99designs/gqlgen generate
    ```
4. Resolvers implementation: edit ```graph/resolvers/{module_name}.resolvers.go```

6. Run:
    ```
    go run server.go 8080
    ```

## Penambahan db Sequence connector
*  config.yaml :
   ```
   dbsequser: root
   dbseqpass: 
   dbseqhost: localhost
   dbseqport: 3306
   dbseqname: his_seq
   ```
*  pemakaian
   ```
   db.Handle.Query(q) 
   di ganti menjadi
   db.HandleSeq.Query(q)
   
   db.HandleSeq.Exec(q)
   diganti menjadi 
   db.HandleSeq.Exed(q)
   ```
