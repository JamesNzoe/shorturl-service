NAME = "shorturl"
DEFAULT_TAG = "shorturl-go:latest"
PACKAGE = "github.com/fpay/lehuipay-shorturl-go"
MAIN = "$(PACKAGE)/entry"
DEFAULT_BUILD_TAG = "1.12-alpine"

TEST_FLAGS= -v ./...
BUILD_FLAGS= -v -o $(NAME) entry/main.go
LINUX_ENV_OS=GOOS=linux
LINUX_ENV_ARCH=GOARCH=amd64
LINUX_ENVS=$(LINUX_ENV_OS) $(LINUX_ENV_ARCH)
LINUX_FLAGS= -mod vendor --ldflags '-extldflags "-static"'

REMOTE_IMAGE = "ccr.ccs.tencentyun.com/fpay/shorturl-go"
REMOTE_TAG = "$(shell git tag -l --sort=-v:refname|head -1|sed -e 's/^v//g')"
ifeq "$(REMOTE_TAG)" ""
	REMOTE_TAG = "latest"
endif
REMOTE_IMAGE_TAG = "$(REMOTE_IMAGE):$(REMOTE_TAG)"

ifeq "$(MODE)" "dev"
	REMOTE_TAG = "1.0.0"
endif

ifeq "$(BUILD_TAG)" ""
	BUILD_TAG = $(DEFAULT_BUILD_TAG)
endif

CL_RED  = "\033[0;31m"
CL_BLUE = "\033[0;34m"
CL_GREEN = "\033[0;32m"
CL_ORANGE = "\033[0;33m"
CL_NONE = "\033[0m"

define color_out
	@echo $(1)$(2)$(CL_NONE)
endef


docker-build:
	$(call color_out,$(CL_BLUE),"Building binary in docker ...")
	@docker run --rm -v "$(PWD)":/go/src/$(PACKAGE) \
		-w /go/src/$(PACKAGE) \
		jerray/golang:$(BUILD_TAG) \
		go build -v -o $(NAME) $(MAIN)
	$(call color_out,$(CL_GREEN),"Building binary ok")

docker: docker-build
	$(call color_out,$(CL_BLUE),"Building docker image ...")
	@docker build -t $(DEFAULT_TAG) .
	$(call color_out,$(CL_GREEN),"Building docker image ok")



push: docker
	@docker tag $(DEFAULT_TAG) $(REMOTE_IMAGE_TAG)
	$(call color_out,$(CL_BLUE),"Pushing image $(REMOTE_IMAGE_TAG) ...")
	@docker push $(REMOTE_IMAGE_TAG)
	$(call color_out,$(CL_ORANGE),"Done")

build:
	@go build -v -o $(NAME) $(MAIN)

test:
	go test $(TEST_FLAGS)

linux:
	@GOOS=linux GOARCH=amd64 go build -v -o $(NAME) $(MAIN)

.PHONY: all
all:
	build
