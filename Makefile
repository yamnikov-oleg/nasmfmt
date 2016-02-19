SOURCES=main.go
GOBUILD=go build
UNIXBIN=bin/nasmfmt
WINBIN=bin/nasmfmt.exe

all: $(SOURCES)
	go install

every_platform: linux_64bit linux_32bit windows_64bit windows_32bit macosx_64bit macosx_32bit

linux_64bit:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(UNIXBIN)
	tar czf bin/nasmfmt_linux64.tar.gz $(UNIXBIN)

linux_32bit:
	GOOS=linux GOARCH=386 $(GOBUILD) -o $(UNIXBIN)
	tar czf bin/nasmfmt_linux32.tar.gz $(UNIXBIN)

windows_64bit:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(WINBIN)
	zip bin/nasmfmt_windows64.zip $(WINBIN)

windows_32bit:
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(WINBIN)
	zip bin/nasmfmt_windows32.zip $(WINBIN)

macosx_64bit:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(UNIXBIN)
	tar czf bin/nasmfmt_macosx64.tar.gz $(UNIXBIN)

macosx_32bit:
	GOOS=darwin GOARCH=386 $(GOBUILD) -o $(UNIXBIN)
	tar czf bin/nasmfmt_macosx32.tar.gz $(UNIXBIN)

clean:
	-rm -rf bin
