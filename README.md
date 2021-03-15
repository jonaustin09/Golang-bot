# Money bot

Bot is designed for personal use only.

Allows tracking expenses semi-automatically using Monobank API (wehbooks) for accounts.

Export data via telegram bot or `/api/export`

Commands:

- `export` - get your data into csv file 📄
- `delete` or `d` - can delete a selected message from log 🌚

## How to run

- Prepare `.env` file with your data.
 
  ```bash
  cp env.example .env
  ```

- Install [pressly/goose](https://github.com/pressly/goose)
  
  ```bash
  go get -u github.com/pressly/goose/cmd/goose
  ```

- Run `make migrate`
- Run `make linux_build`