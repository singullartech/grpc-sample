.PHONY: client
client:
		@docker build -t sample/client ./grpc-client && docker run --network sample --rm --name grpc-client -it -v ${PWD}:/app sample/client sh