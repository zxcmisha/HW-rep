include .env
export

new-user:
	@export NEW_USER=NO && \
	go run main.go


