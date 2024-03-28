# Get version from git hash
git_hash := $(shell git rev-parse --short HEAD || echo 'development')

# Get current date
current_time = $(shell date +"%Y-%m-%d:T%H:%M:%S")

# Add linker flags
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_hash}'

ifeq ($(GOOS),)
	GOOS := linux
endif

ifeq ($(GOARCH),)
	GOARCH := amd64
endif

ifeq ($(PREFIX),)
	PREFIX := /usr/local
endif

.PHONY:
build:
	@echo "Building binaries..."
	GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags=${linker_flags} -o=./bin/opencool-${GOOS}-${GOARCH} ./cmd/opencool/main.go
	ln -sf opencool-${GOOS}-${GOARCH} bin/opencool

clean:
	rm -rf ./bin

install:
	cp ./bin/opencool ${PREFIX}/bin
	mkdir -p /etc/opencool/
	cp ./config.yml /etc/opencool/config.yml
	cp ./opencool.service /etc/systemd/system/opencool.service
	systemctl daemon-reload 