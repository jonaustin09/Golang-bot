from attr import dataclass

from stats_pb2 import LogItemQueryMessage


@dataclass
class LogItem:
    timestamp: int
    amount: float
    category: str

    def as_dict(self):
        return {
            'timestamp': self.timestamp,
            'amount': self.amount,
            'category': self.category,
        }


class LogItemAdapter:
    @staticmethod
    def get_items_as_dict(log_item_query_message: LogItemQueryMessage):
        res = []
        for i in (LogItem(timestamp=l.CreatedAt, amount=l.Amount,
                          category=l.Category) for l in log_item_query_message):
            res.append(i.as_dict())

        return res
