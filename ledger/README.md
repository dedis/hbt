# Introduction

This folder contains the Sawtooth transaction families. For more about
transaction families, see the [Sawtooth
documentation](https://sawtooth.hyperledger.org/docs/1.2/app_developers_guide/overview.html):

> Sawtooth separates the application level from the core system level with
> **transaction families**, which allows application developers to write in the
> languages of their choice. Each application defines the custom transaction
> families for its unique requirements.

In short, transaction families are our custom smart contracts that will be
deployed on the Sawtooth network.

In a transaction family folder, you will find the following folders:

- **cmd**: contains the CLIs to 1) register the transaction family, 2) send
  transactions
- **handler**: holds the smart contract's logic
- **payload**: defines the contend of the transaction sent to the smart contract
- **state**: defines the state of the smart contract

## Prerequisites

1) Install **docker** and **docker-compose**:

https://docs.docker.com/get-docker/

2) Install **zeromq**:

https://zeromq.org/download/

3) Install **openssl**, and make sure PKG_CONFIG_PATH is set. For example on Mac:

```sh
brew install openssl
export PKG_CONFIG_PATH="/usr/local/opt/openssl@3/lib/pkgconfig"
```

## Setup

1) Run docker-compose:

```sh
docker-compose -f sawtooth-default.yml up
```

If is was previously running, stop it first with:

```sh
docker-compose -f sawtooth-default.yml down
```

2) Run the smart contract handler:

```sh
cd cmd/handler
go run . -vv
```

3) Use the client CLI to interact with the smart contract. For example:

```sh
cd cmd/client
go run . -k /tmp/private.key add --id "XX"
```