go env -w CGO_ENABLED=0
go build -o .\LogCenter.exe -ldflags "-s -w" .\main.go
upx -9 .\LogCenter.exe