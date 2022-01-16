hello:
	echo "Hello"

build:
	go build -o bin/main cmd/main.go

run:
	go run cmd/main.go

run-order:
	go run pkg/order/main.go

run-inventory:
	go run pkg/inventory/main.go
	
build-inventory-image:
	docker build -t rfashwal/inventory pkg/inventory/.