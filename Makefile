PROTO_TAG ?= v0.0.2
PROTO_NAME := user.proto
PROTO_PATH := github.com/esclient/protos
PROTO_MODULE := $(shell go list -m)

TMP_DIR := .proto
OUT_DIR := src/userservice/grpc

.PHONY: clean fetch-proto gen-stubs update

ifeq ($(OS),Windows_NT)
MKDIR    = powershell -Command "New-Item -ItemType Directory -Force -Path"
RM       = powershell -NoProfile -Command "Remove-Item -Path '$(TMP_DIR)' -Recurse -Force"
DOWN     = powershell -Command "Invoke-WebRequest -Uri"
DOWN_OUT = -OutFile
else
MKDIR    = mkdir -p
RM       = rm -rf $(TMP_DIR)
DOWN     = wget
DOWN_OUT = -O
endif

clean:
	$(RM)

fetch-proto:
	$(MKDIR) "$(TMP_DIR)"
	$(DOWN) "https://raw.githubusercontent.com/esclient/protos/$(PROTO_TAG)/$(PROTO_NAME)" $(DOWN_OUT) "$(TMP_DIR)/$(PROTO_NAME)"

gen-stubs: fetch-proto
	$(MKDIR) "$(OUT_DIR)"
	protoc \
		--proto_path="$(TMP_DIR)" \
		--go_out="$(OUT_DIR)" \
		--go_opt=module=$(PROTO_MODULE) \
		--go-grpc_out="$(OUT_DIR)" \
		--go-grpc_opt=module=$(PROTO_MODULE) \
		"$(TMP_DIR)/$(PROTO_NAME)"

update: gen-stubs clean