#!/usr/bin/env bash

protoc --proto_path=. --go_out=paths=source_relative:. *.proto
