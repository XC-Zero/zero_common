go env -w GOOS= linux
go env -w ARCH= amd64
go build -o mdns ../mdns.go
docker image build -t mdns:v0.0.1 .