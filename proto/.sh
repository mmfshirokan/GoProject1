#!/bin/sh

protoc --proto_path=/home/andreishyrakanau/projects/project1/GoProject1/proto/ --go_out=/home/andreishyrakanau/projects/project1/GoProject1/proto/ *.proto;

#protoc --proto_path=/home/andreishyrakanau/projects/project1/GoProject1/proto/ --go_out=/home/andreishyrakanau/projects/project1/GoProject1/proto/ /home/andreishyrakanau/projects/project1/GoProject1/proto/userRPC.proto;
#protoc --proto_path=/home/andreishyrakanau/projects/project1/GoProject1/proto/ --go_out=/home/andreishyrakanau/projects/project1/GoProject1/proto/ /home/andreishyrakanau/projects/project1/GoProject1/proto/passwordRPC.proto;
#protoc --proto_path=/home/andreishyrakanau/projects/project1/GoProject1/proto/ --go_out=/home/andreishyrakanau/projects/project1/GoProject1/proto/ /home/andreishyrakanau/projects/project1/GoProject1/proto/tokenRPC.proto;