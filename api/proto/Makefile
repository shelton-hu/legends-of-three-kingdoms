# PROTOC_INC_PATH=/usr/local/include
default: generate

PROTO_FILES=$(sort $(wildcard ./*.proto))
GOOGLEAPIS_PATH=third_party/googleapis/
GRPC_GATEWAY_PATH=third_party/grpc-gateway/

PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(GRPC_GATEWAY_PATH)

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  generate          to generate grpc/grpc-gateway/swagger code"
	@echo "  help              to show this help informatrion"

.PHONY: generate
generate:
	rm -rf ../../pkg/pb/*
	rm -rf ../../docs/api.html
	rm -rf ../../docs/api.md
	rm -rf ../../pkg/apigateway/spec/api.swagger.json

	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../../pkg/pb $(PROTO_FILES)
	@echo "generate pb files ok"

	protoc $(PROTOC_FLAGS) --doc_out=../../docs/ --doc_opt=markdown,api.md $(PROTO_FILES)
	protoc $(PROTOC_FLAGS) --doc_out=../../docs/ --doc_opt=html,api.html $(PROTO_FILES)
	@echo "generate docs files ok"

	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:../../pkg/pb $(PROTO_FILES)
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:../../pkg/apigateway/spec $(PROTO_FILES)
	cd ../../pkg/apigateway/spec/ && cat ./*.swagger.json 0.swagger.json | jq --slurp 'reduce .[] as $$item ({}; . * $$item)' | jq -f ./preprocess.jq > ./api.swagger.json
	cd ../../pkg/apigateway/spec/ && find . | grep .swagger.json | grep -v "api" | xargs rm
	@echo "generate apigateway files ok"

	@echo "finished"