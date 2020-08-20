build:
	set GOOARCH=386
	go build -o bin/synmed-reader-x86.exe --buildmode=exe main.go
	set GOOARCH=amd64
	go build -o bin/synmed-reader-amd64.exe --buildmode=exe main.go
run:
	go run main.go
package:
	make build
	copy LICENSE.md bin\LICENSE.md
	copy README.md bin\README.md
