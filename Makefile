.PHONHY: test, test_iter1

APP_NAME := urlshort

test:
	go test ./... -v -count=1

test_iter1:
	go build -o ./bin/shortener ./cmd/shortener
	shortenertest -test.v -test.run=^TestIteration1 -binary-path=cmd/shortener/shortener