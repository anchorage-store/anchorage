.PHONY: test

test:
	go test ./... -v

usr:
	go run ./cmd/usrhttpd
