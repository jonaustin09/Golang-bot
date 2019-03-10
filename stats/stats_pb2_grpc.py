# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import stats_pb2 as stats__pb2


class StatsStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.GetAllTimeByMonthStat = channel.unary_unary(
        '/stats.Stats/GetAllTimeByMonthStat',
        request_serializer=stats__pb2.LogItemQueryMessage.SerializeToString,
        response_deserializer=stats__pb2.ImageMessage.FromString,
        )
    self.GetAllTimeCategoryStat = channel.unary_unary(
        '/stats.Stats/GetAllTimeCategoryStat',
        request_serializer=stats__pb2.LogItemQueryMessage.SerializeToString,
        response_deserializer=stats__pb2.ImageMessage.FromString,
        )


class StatsServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def GetAllTimeByMonthStat(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetAllTimeCategoryStat(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_StatsServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'GetAllTimeByMonthStat': grpc.unary_unary_rpc_method_handler(
          servicer.GetAllTimeByMonthStat,
          request_deserializer=stats__pb2.LogItemQueryMessage.FromString,
          response_serializer=stats__pb2.ImageMessage.SerializeToString,
      ),
      'GetAllTimeCategoryStat': grpc.unary_unary_rpc_method_handler(
          servicer.GetAllTimeCategoryStat,
          request_deserializer=stats__pb2.LogItemQueryMessage.FromString,
          response_serializer=stats__pb2.ImageMessage.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'stats.Stats', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))