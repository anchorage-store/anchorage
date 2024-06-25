.PHONY: test

include ./conf/local.env
export

test:
	go test ./... -v

usr:
	go run ./cmd/usrhttpd

usr-migrate:
	go run ./cmd/usrhttpd -migrate=true
