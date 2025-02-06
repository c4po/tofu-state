BINARY_NAME=tofu-state-backend
# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildTime=${BUILD_TIME}"

# Build directories
BUILD_DIR=build
LINUX_AMD64_DIR=$(BUILD_DIR)/linux_amd64
DARWIN_AMD64_DIR=$(BUILD_DIR)/darwin_amd64
DARWIN_ARM64_DIR=$(BUILD_DIR)/darwin_arm64

# Build for Linux (amd64)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(LINUX_AMD64_DIR)/$(BINARY_NAME) ./cmd/server

# Build for macOS (both amd64 and arm64)
build-darwin:
	# CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DARWIN_AMD64_DIR)/$(BINARY_NAME) ./cmd/server
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DARWIN_ARM64_DIR)/$(BINARY_NAME) ./cmd/server

# Install dependencies
deps:
	$(GOGET) -v ./...

# Run go fmt against code
fmt:
	$(GOCMD) mod tidy
	$(GOCMD) fmt ./...

# Run go vet against code
vet:
	$(GOCMD) vet ./...
