#!/bin/sh

go get
go get -v github.com/axw/gocov/gocov
go get -v github.com/mattn/goveralls
go get -v gopkg.in/check.v1

go fmt github.com/stojg/vivere/.
go test github.com/stojg/vivere/.

