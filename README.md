# p2pderivatives-oracle

Repository for the P2PDerivatives oracle

## Requirements

Confirmed working with `Go` v1.14.  
`Postman/newman` for integration testing.

## Getting started

Run `make setup` to setup the repository.
You will need to setup a `postgresql` database connection in the configuration file (or via environment variables).  
You can easily setup a running database using `docker-compose up db`  
Once that is done, the server can be run locally using `make run-local-server`.

## Integration Test

The integration tests uses the go REST client library [`Resty`](https://github.com/go-resty/resty)

## Running using Docker

### Docker Compose

You can easily start and build the docker environment using `docker-compose`  
To build from scratch the server use: `docker-compose up --build`

### Building the image

In the root of the repository:

`docker build -t p2pd-server .`

The name `p2pd-server` can be changed to any other docker compliant name.

### Running the container

Once built, you can start running the server (a database connection is necessary):
`docker-compose up`
To run the container in the background use the `-d` flag.

### Docker configuration

You can override any variable from the configuration file `.yaml` using `environment variables` with this format `APPNAME_MY_PATH_TO_VARIABLE`  
Exemple:  

- To override the `database.host` property in configuration file : `P2PD_DATABASE_HOST=mynewhost`
