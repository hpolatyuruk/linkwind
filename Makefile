GCFLAGS=-gcflags=-trimpath=$GOPATH
ASMFLAGS=-asmflags=-trimpath=$GOPATH

build:
	$(eval file := $(FILE))
	$(eval variables := $(shell cat ${file}))
	$(eval combinedflags := $(foreach v,$(variables),-X main.$(v)))
	$(eval combinedflags := $(foreach v,$(variables),-X main.$(v)))
	$(eval LDFLAGS=-ldflags "$(combinedflags)")
ifeq ($(FILE), .env.dev)
	go build $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) -o main ./settings
else
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) -o main ./settings
endif

build-dev:
	make build FILE=".env.dev"

build-staging:   
	make build FILE=".env.staging"

build-prod:
	make build FILE=".env.prod"

