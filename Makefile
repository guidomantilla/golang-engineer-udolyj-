.PHONY: phony
phony-goal: ; @echo $@

prepare:
	go install github.com/kisielk/godepgraph@latest
	go install github.com/incu6us/goimports-reviser/v3@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/cweill/gotests/gotests@latest
	go mod download
	go mod tidy

update-dependencies:
	go get -u ./...
	go get -t -u ./...
	go mod tidy

sonarqube:
	sonar-scanner

run-vanilla:
	go run vanilla/main.go serve

