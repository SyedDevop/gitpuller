run:build
	@./bin/gitpuller get SyedDevop/linux-setup

run1:build
	@./bin/gitpuller get SyedDevop/fiyat_list

run2:build
	@./bin/gitpuller get SyedDevop/pc-info

run3:build
	@./bin/gitpuller get SyedDevop/large-file

clean:
	@rm -rf ./temp/.* ./temp/*

build:
	@go build  -o ./bin/
