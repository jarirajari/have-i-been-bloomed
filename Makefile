#!/bin/bash

N := 501636842

# With Docker NEVER download pwned passwords file, it will be copied as "pwned-passwords.txt.7z"

test-filter:
	bloom --gzip create -p 1e-6 -n 100 pwned-passwords.bloom.test.gz
	7z x pwned-passwords.txt.7z -so | awk -F":" '{print $$1}' | head -n 100 | bloom --gzip insert pwned-passwords.bloom.test.gz

# Create the Bloom filter
pwned-passwords.bloom:
	bloom --gzip create -p 1e-6 -n ${N} pwned-passwords.bloom.gz
	7z x pwned-passwords.txt.7z -so | awk -F":" '{print $$1}' | bloom --gzip insert pwned-passwords.bloom.gz

bloom-filter: pwned-passwords.bloom

bloom-tool:
	go get github.com/dcso/bloom
	go install github.com/dcso/bloom/bloom

test: bloom-tool server test-filter

run:
	hibb

run-test:
	hibb -f pwned-passwords.bloom.test.gz

server:
	go get ./...
	go install ./...

all: bloom-tool server bloom-filter

.DEFAULT_GOAL := all

