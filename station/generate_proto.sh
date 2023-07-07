rm -rf ./proto
mkdir -p ./proto
python3 -m grpc_tools.protoc -I. ./storageServer.proto --python_out=./proto --grpc_python_out=./proto
touch ./proto/__init__.py

# manually add .proto in the import statement i communication_pb2_grpc.py