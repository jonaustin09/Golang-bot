import csv
from datetime import datetime
import uuid
import re

categories = {}
with open('money - log.csv', newline='') as csvfile:
    reader = csv.DictReader(csvfile)

    with open('new_money - log.csv', 'w+', newline='') as csvfile_w:
        fieldnames = ['uuid', 'date', 'item', 'value', 'category', 'telegram_user_id']

        writer = csv.DictWriter(csvfile_w, fieldnames=fieldnames)
        writer.writeheader()

        date = None
        for i, row in enumerate(reader):
            data = row['date']
            timestamp = int(datetime.strptime(data, '%d.%m.%Y').timestamp())
            row_new = row.copy()
            row_new['date'] = timestamp
            row_new['uuid'] = uuid.uuid4().hex
            row_new['telegram_user_id'] = 154701187
            print(row)
            row_new['value'] = float(re.sub('[^A-Za-z0-9\.\,]+', '', row['value']).replace(',', '.'))
            index = categories.get(row_new['category'])
            if index is None:
                categories[row_new['category']] = len(categories) + 1
                index = categories[row_new['category']]
            row_new['category'] = index
            writer.writerow(row_new)
print(categories)
