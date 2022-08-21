
# compile: python3 -m grpc_tools.protoc --proto_path=. ./communication.proto --python_out=./proto --grpc_python_out=./proto
from concurrent import futures
import logging

import grpc

import communication_pb2_grpc
from communication_pb2 import Reply


class CommunicatorServicer(communication_pb2_grpc.CommunicatorServicer):
    def GetFeature(self, request, context):
        print(request.msg)
        return Reply(msg="received in server")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    communication_pb2_grpc.add_CommunicatorServicer_to_server(
        CommunicatorServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig()
    serve()
