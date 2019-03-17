### What if you want to try:
https://t.me/your_money_helper_bot
Commands:
- `stat_all_by_month` - get statistic based aggerated by month🤓
- `stat_by_category` - get statistic of your spending based on categories🤓
- `export` - get your data into csv file 📄
- `delete` - can delete selected message from log 🌚

### How to deploy?
1) run `docker-compose -f docker-compose.yaml up`

#### How to generate proto:
run command in container: 
```bash
protoc -I=stats/proto stats/proto/stats.proto --go_out=plugins=grpc:stats
```

#### How to create new migrations:
Read https://github.com/pressly/goose#usage
run command in container: 
```bash
cd migrations
goose sqlite3 ../db.sqlite3 create <name> sql
```
