EXECUTABLE=koe
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
PI=$(EXECUTABLE)_pi_arm
VERSION=$(git describe --tags --always --long --dirty)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

pi: $(PI) ## Build for Darwin (macOS)


$(WINDOWS):
	CGO_ENABLED=0 GOOS=windows go build -v -o $(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go

$(LINUX):
	CGO_ENABLED=0 GOOS=linux go build -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go

$(DARWIN):
	CGO_ENABLED=0 GOOS=darwin go build -v -o $(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go

$(PI):
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -o $(PI) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go


# GOOS=linux CC="/usr/local/bin/arm-linux-musleabi-gcc" GOARCH=arm CGO_ENABLED=1 go build -o koe_pi_arm -ldflags "-linkmode external -extldflags -static" main.go

