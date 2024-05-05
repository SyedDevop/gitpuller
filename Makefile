run:build
	@./bin/gitpuller get SyedDevop/Timer-dash ./temp -p

run1:build
	@./bin/gitpuller get SyedDevop/fiyat_list ./temp -p

run2:build
	@./bin/gitpuller get SyedDevop/pc-info ./temp -p

run3:build
	@./bin/gitpuller get SyedDevop/large-file ./temp -p

run4:build
	@./bin/gitpuller get SyedDevop/tax-care-admin-dashboard ./temp -p

run5:build
	@./bin/gitpuller get SyedDevop/linux-setup ./temp -p
clean:
	@rm -rf ./temp/.* ./temp/*

build:
	@go build  -o ./bin/
