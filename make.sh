#!/bin/sh

gofmt -s -w ./ai
gofmt -s -w ./engine
gofmt -s -w ./net
gofmt -s -w ./physics
gofmt -s -w ./vec

cd ai && go test
cd ../engine && go test
cd ../net && go test
cd ../physics && go test
cd ../vec && go test
go get
