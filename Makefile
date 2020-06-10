
run_bot_bash:
	docker-compose -f docker-compose.yaml run --rm bot bash $(ARGS)

run_stats_bash:
	docker-compose -f docker-compose.yaml run --rm stats bash $(ARGS)

gen_proto:
	$(MAKE) run_bot_bash \
	  ARGS='-c "protoc -I=stats/proto stats/proto/stats.proto --go_out=plugins=grpc:stats --go_opt=paths=source_relative"'
	$(MAKE) run_stats_bash \
	  ARGS='-c "poetry run python -m grpc_tools.protoc -I=proto --python_out=ploting --python_grpc_out=ploting proto/stats.proto"'
