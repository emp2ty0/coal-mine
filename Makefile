include .env
export

generate-server:
	openapi-generator-cli generate \
	-i spec.yaml \
	-g go-server \
	-o ./api-server \
	--additional-properties=packageName=minerapi

run-app:
	go run cmd/server/main.go