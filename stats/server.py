import io

import grpc
import time
import hashlib
import stats_pb2
import stats_pb2_grpc
from concurrent import futures

from adapter import LogItemAdapter
from ploting import get_all_time_by_month_stat
from ploting import get_all_time_category_stat

adapter = LogItemAdapter()


class StatsServicer(stats_pb2_grpc.StatsServicer):
    def __init__(self, *args, **kwargs):
        # TODO: add settings file
        self.server_port = 50051

    @staticmethod
    def to_message(plt):
        fig = plt.get_figure()

        buf = io.BytesIO()
        fig.savefig(buf, format='png', bbox_inches="tight")
        buf.seek(0)

        return stats_pb2.ImageMessage(res=buf.getvalue())

    def GetAllTimeByMonthStat(self, request, context):
        data = adapter.get_items_as_dict(request.LogItems)
        plt = get_all_time_by_month_stat(data)

        return self.to_message(plt)

    def GetAllTimeCategoryStat(self, request, context):
        data = adapter.get_items_as_dict(request.LogItems)
        plt = get_all_time_category_stat(data)

        return self.to_message(plt)

    def start_server(self):
        """
        Function which actually starts the gRPC server, and preps
        it for serving incoming connections
        """
        stats_server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10))

        # This line can be ignored
        stats_pb2_grpc.add_StatsServicer_to_server(
            StatsServicer(), stats_server)

        # bind the server to the port defined above
        stats_server.add_insecure_port('[::]:{}'.format(self.server_port))

        # start the server
        stats_server.start()
        print('Stats Server running ...')

        try:
            # need an infinite loop since the above
            # code is non blocking, and if I don't do this
            # the program will exit
            while True:
                time.sleep(60 * 60 * 60)
        except KeyboardInterrupt:
            stats_server.stop(0)
            print('Digestor Server Stopped ...')


StatsServicer().start_server()
