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