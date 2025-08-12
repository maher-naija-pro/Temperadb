# Makefile for TimeSeriesDB - Test and Coverage only

# Go command
GOCMD=go

# Test targets
.PHONY: test
test:
	$(GOCMD) test -v ./...

# Test with coverage
.PHONY: coverage
coverage:
	$(GOCMD) test -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean coverage files
.PHONY: clean
clean:
	rm -f coverage.out coverage.html

# Default target
.DEFAULT_GOAL := test
