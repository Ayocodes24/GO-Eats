PROTO_DIR  := proto
PKG_DIR    := pkg/proto
GOPATH_BIN := $(shell go env GOPATH)/bin
PATH       := $(GOPATH_BIN):$(PATH)

.PHONY: proto build-all run-all setup-dbs

proto:
	protoc --go_out=$(PKG_DIR)/user       --go_opt=paths=source_relative \
	       --go-grpc_out=$(PKG_DIR)/user  --go-grpc_opt=paths=source_relative \
	       --proto_path=$(PROTO_DIR)/user \
	       $(PROTO_DIR)/user/user.proto
	protoc --go_out=$(PKG_DIR)/restaurant       --go_opt=paths=source_relative \
	       --go-grpc_out=$(PKG_DIR)/restaurant  --go-grpc_opt=paths=source_relative \
	       --proto_path=$(PROTO_DIR)/restaurant \
	       $(PROTO_DIR)/restaurant/restaurant.proto
	@echo "Proto generation complete."

build-all:
	go build ./cmd/user-svc/...
	go build ./cmd/restaurant-svc/...
	go build ./cmd/order-svc/...
	go build ./cmd/gateway/...
	@echo "All services built."

setup-dbs:
	bash scripts/setup-dbs.sh

run-all:
	bash scripts/run-all.sh
