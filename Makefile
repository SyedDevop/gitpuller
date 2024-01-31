run:
	@go run . SyedDevop/linux-setup


run2:
	@go run . SyedDevop/large-file

clean:
	@rm -rf ./temp/*

build:
	@go build  -o ./bin/
