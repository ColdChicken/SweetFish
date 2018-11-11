.PHONY: all be nimo fe webpack clean

GOPATH :=
ifeq ($(OS),Windows_NT)
	GOPATH := $(CURDIR)/_vender;$(CURDIR)
else
	GOPATH := $(CURDIR)/_vender:$(CURDIR)
endif

export GOPATH

all: be nimo fe

be:
	go install be/cmd/sweetfish

nimo:
	go install be/cmd/nimo

webpack:
	cd src/fe && npm run build

fe: webpack
	echo '{{define "index"}}' > src/fe/dist/_index.html
	cat src/fe/dist/index.html >> src/fe/dist/_index.html
	echo "{{end}}" >> src/fe/dist/_index.html
	rm -rf src/fe/dist/index.html
	rm -rf src/fe/dist/robots.html
	mv src/fe/dist/_index.html src/fe/dist/index.html
	cp src/fe/src/assets/login.html src/fe/dist/login.html

clean:
	rm -rf bin
	rm -rf pkg
