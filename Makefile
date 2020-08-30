.PHONY: build clean deploy

build:
	cd app && env GOOS=linux packr2 build -ldflags="-s -w" -o ../bin/app

clean:
	rm -rf ./bin

update:
	go get -u github.com/elliotchance/ok@$(shell curl -s https://api.github.com/repos/elliotchance/ok/releases/latest | jq -r .tag_name)

deploy: clean update build
	sls create_domain
	sls deploy --verbose
