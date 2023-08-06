#!/bin/sh

# todo adapt file
# todo check if library is added correctly

rm -rf ./esp_humidity/proto-ino
mkdir -p ./esp_humidity/proto-ino

ln ../server/proto/humidityStorage.options ./humidityStorage.options

GENERATOR_PATH=/home/tobias/Downloads/nanopb-0.4.7-linux-x86/generator/protoc-gen-nanopb
PROTO_PATH=../server/proto/
PROTO_OUT_PATH=./esp_humidity/


#protoc --proto_path=$PROTO_PATH \
#--plugin=protoc-gen-nanopb=$GENERATOR_PATH \
#      --nanopb_out=-v:. \
#      humidityStorage.proto

# to generate c++ --cpp_out=./
protoc --plugin=protoc-gen-nanopb=$GENERATOR_PATH \
      --experimental_allow_proto3_optional=false \
      --proto_path=$PROTO_PATH \
      --nanopb_out=-v:$PROTO_OUT_PATH  \
      humidityStorage.proto

rm ./humidityStorage.options