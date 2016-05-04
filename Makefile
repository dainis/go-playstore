VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

updatedeps:
	go get -u github.com/kardianos/govendor
	govendor fetch +vendor

# test runs the unit tests and vets the code
test:
	go test -timeout=30s -parallel=4
	@$(MAKE) vet

# testrace runs the race checker
testrace:
	go test -race

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go list -f '{{if gt (len .TestGoFiles) 0}}"go test -covermode count -coverprofile {{.Name}}.coverprofile -coverpkg ./... {{.ImportPath}}"{{end}}' ./... | xargs -I {} bash -c {}
	gocovmerge `ls *.coverprofile` > coverage.out
	go tool cover -html=coverage.out
	rm coverage.out
	rm *.coverprofile

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) . ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for reviewal."; \
	fi

updateproto:
	wget https://raw.githubusercontent.com/dweinstein/google-play-proto/master/googleplay.proto -O protobuf/googleplay.proto
	protoc -I=protobuf --go_out=protobuf protobuf/googleplay.proto

.PHONY: updatedeps vet testrace test cover
