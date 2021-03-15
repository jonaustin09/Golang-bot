package main

import (
	"github.com/dobrovolsky/money_bot/moneybot"
)

func main() {
	app := moneybot.BuildApp()
	app.Start()
}
