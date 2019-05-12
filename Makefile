server.start:
	@go run app/server/*.go

docker.server.build:
	@docker build -f app/server/Dockerfile -t lynlab/luppiter-server .

docker.server.start:
	@docker run --rm -v "$$HOME/.aws:/.aws" -e "HOME=/" -p 1323:1323 lynlab/luppiter-server
