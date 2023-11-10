go env -w GOOS= linux
go build -o tx-erp ./main.go
docker image build -t tx-erp:v0.0.1 .