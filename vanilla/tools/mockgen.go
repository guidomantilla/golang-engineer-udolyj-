package tools

//go:generate mockgen -package config -destination ../pkg/config/mocks.go github.com/sethvargo/go-envconfig Lookuper
//go:generate mockgen -package datasource -destination ../pkg/datasource/mocks.go -source ../pkg/datasource/types.go
//go:generate mockgen -package=environment -source ../pkg/environment/types.go -destination ../pkg/environment/mocks.go
//go:generate mockgen -package=log -source ../pkg/log/types.go -destination ../pkg/log/mocks.go
//go:generate mockgen -package=properties -source ../pkg/properties/types.go -destination ../pkg/properties/mocks.go
//go:generate mockgen -package=security -source ../pkg/security/types.go -destination ../pkg/security/mocks.go
//go:generate mockgen -package=server -destination ../pkg/server/mocks.go github.com/qmdx00/lifecycle Server

//go:generate mockgen -package=rpc -destination ../core/endpoints/rpc/mocks.go -source ../core/endpoints/rpc/types.go
//go:generate mockgen -package=services -destination ../core/services/mocks.go -source ../core/services/types.go
