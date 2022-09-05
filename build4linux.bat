go env -w GOOS=linux

go build -ldflags "-w" -o mindoc main.go

go env -w GOOS=windows
