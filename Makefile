server.start:
	@go run app/server/*.go

docker.server.build:
	@docker build -f app/server/Dockerfile -t lynlab/luppiter-server .

docker.server.run:
	@docker run --rm -v "$$HOME/.aws:/.aws" -e "HOME=/" -p 8080:8080 lynlab/luppiter-server
