package main

import (
	"time"

	"urinal/internal/transport/tg"
	"urinal/pkg/scheduler"
)

func main() {
	bot := tg.NewTgBot(time.Second * 5)
	scheduler := scheduler.NewScheduler(bot.GetFunc)
	scheduler.Run()

	bot.Run()
}
