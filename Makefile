default: bin

# Misc
branch := $(shell git rev-parse --abbrev-ref HEAD)
ifeq (Darwin,$(shell uname))
	shasum := shasum
	tarxform := -s
else
	shasum := sha1sum
	tarxform := --show-transformed-names --transform=s
endif

# Dist vars
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
binary := zipurls
dist_dir := dist
dist_binary := $(dist_dir)/$(binary)
version := $(shell git describe --always --tags --dirty) # --abbrev=40)
short_sha := $(shell git rev-parse --short $(shell git describe --always --tags))

deps:
	go get -v ./...

# Clean up
.PHONY: clean
clean:
	rm -rf app/
	rm -f $(dist_dir)/*


# Development

# Build linux/amd64 binary.
.PHONY: dist
dist:
	@GOOS=linux GOARCH=amd64 make bin

# Build binary based on local environment's GOOS and GOARCH.
bin: dist-dir
	go build -v -ldflags "-X main.version $(version)" -o $(dist_binary) .
	@echo $(version) > $(dist_dir)/VERSION
	@echo '==>' built $(binary) $(GOOS) $(GOARCH) version $(version)

dist-dir:
	@mkdir -p $(dist_dir)

# builds the slug (./ PREFIX IS REQUIRED)
slug: dist
	tar -c -z -v -f slug.tgz $(tarxform)/$(dist_dir)/app/ ./$(dist_dir)

# fetches the slug from latest green master build on CI
get-slug:
	cart -branch $(branch) slug.tgz

deploy-slug: slug.tgz
	slugger -app stable-eu && slugger -app stable-us

deploy: get-slug deploy-slug

