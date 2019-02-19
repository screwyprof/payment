# GopherPay

<p align="center">
    <img alt="Gopher" src="https://raw.githubusercontent.com/screwyprof/payment/master/assets/gopher.png" width="250">
</p>

<p align="center">
 A show case of DDD/CQRS/Clean Architecture
</p>

[![CircleCI](https://circleci.com/gh/screwyprof/payment/tree/master.svg?style=svg)](https://circleci.com/gh/screwyprof/payment/tree/master)
[![Codecov](https://codecov.io/gh/screwyprof/payment/branch/master/graph/badge.svg)](https://codecov.io/gh/screwyprof/payment)
[![Go Report Card](https://goreportcard.com/badge/github.com/screwyprof/payment)](https://goreportcard.com/report/github.com/screwyprof/payment)
[![Codebeat badge](https://codebeat.co/badges/ad61b532-8ced-4c61-99a2-448fad6950da)](https://codebeat.co/projects/github-com-screwyprof-payment-master)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/screwyprof/payment.svg)

## Getting help
```bash
make help
```

## Building Docker Image
```bash
make docker-build
```

## Running service in Docker
```bash
make docker-run
```

## Building Locally Quickly
To build the application, run tests and install it on a local machine just run the following:
```bash
make
```

Otherwise you can do it manually by running the commands from the sections below.

## Installing dependencies
```bash
make deps
```

## Running unit-tests
```bash
make unit-test
```

## Running integration-tests
```bash
make integration-test
```

## Building the application
```bash
make build
```

## Installing application
The following command will install the application to $GOPATH/bin
> Note: Don't forget to add $GOPATH/bin to your $PATH
```bash
make install
```

## Running the application
If the application was installed correctly just run:
```bash
gopherpay
```

## Exploring the Examples
Just look at the [example](./example/) folder

## Getting API information (Swagger)
After starting the service type in your browser:
`http://localhost:8080/swagger/index.html`
