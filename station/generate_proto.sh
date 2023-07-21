rm -rf ./proto
mkdir -p ./proto


python -m grpc_tools.protoc \
	--proto_path=. \
	--python_out=./proto/ \
	--grpc_python_out=./proto/ \
	../server/proto/storageServer.proto

RUN python fix_proto_imports.py proto/storageServer_pb2_grpc.py proto.

touch ./proto/__init__.py
