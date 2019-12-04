import io

import asyncio

from grpclib.utils import graceful_exit
from grpclib.server import Server

import stats_pb2
import stats_grpc

from adapter import LogItemAdapter
from ploting import get_month_stat
from ploting import get_category_stat
from ploting import get_month_amount_stat

adapter = LogItemAdapter()


class Stater(stats_grpc.StatsBase):

    @staticmethod
    async def send_plot(plt, stream):
        fig = plt.get_figure()

        buf = io.BytesIO()
        fig.savefig(buf, format='png', bbox_inches="tight")
        buf.seek(0)

        plt.cla()
        fig.clear()

        await stream.send_message(stats_pb2.ImageMessage(res=buf.getvalue()))

    async def GetMonthStat(self, stream):
        request: stats_pb2.LogItemQueryMessage = await stream.recv_message()
        data = adapter.get_items_as_dict(request.LogMessagesAggregated)
        plt = get_month_stat(data)

        await self.send_plot(plt, stream)

    async def GetMonthAmountStat(self, stream):
        request: stats_pb2.LogItemQueryMessage = await stream.recv_message()
        data = adapter.get_items_as_dict(request.LogMessagesAggregated)
        plt = get_month_amount_stat(data)

        await self.send_plot(plt, stream)

    async def GetCategoryStat(self, stream):
        request: stats_pb2.LogItemQueryMessage = await stream.recv_message()
        data = adapter.get_items_as_dict(request.LogMessagesAggregated)
        plt = get_category_stat(data)

        await self.send_plot(plt, stream)


async def main(*, host='0.0.0.0', port=50051, loop=None):
    loop = loop or asyncio.get_event_loop()

    server = Server([Stater()], loop=loop)
    with graceful_exit([server], loop=loop):
        await server.start(host, port)
        print(f'Serving on {host}:{port}')
        await server.wait_closed()


if __name__ == '__main__':
    asyncio.run(main())
