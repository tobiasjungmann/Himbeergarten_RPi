rm -rf ./proto
mkdir -p ./proto
#python3 -m grpc_tools.protoc -I. ./communication.proto --python_out=./proto --grpc_python_out=./proto
python3 -m grpc_tools.protoc -I proto --python_out=. --grpc_python_out=. proto/some/folder/*.proto
touch ./proto/__init__.py