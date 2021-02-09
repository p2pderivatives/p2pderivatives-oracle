# P2PDerivatives Oracle

This repository hosts the code for an oracle intended to be used for Discreet Log Contracts (DLC).
It design to work with price feeds, serving data on a regular basis that can be decided through a configuration.
See [here](./test/config/default.release.yml) for an example configuration.
At the moment only the CryptoCompare data feed is supported, but adding support for other feeds can be done by implementing the `DataFeed` interface.
It is possible to run the oracle without a CryptoCompare API key but only a few requests can be made.

The description of the API can be found [here](./api/README.md).

## Trying it out

To quickly try out the oracle:

```bash
make gen
docker-compose up
```

The oracle should then be running on port 8080 and you can access it through your web browser.
To see the list of assets (by default btc/usd and btc/jpy):
- localhost:8080/asset

To get an announcement for an event (replace the date as adequate):
- http://localhost:8080/asset/btcusd/announcement/2021-02-15T05:32:00Z

To get an attestation for an event (replace the date as adequate):
- http://localhost:8080/asset/btcusd/attestation/2021-02-08T05:32:00Z

## Requirements

Confirmed working with `Go` v1.14.

## Getting started

Run `make setup` to setup the repository.
You will need to setup a `postgresql` database connection in the configuration file (or via environment variables).  
You can easily setup a running database using `docker-compose up db`  
Once that is done, the server can be run locally using `make run-local-server`.

## Integration Test

The integration tests uses the go REST client library [`Resty`](https://github.com/go-resty/resty).
You can run integration tests by using `make integration-test` on local (with an oracle server running).
or by using  
`gotestsum -- -tags=integration -parallel=4 ./test/integration/... -config-file-name <config-file> -oracle-base-url <oracle-url>`

## Running using Docker

### Docker Compose

You can easily start and build the docker environment using `docker-compose`  
To build from scratch the server use: `docker-compose up --build`

### Building the image

In the root of the repository:

`docker build -t p2pd-oracle .`

The name `p2pd-oracle` can be changed to any other docker compliant name.

### Running the container

Once built, you can start running the server (a database connection is necessary):
`docker-compose up`
To run the container in the background use the `-d` flag.

### Docker configuration

You can override any variable from the configuration file `.yaml` using `environment variables` with this format `APPNAME_MY_PATH_TO_VARIABLE`  
Example:  

- To override the `database.host` property in configuration file : `P2PDORACLE_DATABASE_HOST=mynewhost`
