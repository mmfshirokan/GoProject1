generate_grpc:
	cd ./proto &&  protoc --go_out=. --go-grpc_out=. *.proto && cd ../
docker: 
	docker-compose up --build