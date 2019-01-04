.PHONY: clean
clean:
	rm -rf dedupe

.PHONY: test
test:
	go test -v -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html
	open cover.html
