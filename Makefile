NAME     := qiiotd
VERSION  := 0.2.0
REVISION := $(shell git rev-parse --short HEAD)
SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\""

bin/$(NAME): $(SRCS) format
	GO111MODULE=on go build $(LDFLAGS) -o bin/$(NAME)

linux: $(SRCS) format
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(NAME)

format:
	GO111MODULE=on go fmt $(SRCS)

clean:
	rm -rf bin/*

install:
	GO111MODULE=on go install $(LDFLAGS)

.PHONY: format, clean, install
