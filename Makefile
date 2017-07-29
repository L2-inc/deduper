deduper: main.go
	go build
test: main.go main_test.go
	go test -v -cover
clean:
	rm deduper
