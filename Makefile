GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-s -w"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all:clean potter key

potter:
	$(GOBUILD) potter.go 

clean:
	$(GOCLEAN)
key:
	@ssh-keygen -N '' -t ed25519 -f potter.key

