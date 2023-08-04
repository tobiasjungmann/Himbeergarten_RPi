rm -rf ./esp_humidity/proto-ino
mkdir -p ./esp_humidity/proto-ino

protoc --plugin=protoc-gen-nanopb=/home/tobias/Downloads/nanopb-0.4.7-linux-x86/generator/protoc-gen-nanopb --experimental_allow_proto3_optional=true --proto_path=../server/proto/ --nanopb_out=./esp_humidity/ --cpp_out=./ humidityStorage.proto