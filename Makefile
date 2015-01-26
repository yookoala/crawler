# This Makefile is basically for unit testing only
# To build and install this library, please use normal commands
# such as "go get" or "go install"

export GOPATH=$(shell pwd)/_gopath
export DATA=$(shell pwd)/_data
export EXAMPLE=$(shell pwd)/_examples

all: fmt test

test: test.main test.example.sqlite3

test.main:
	@echo "Main Tests"
	@echo "----------"
	go version
	go test
	@echo

test.example.sqlite3: \
	_test.database \
	cachedfetcher \
	_examples/run-all
	@echo "Run Examples on Sqlite3"
	@echo "-----------------------"
	./_examples/run-all -driver "sqlite3" -db "file:./_data/test.sqlite3.db"
	@echo

test.example.mysql: \
	cachedfetcher \
	_examples/run-all
	@echo "Run Examples on MySQL / MariaDB"
	@echo "-------------------------------"
	./_examples/run-all -driver "mysql" -db "${MYSQL}"
	@echo

fmt:
	@echo "Format the source files"
	@echo "-----------------------"
	go fmt
	cd _examples && go fmt
	@echo

clean:
	rm -Rf _gopath

_examples/run-all: \
	_gopath/src \
	_gopath/src/github.com/mattn/go-sqlite3 \
	_gopath/src/github.com/go-sql-driver/mysql \
	cachedfetcher
	@echo "Build Example(s) runner"
	@echo "-----------------------"
	cd _examples && go build -o ${EXAMPLE}/run-all
	@echo

_gopath/src:
	@echo "Create testing GOPATH"
	@echo "---------------------"
	mkdir -p _gopath/src
	@echo

cachedfetcher: _gopath/src/github.com/yookoala/cachedfetcher
	@echo "Install cachedfetcher"
	@echo "-------------------"
	rm -Rf _gopath/pkg/*/github.com/yookoala
	go install github.com/yookoala/cachedfetcher
	@echo

_gopath/src/github.com/yookoala/cachedfetcher:
	@mkdir -p _gopath/src/github.com/yookoala
	@cd _gopath/src/github.com/yookoala && ln -s ../../../../. cachedfetcher

_gopath/src/github.com/mattn/go-sqlite3:
	@echo "Install go-sqlite3"
	@echo "------------------"
	sqlite3 --version
	go get -u github.com/mattn/go-sqlite3
	@echo

_gopath/src/github.com/go-sql-driver/mysql:
	@echo "Install go-sql-driver/mysql"
	@echo "---------------------------"
	go get -u github.com/go-sql-driver/mysql
	@echo

_test.database:
	@echo "Create Example Database"
	@echo "-----------------------"
	cat _data/setup_sqlite3.sql | sqlite3 _data/test.sqlite3.db
	@echo

.PHONY: test test.main test.example.sqlite3 test.example.mysql
.PHONY: _test.database cachedfetcher clean
