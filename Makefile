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

define newline


endef

define HTML_HEADER_
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>$(1)</title>
</head>
<body>

endef

# HTML_HEADER and FOOTER usage example:
# printf "$(call HTML_HEADER,Doc1)"  > doc/g1/doc1.html
# markdown $(SOURCE_PATH)/doc/group1/doc1.md >> doc/g1/doc1.html
# printf "$(HTML_FOOTER)" >> doc/g1/doc1.html
define HTML_HEADER
$(subst $(newline),\n,$(call HTML_HEADER_,$(1)))
endef

define HTML_FOOTER_
</body>
</html>
endef
HTML_FOOTER:=$(subst $(newline),\n,${HTML_FOOTER_})


.PHONY: doc
doc: M2H:=-f markdown -t html -s
doc:
	mkdir -p doc/g1
	pandoc $(SOURCE_PATH)/test/group1/doc1.md $(M2H) -o doc/g1/doc1.html

.PHONY: design_doc
design_doc: SRC:=$(SOURCE_PATH)/fixhref/fixhref.go
design_doc: DST:=doc/design.md
design_doc: DST_HTML:=doc/design.html
design_doc:
	awk '/PROSE_BEGIN/{flag=1;next}/PROSE_END/{flag=0}flag' $(SRC) > $(DST)
	pandoc $(DST) -f markdown -t html -s -o $(DST_HTML)
