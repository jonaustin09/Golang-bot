### Generate proto code:
```
poetry run python -m grpc_tools.protoc \
-I=proto \
--python_out=ploting \
--grpc_python_out=ploting \
proto/stats.proto
```
