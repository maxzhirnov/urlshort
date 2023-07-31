.PHONHY: test

APP_NAME := urlshort

test:
	go test ./... -v -count=1
