package tg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"urinal/internal/models"
	"urinal/pkg/cache"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1", "0"),
		tgbotapi.NewInlineKeyboardButtonData("2", "1"),
		tgbotapi.NewInlineKeyboardButtonData("3", "2"),
		tgbotapi.NewInlineKeyboardButtonData("4", "3"),
		tgbotapi.NewInlineKeyboardButtonData("5", "4"),
		tgbotapi.NewInlineKeyboardButtonData("6", "5"),
		tgbotapi.NewInlineKeyboardButtonData("7", "6"),
		tgbotapi.NewInlineKeyboardButtonData("8", "7"),
		tgbotapi.NewInlineKeyboardButtonData("9", "8"),
	),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Optimal place", "Optimal place")),
)

type TgBot struct {
	bot         *tgbotapi.BotAPI
	mapModel    map[int64]*models.Toilet
	cache       *cache.Cache[int64, string]
	cacheNumber *cache.Cache[int64, []int]
}

func NewTgBot(timeoutCache time.Duration) *TgBot {
	token, ok := os.LookupEnv("TG_TOKEN")
	if !ok {
		log.Fatal("Not found tg token")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	cacheInstans := cache.NewCache[int64, string](timeoutCache)
	cacheNumber := cache.NewCache[int64, []int](timeoutCache * 2)
	return &TgBot{
		bot:         bot,
		mapModel:    make(map[int64]*models.Toilet, 100),
		cache:       cacheInstans,
		cacheNumber: cacheNumber,
	}
}

func (b *TgBot) GetFunc() error {
	return b.cache.CheckTime()
}

func (b *TgBot) Run() {
	u := tgbotapi.NewUpdate(0)
	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			cache, ok := b.cache.Get(update.Message.Chat.ID)
			if ok {
				if cache == "/help" {
					v, err := strconv.Atoi(update.Message.Text)
					if err != nil {
						b.sendMessage(update.Message.From.ID, "You must enter a number")
					}
					mapModel, ok := b.mapModel[update.Message.Chat.ID]
					if ok {
						mapModel.AddUrinal(v)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Specify the seats that are occupied")
						msg.ReplyMarkup = numericKeyboard
						if _, err = b.bot.Send(msg); err != nil {
							panic(err)
						}
						b.mapModel[update.Message.Chat.ID] = mapModel
						continue
					} else {
						model := models.NewToilet()
						model.AddUrinal(v)
						b.mapModel[update.Message.Chat.ID] = model
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Specify the seats that are occupied")
						msg.ReplyMarkup = numericKeyboard
						if _, err = b.bot.Send(msg); err != nil {
							panic(err)
						}
						continue
					}
				}
			} else {
				b.cache.AddElement(update.Message.Chat.ID, "/help")
				b.sendMessage(update.Message.From.ID, "Enter a quantity urinal")
				continue
			}
		} else if update.CallbackQuery != nil {
			model, ok := b.mapModel[update.CallbackQuery.Message.Chat.ID]
			if !ok {
				b.sendMessage(update.CallbackQuery.Message.Chat.ID, "Try again /help")
			}
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.bot.Request(callback); err != nil {
				panic(err)
			}

			cache, ok := b.cacheNumber.Get(update.CallbackQuery.Message.Chat.ID)
			if !ok {
				array := make([]int, 0, 100)
				b.cacheNumber.AddElement(update.CallbackQuery.Message.Chat.ID, array)
			}

			if update.CallbackQuery.Data == "Optimal place" {
				fmt.Println("Result")
				array, _ := b.cacheNumber.Get(update.CallbackQuery.Message.Chat.ID)
				for _, v := range array {
					model.TakePlace(v)
				}
				result, err := model.GetOptimalPlace()
				if err != nil {
					b.sendMessage(update.CallbackQuery.Message.Chat.ID, err.Error())
					continue
				}
				b.sendMessage(update.CallbackQuery.Message.Chat.ID, strconv.Itoa(result+1))
				delete(b.mapModel, update.CallbackQuery.Message.Chat.ID)
				continue
			}
			v, err := strconv.Atoi(callback.Text)
			if err != nil {
				b.sendMessage(update.CallbackQuery.Message.Chat.ID, "You must enter a number")
			}
			cache = append(cache, v)
			b.cacheNumber.AddElement(update.CallbackQuery.Message.Chat.ID, cache)
		}
	}
}

func (b *TgBot) sendMessage(id int64, text string) {
	repl := tgbotapi.NewMessage(id, text)
	b.bot.Send(repl)
}

func (b *TgBot) SetUrinals(userId int64, count int) {
	model := models.NewToilet()
	model.AddUrinal(count)
	b.mapModel[userId] = model
}

func (b *TgBot) SetPlace(userId int64, arrayPlace []int) {
	model, ok := b.mapModel[userId]
	if !ok {
		return
	}
	for _, v := range arrayPlace {
		model.TakePlace(v)
	}
	b.mapModel[userId] = model
}

func (b *TgBot) GetResult(userId int64) (int, error) {
	model, ok := b.mapModel[userId]
	if !ok {
		return 0, errors.New("Not found model")
	}
	delete(b.mapModel, userId)
	return model.GetOptimalPlace()
}
