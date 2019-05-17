#### How to generate proto:
run command in container: 
```bash
poetry run python -m grpc_tools.protoc -I=proto --python_out=ploting --python_grpc_out=ploting proto/stats.proto
```
