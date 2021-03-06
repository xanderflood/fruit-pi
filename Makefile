commands := $(shell ls cmd/ | awk '{split($$0,a,"/"); print a[1]}' | tr '\n' ' ')

$(commands):
	@echo "Building command $@"
	go build cmd/$@/main.go
	mv main build/$@/$@
	cd ./build/$@ && godotenv -f ../../.env ./$@

gen: ; go generate ./...
