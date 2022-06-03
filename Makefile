EXECUTABLE=koe
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
VERSION=$(git describe --tags --always --long --dirty)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)


$(WINDOWS):
	env GOOS=windows go build -v -o $(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go

$(LINUX):
	env GOOS=linux go build -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go

$(DARWIN):
	env GOOS=darwin go build -v -o $(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go
