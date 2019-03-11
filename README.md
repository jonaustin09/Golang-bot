### How to deploy?
1) run `docker-compose -f docker-compose.yaml up`

#### How to generate proto:
run command in container: 
```bash
protoc -I=stats/proto stats/proto/stats.proto --go_out=plugins=grpc:stats
```