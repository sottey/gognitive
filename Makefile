# Location of the relifevc module
RELIFEDIR = examples/relifevc

# Build output
BINARY = $(RELIFEDIR)/relifevc

# Default target
all: build

# Clean build artifacts and Go cache
clean:
	go clean -cache -modcache -i -r
	rm -f $(BINARY)

# Build relifevc from within its own module directory
build:
	cd $(RELIFEDIR) && go build -o relifevc

# Rebuild after clean
rebuild: clean build
