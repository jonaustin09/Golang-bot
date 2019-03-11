### How to deploy?
1) run `docker-compose -f docker-compose.yaml run bot`

How to generate proto:
run command in container: ```protoc -I=stats/proto stats/proto/stats.proto --go_out=plugins=grpc:stats```