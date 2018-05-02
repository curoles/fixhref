export GOPATH=$(PWD)
export PATH:=$(PATH):$(GOPATH)/bin

.PHONY: all
all:
	@go install github.com/curoles/fixhref
	ls -lh $(GOPATH)/bin

.PHONY: help
help:
	@echo "env    - show environment variables."
	@echo "export - show what variables to export into shell."
	@echo "get    - call 'go get' to download all required Go packages."

.PHONY: env
env:
	@echo Go environment:
	@echo ==============
	@go env
	@echo Shell environment:
	@echo =================
	@echo CWD=$(PWD)
	@echo PATH=$(PATH)
	@echo GOPATH=`go env GOPATH`

.PHONY: export
export:
	@echo export GOPATH=$(GOPATH)
	@echo "export PATH=\$$PATH:\$$GOPATH/bin"

.PHONY: get
get:
	go get golang.org/x/net/html
	#go get golang.org/x/crypto/ssh/terminal
	#go get github.com/braintree/manners
	#go get github.com/mattn/go-sqlite3
	#go install github.com/mattn/go-sqlite3
	#go get -u -v github.com/curoles/go-fun
	#go get -u -v github.com/curoles/answer42

#.PHONY: godoc
#godoc:
#	#go doc -u github.com/curoles/fixhref Foo
#	@echo "<html><head></head><body>" > 1.html
#	godoc -html -index github.com/curoles/fixhref/fixhref >> 1.html
#	@echo "</body></html>" >> 1.html

.PHONY: doc
doc:
	mkdir -p doc/g1
	@echo "<html><head></head><body>" > doc/g1/doc1.html
	markdown $(SOURCE_PATH)/doc/group1/doc1.md >> doc/g1/doc1.html
	@echo "</body></html>" >> doc/g1/doc1.html
