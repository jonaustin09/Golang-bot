b:
	docker-compose -f docker-compose.yaml run --rm bot bash
gen_proto:
	docker-compose -f docker-compose.yaml run --rm bot bash -c "cd /go/src/money/ && protoc -I stats stats/stats.proto --go_out=."
build:
	docker-compose -f docker-compose.yaml run --rm bot bash -c "cd /go/src/money/ && rm --f build && go build -o build"