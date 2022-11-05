# compile: python3 -m grpc_tools.protoc --proto_path=. ./communication.proto --python_out=./proto --grpc_python_out=./proto
import os
import sys
from concurrent import futures
import logging
import socket
import grpc

from proto.communication_pb2 import GPIOReply
from proto.communication_pb2_grpc import CommunicatorServicer, add_CommunicatorServicer_to_server


class CommunicatorServicer(CommunicatorServicer):
    def outletOn(self, request, context):
        print("Received SayHello",request.outletId)
        return GPIOReply(on=True)

#def SayHello(self, request, context):


# todo problem: innerhalb des files fehlt einmal der richtige import mit.proto
def serve(address):
    print("Starting Server...")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_CommunicatorServicer_to_server(
        CommunicatorServicer(), server)
    concataddress=address+":8010"
    print("Server at: ",concataddress)
    server.add_insecure_port(concataddress)
    server.start()
    print("Waiting for connections...")
    server.wait_for_termination()
    print("Terminated.")

def get_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.settimeout(0)
    try:
        # doesn't even have to be reachable
        s.connect(('10.255.255.255', 1))
        IP = s.getsockname()[0]
    except Exception:
        IP = '127.0.0.1'
    finally:
        s.close()
    return IP

if __name__ == '__main__':
    print(get_ip())
    logging.basicConfig()
    serve(get_ip())
