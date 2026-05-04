test:
	go test -cover ./...

test-v:
	go test -cover -v ./...

.PHONY: test test-v