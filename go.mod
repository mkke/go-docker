module github.com/mkke/go-docker

go 1.15

replace github.com/mkke/go-mlog => ../go-mlog

require (
	docker.io/go-docker v1.0.0
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/mkke/go-log v0.0.0-20201114112904-fa93cd4c3b73
	github.com/mkke/go-mlog v0.0.0-20201115114057-047aed649499
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1 // indirect
)
