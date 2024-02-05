run:
	@go run . SyedDevop/linux-setup


run1:
	@go run . SyedDevop/gitpuller

run2:
	@go run . SyedDevop/large-file

clean:
	@rm -rf ./temp/*

build:
	@go build  -o ./bin/
