default: reviewable

reviewable: build fmt generate test ## Run before committing.

GOBIN = $(shell pwd)/bin
PROVIDER_SRC_DIR := ./internal/provider/...
CONTAINER_COMPOSE_ENGINE ?= $(shell docker compose version >/dev/null 2>&1 && echo 'docker compose' || echo 'docker-compose')

build: ## Build the provider binary.
	go mod tidy
	GOBIN=$(GOBIN) go install

generate: ## Generate documentation.
	PATH="$(GOBIN):$(PATH)" go generate --tags tfplugindocs ./...

ifdef RUN
TESTARGS += -test.run $(RUN)
endif

test: ## Run unit tests.
	go test -cover $(TESTARGS) $(PROVIDER_SRC_DIR)

fmt: tool-golangci-lint tool-terraform tool-shfmt ## Format files and fix issues.
	gofmt -s -w -e .
	$(GOBIN)/golangci-lint run --build-tags acceptance --fix
	$(GOBIN)/terraform fmt -recursive -list ./examples ./playground
	$(GOBIN)/shfmt -l -s -w ./examples

lint: lint-golangci lint-examples-tf lint-examples-sh lint-generated

lint-golangci: tool-golangci-lint ## Run golangci-lint linter (same as fmt but without modifying files).
	PATH="$(GOBIN):$(PATH)" golangci-lint run --build-tags acceptance

lint-examples-tf: tool-terraform ## Run terraform linter on examples (same as fmt but without modifying files).
	PATH="$(GOBIN):$(PATH)" terraform fmt -recursive -check -diff ./examples ./playground

lint-examples-sh: tool-shfmt ## Run shell linter on examples (same as fmt but without modifying files).
	PATH="$(GOBIN):$(PATH)" shfmt -l -s -d ./examples

lint-generated: generate ## Check that "make generate" was called. Note this only works if the git workspace is clean.
	@echo "Checking git status"
	@[ -z "$(shell git status --short)" ] || { \
		echo "Error: Files should have been generated:"; \
		git status --short; echo "Diff:"; \
		git --no-pager diff HEAD; \
		echo "Run \"make generate\" and try again"; \
		exit 1; \
	}
	@echo "validate documentation"
	PATH="$(GOBIN):$(PATH)" go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate

LAKEKEEPER_ENDPOINT ?= http://localhost:8181
LAKEKEEPER_AUTH_URL ?= http://localhost:30080/realms/iceberg/protocol/openid-connect/token
LAKEKEEPER_CLIENT_ID ?= lakekeeper-admin
LAKEKEEPER_CLIENT_SECRET ?= KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL

testacc-up: ## Launch a Lakekeeper instance.
	cd run; $(CONTAINER_COMPOSE_ENGINE) up -d
	LAKEKEEPER_ENDPOINT=$(LAKEKEEPER_ENDPOINT) LAKEKEEPER_AUTH_URL=$(LAKEKEEPER_AUTH_URL) LAKEKEEPER_CLIENT_ID=$(LAKEKEEPER_CLIENT_ID) LAKEKEEPER_CLIENT_SECRET=$(LAKEKEEPER_CLIENT_SECRET) ./scripts/await-healthy.sh
	
testacc-down: ## Teardown a Lakekeeper instance.
	cd run; $(CONTAINER_COMPOSE_ENGINE) down --volumes

testacc: ## Run acceptance tests against a Lakekeeper instance.
	TF_ACC=1 LAKEKEEPER_ENDPOINT=$(LAKEKEEPER_ENDPOINT) LAKEKEEPER_AUTH_URL=$(LAKEKEEPER_AUTH_URL) LAKEKEEPER_CLIENT_ID=$(LAKEKEEPER_CLIENT_ID) LAKEKEEPER_CLIENT_SECRET=$(LAKEKEEPER_CLIENT_SECRET) go test -cover --tags acceptance -v $(PROVIDER_SRC_DIR) $(TESTARGS) -timeout 40m

testacc-flakey: ## Run flakey acceptance tests against a Lakekeeper instance.
	TF_ACC=1 LAKEKEEPER_ENDPOINT=$(LAKEKEEPER_ENDPOINT) LAKEKEEPER_AUTH_URL=$(LAKEKEEPER_AUTH_URL) LAKEKEEPER_CLIENT_ID=$(LAKEKEEPER_CLIENT_ID) LAKEKEEPER_CLIENT_SECRET=$(LAKEKEEPER_CLIENT_SECRET) go test -cover --tags flakey -v $(PROVIDER_SRC_DIR) $(TESTARGS) -timeout 40m

testacc-settings: ## Run application settings acceptance tests against a Lakekeeper instance.
	TF_ACC=1 LAKEKEEPER_ENDPOINT=$(LAKEKEEPER_ENDPOINT) LAKEKEEPER_AUTH_URL=$(LAKEKEEPER_AUTH_URL) LAKEKEEPER_CLIENT_ID=$(LAKEKEEPER_CLIENT_ID) LAKEKEEPER_CLIENT_SECRET=$(LAKEKEEPER_CLIENT_SECRET) go test -cover --tags settings -v $(PROVIDER_SRC_DIR) $(TESTARGS) -timeout 40m

# TOOLS
# Tool dependencies are installed into a project-local /bin folder.

tool-golangci-lint:
	@mkdir -p $(GOBIN)
	@[ -f $(GOBIN)/golangci-lint ] || { curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOBIN) v2.2.0; }

tool-shfmt:
	@$(call install-tool, mvdan.cc/sh/v3/cmd/shfmt)

define install-tool
	GOBIN=$(GOBIN) go install $(1)
endef

playground: tool-terraform
	@cd playground; \
		TF_CLI_CONFIG_FILE=./.terraformrc $(GOBIN)/terraform init; \
		$(GOBIN)/terraform apply -auto-approve

playground-destroy:
	@cd playground; \
		$(GOBIN)/terraform destroy -auto-approve; \
		rm -rf ./terraform.tfstate ./terraform.tfstate.backup .terraform .terraform.lock.hcl

TERRAFORM_VERSION = v1.9.8
tool-terraform:
	@# See https://github.com/hashicorp/terraform/issues/30356
	@[ -f $(GOBIN)/terraform ] || { mkdir -p tmp; cd tmp; rm -rf terraform; git clone --branch $(TERRAFORM_VERSION) --depth 1 https://github.com/hashicorp/terraform.git; cd terraform; GOBIN=$(GOBIN) go install; cd ../..; rm -rf tmp; }

clean: playground-destroy testacc-down
	go clean -testcache
	rm -rf bin/
