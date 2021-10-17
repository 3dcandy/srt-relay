srtrelay:
	HOME=$$(pwd) git config --global http.sslVerify false
	mkdir -p $$(pwd)/gopath
	HOME=$$(pwd) GOPATH=$$(pwd)/gopath go build -o srtrelay

install: srtrelay
	mkdir -p $$(pwd)/debian/srtrelay/usr/bin
	install -m 0755 srtrelay $$(pwd)/debian/srtrelay/usr/bin 

