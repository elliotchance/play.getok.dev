.PHONY: build clean deploy

build:
	cd app && env GOOS=linux packr2 build -ldflags="-s -w" -o ../bin/app

clean:
	rm -rf ./bin

deploy: clean build
	sls create_domain
	sls deploy --verbose
