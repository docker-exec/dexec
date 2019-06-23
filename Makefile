.PHONY : all test release

all: | test

test:
	@go test

release:
	@.script/release.sh
