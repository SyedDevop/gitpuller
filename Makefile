run:build
	@./bin/gitpuller get SyedDevop/linux-setup ./temp -p

run1:build
	@./bin/gitpuller get SyedDevop/fiyat_lis ./temp -p

run2:build
	@./bin/gitpuller get SyedDevop/pc-inf ./temp -p

run3:build
	@./bin/gitpuller get SyedDevop/large-fil ./temp -p

clean:
	@rm -rf ./temp/.* ./temp/*

build:
	@go build  -o ./bin/
