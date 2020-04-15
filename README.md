# p2pderivatives-oracle

Repository for the P2PDerivatives oracle

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
`go test ./test/integration/... -appname p2pdoracle -e integration -abs-config $(pwd)/test/config -oracle-base-url <my-url>`

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
