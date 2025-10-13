PROTO_TAG ?= v0.0.11
PROTO_NAME := user.proto
PROTO_PATH := github.com/esclient/protos
PROTO_MODULE := $(shell go list -m)

TMP_DIR := .proto
OUT_DIR := api/userservice

ENV_FILE := .env
DOCKER_PORT := 50125

.PHONY: clean fetch-proto gen-stubs update run stop rebuild logs test-cover sonar sonar-docker

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


# --- Docker commands ---
run: build
	docker run -d \
		--name user-service \
		--env-file $(ENV_FILE) \
		-p ${DOCKER_PORT}:${DOCKER_PORT} \
		--restart unless-stopped \
		user-service
	docker logs -f user-service

build:
	docker build -t user-service .

stop:
	-docker stop user-service
	-docker rm user-service

rebuild: stop
	docker build -t user-service .
	$(MAKE) run

logs:
	docker logs -f user-service

# --- Test/Coverage and Sonar ---
test-cover:
	go test ./... -coverprofile=coverage.out

sonar: test-cover
	sonar-scanner

sonar-docker: test-cover
	@if [ -z "$$SONAR_HOST_URL" ] || [ -z "$$SONAR_TOKEN" ]; then \
		echo "Set SONAR_HOST_URL and SONAR_TOKEN environment variables"; \
		exit 1; \
	fi
	docker run --rm -e SONAR_HOST_URL -e SONAR_TOKEN -v "$(shell pwd):/usr/src" sonarsource/sonar-scanner-cli
