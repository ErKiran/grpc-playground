### Compile GRPC
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/calculator.proto

protoc -I . --grpc-gateway_out /proto \
    --grpc-gateway_opt paths=source_relative \
    proto/calculator.proto