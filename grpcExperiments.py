# compile: python3 -m grpc_tools.protoc --proto_path=. ./communication.proto --python_out=./proto --grpc_python_out=./proto
import os
import sys
from concurrent import futures
import logging

import grpc

from proto.communication_pb2 import Reply
from proto.communication_pb2_grpc import CommunicatorServicer, add_CommunicatorServicer_to_server


class CommunicatorServicer(CommunicatorServicer):
    def GetFeature(self, request, context):
        print(request.msg)
        return Reply(msg="received in server")

# todo problem: innerhalb des files fehlt einmal der richtige import mit.proto
def serve():
    print("Starting Server...")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_CommunicatorServicer_to_server(
        CommunicatorServicer(), server)
    server.add_insecure_port('127.0.0.1:8000')
    server.start()
    print("Waiting for connections...")
    server.wait_for_termination()
    print("Terminated.")


if __name__ == '__main__':
    logging.basicConfig()
    serve()
