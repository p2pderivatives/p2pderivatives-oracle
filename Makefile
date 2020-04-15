setup: install deps gen-ssl-certs
	echo "setup done"

install:
	go mod download

deps:
	go mod tidy

vendor:
	go mod vendor

gen:
	@make gen-ssl-certs

gen-ssl-certs:
	mkdir -p certs/db
	$(eval CERT_TEMP=$(shell mktemp -d))
	openssl req -new -text -passout pass:abcd -subj /CN=localhost -out ${CERT_TEMP}/db.req -keyout ${CERT_TEMP}/privkey.pem
	openssl rsa -in ${CERT_TEMP}/privkey.pem -passin pass:abcd -out certs/db/db.key
	openssl req -x509 -in ${CERT_TEMP}/db.req -text -key certs/db/db.key -out certs/db/db.crt
	chmod 600 certs/db/db.key

oracle:
	mkdir -p bin
	go build -o ./bin/oracle ./cmd/p2pdoracle/main.go

unit-test:
	go test $(shell go list ./... | grep -v /integration/)

integration-test:
	go test ./test/integration/... -appname p2pdoracle -e integration -abs-config $(shell pwd)/test/config

run-server-local:
	@make server
	./bin/oracle -config ./test/config -appname p2pdoracle -e integration -migrate

docker:
	docker build -t docker.pkg.github.com/cryptogarageinc/p2pderivatives-oracle/server .

help:
	@make2help $(MAKEFILE_LIST)
