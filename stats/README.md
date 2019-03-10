### Generate proto code:
```
poetry run python -m grpc_tools.protoc \
-I. \
--python_out=. \
--grpc_python_out=. \
stats.proto
```
