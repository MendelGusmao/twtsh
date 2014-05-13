all:
	@go build

test:
	@find . -name "*_test.go" | xargs dirname | xargs go test $(O)

cov:
	@make test O="$(O) -cover"
