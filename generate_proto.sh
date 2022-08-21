mkdir -p ./proto
python3 -m grpc_tools.protoc --proto_path=. ./communication.proto --python_out=./proto --grpc_python_out=./proto