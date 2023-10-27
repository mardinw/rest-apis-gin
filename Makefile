all:
	make build
	make cleanpkg
	make package
	make clean


build:
	@echo "Build the binary"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o apps main.go

cleanpkg:
	@echo "Clean Package"
	rm deployment.zip

package:
	@echo "Create a ZIP file"
	zip deployment.zip main

clean:
	@echo "Cleaning up"
	rm main
