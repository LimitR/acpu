package main

import (
	"time"

	"urinal/internal/transport/tg"
	"urinal/pkg/scheduler"
)

func main() {
	bot := tg.NewTgBot(time.Second * 5)
	scheduler := scheduler.NewScheduler(time.Millisecond*100, bot.GetFunc, func(err error) {})
	scheduler.Run()

	bot.Run()
}
