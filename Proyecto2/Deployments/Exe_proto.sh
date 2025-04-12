# instalar proto
#go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

#export PATH="$PATH:$(go env GOPATH)/bin" # Exportamos al path

go mod init client.go # .mod

# dependencias
go get google.golang.org/protobuf@latest
go get google.golang.org/grpc@latest

go mod tidy

# Generar Código
protoc --go_out=. --go-grpc_out=. tweet.proto
