#!/home/pi/bin/python
# -*- coding: utf-8 -*-
import grpc
import proto.storageServer_pb2
import proto.storageServer_pb2_grpc

if __name__ == '__main__':
    print("starting")
    print("Will try to greet world ...")
    with grpc.insecure_channel('192.168.178.97:12346') as channel:
        stub = proto.storageServer_pb2_grpc.StorageServerStub(channel)
        response = stub.storeHumidityEntry(proto.storageServer_pb2.StoreHumidityRequest(requestNumber=1,humidity=7))
    print("Greeter client received.")