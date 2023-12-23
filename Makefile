.PHONY:docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook
	@docker rmi -f my_space/webook:v0.0.1
	@docker build -t my_space/webook:v0.0.1 .