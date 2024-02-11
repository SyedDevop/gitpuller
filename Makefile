run:build
	@./bin/gitpuller SyedDevop/linux-setup

run1:build
	@./bin/gitpuller SyedDevop/fiyat_list

run2:build
	@./bin/gitpuller SyedDevop/pc-info

run3:build
	@./bin/gitpuller SyedDevop/large-file

clean:
	@rm -rf ./temp/.* ./temp/*

build:
	@go build  -o ./bin/
