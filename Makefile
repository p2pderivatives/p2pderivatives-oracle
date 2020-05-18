.PHONY: setup install install-tools deps gen

setup: install install-tools deps gen
	@echo "setup done"

install:
	go mod download

install-tools: install
	$(shell cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %)

deps:
	go mod tidy

vendor:
	go mod vendor

gen: gen-ssl-certs gen-oracle-key

KEY_DIR = certs/oracle
gen-oracle-key:
	mkdir -p $(KEY_DIR)
	openssl rand -base64 32 > $(KEY_DIR)/pass.txt
	openssl ecparam -genkey -name secp256k1 | openssl ec -aes256 -passout file:$(KEY_DIR)/pass.txt -out $(KEY_DIR)/key.pem

DB_CERTS_DIR = certs/db
gen-ssl-certs:
	mkdir -p $(DB_CERTS_DIR)
	$(eval CERT_TEMP=$(shell mktemp -d))
	openssl req -new -text -passout pass:abcd -subj /CN=localhost -out ${CERT_TEMP}/db.req -keyout ${CERT_TEMP}/privkey.pem
	openssl rsa -in ${CERT_TEMP}/privkey.pem -passin pass:abcd -out $(DB_CERTS_DIR)/db.key
	openssl req -x509 -in ${CERT_TEMP}/db.req -text -key $(DB_CERTS_DIR)/db.key -out $(DB_CERTS_DIR)/db.crt
	chmod 600 $(DB_CERTS_DIR)/db.key

MOCK_DIR = test/mock
gen-mock:
	mkdir -p $(MOCK_DIR)
	mockgen -source internal/dlccrypto/crypto_service.go -destination $(MOCK_DIR)/dlccrypto/mock_crypto_service.go
	mockgen -source internal/datafeed/datafeed.go -destination $(MOCK_DIR)/datafeed/mock_datafeed.go

oracle:
	mkdir -p bin
	go build -o ./bin/oracle ./cmd/p2pdoracle/main.go

unit-test:
	gotestsum -- -cover $(shell go list ./... | grep -v /integration/)

integration-test:
	gotestsum -- ./test/integration/... -appname p2pdoracle -e integration -abs-config $(shell pwd)/test/config

run-server-local: oracle
	./bin/oracle -config ./test/config -appname p2pdoracle -e integration -migrate

docker:
	docker build -t docker.pkg.github.com/cryptogarageinc/p2pderivatives-oracle/server .

help:
	@make2help $(MAKEFILE_LIST)
