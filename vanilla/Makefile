.PHONY: phony
phony-goal: ; @echo $@

build: validate
	docker compose -f docker/docker-compose.yml -p bank-vanilla up --detach

validate: generate sort-import format vet lint coverage

generate:
	go generate ./...

sort-import:
	goimports-reviser -rm-unused -set-alias -format -recursive cmd
	goimports-reviser -rm-unused -set-alias -format -recursive core
	goimports-reviser -rm-unused -set-alias -format -recursive pkg

format:
	go fmt ./...

vet:
	go vet ./cmd/... ./core/... ./pkg/...

lint:
	golangci-lint run ./cmd/... ./core/... ./pkg/...

test:
	go test -covermode count -coverprofile coverage.out.tmp.01 ./core/... ./pkg/...
	cat coverage.out.tmp.01 | grep -v "mocks.go" > coverage.out

coverage: test
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

