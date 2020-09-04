
BINARY = tamarin
SOURCES = $(wildcard monkey/*.go) $(wildcard monkey/*/*.go)
GO = go

.PHONY: build
build: $(BINARY)

$(BINARY): $(SOURCES)
	$(GO) build -o $@ ./monkey
