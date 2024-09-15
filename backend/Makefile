build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go 
deploy: build
	sls deploy
clean:
	rm -rf ./bin ./vendor Gopkg.lock ./serverless