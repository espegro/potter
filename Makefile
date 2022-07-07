GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-s -w"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all:clean potter

timespotter:
	$(GOBUILD) potter.go 

clean:
	$(GOCLEAN)
