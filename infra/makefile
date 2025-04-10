gen-proto:
		protoc -I=./api/pb --go_out=./api/pb --go_opt=paths=source_relative --go-grpc_out=./api/pb --go-grpc_opt=paths=source_relative ./api/pb/*.proto

gen-proto-python:
		python -m grpc_tools.protoc -I./pb --python_out=./pb --pyi_out=./pb --grpc_python_out=./pb ./pb/*.proto