package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Alexsilvacodes/EsiosBot/esios"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func sendMessage(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func priceColor(price float64) string {
	if price < 0.1 {
		return "🟢"
	} else if price >= 0.1 && price <= 0.15 {
		return "🟠"
	} else {
		return "🔴"
	}
}

func main() {
	err := godotenv.Load(".env")
	token := os.Getenv("TELEGRAM_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err) // You should add better error handling than this!
	}

	log.Printf("Started '%s'", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				bot.Request(tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.MessageID,
				})

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = "⚡️🇪🇸 Bienvenid@ al bot de precios oficiales de Red Eléctrica Nacional. Plataforma fuente: " +
					"ESIOS.\n Para recibir la lista de precios por hora ejecuta el comando /precios_hora" +
					"\n Para recibir los picos de precio ejecuta el comando /precios_pico"

				sendMessage(msg, bot)
			case "/precios_pico":
				bot.Request(tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.MessageID,
				})

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				prices := esios.GetPrice()
				priceStr := "💰 Precios " + prices.PVPC[0].Dia + "\n\n" +
					"_Consultado: " + time.Now().Format("15:04:05") + "_ \n\n"
				cheapPrice := prices.PVPC[0]
				cheapPricePCB := 0.0
				cheapPriceStr := strings.Replace(prices.PVPC[0].PCB, ",", ".", 1)
				if s, err := strconv.ParseFloat(cheapPriceStr, 32); err == nil {
					cheapPricePCB = s / 1000
				}
				expensivePrice := prices.PVPC[0]
				expensivePricePCB := 0.0
				expensivePriceStr := strings.Replace(prices.PVPC[0].PCB, ",", ".", 1)
				if s, err := strconv.ParseFloat(expensivePriceStr, 32); err == nil {
					expensivePricePCB = s / 1000
				}
				for _, price := range prices.PVPC {
					if cheapPrice.PCB > price.PCB {
						cheapPrice = price
						if s, err := strconv.ParseFloat(strings.Replace(price.PCB, ",", ".", 1), 32); err == nil {
							cheapPricePCB = s / 1000
						}
					}

					if expensivePrice.PCB < price.PCB {
						expensivePrice = price
						if s, err := strconv.ParseFloat(strings.Replace(price.PCB, ",", ".", 1), 32); err == nil {
							expensivePricePCB = s / 1000
						}
					}
				}
				priceStr += "🟢🕐 " + cheapPrice.Hora + ": " + fmt.Sprintf("%f", cheapPricePCB) + " € / kWh\n"
				priceStr += "🔴🕐 " + expensivePrice.Hora + ": " + fmt.Sprintf("%f", expensivePricePCB) + " € / kWh\n"
				msg.Text = priceStr
				msg.ParseMode = tgbotapi.ModeMarkdown

				sendMessage(msg, bot)
			case "/precios_hora":
				bot.Request(tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.MessageID,
				})

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				prices := esios.GetPrice()
				priceStr := "💰 Precios Hora " + prices.PVPC[0].Dia + "\n\n" +
					"_Consultado: " + time.Now().Format("15:04:05") + "_ \n\n"
				for _, price := range prices.PVPC {
					pcbStr := strings.Replace(price.PCB, ",", ".", 1)
					pcb := 0.0
					if s, err := strconv.ParseFloat(pcbStr, 32); err == nil {
						pcb = s / 1000
					}
					priceStr += "🕐 " + price.Hora + " " + priceColor(pcb) + ": " + fmt.Sprintf("%f", pcb) + " € / kWh\n"
				}
				msg.Text = priceStr
				msg.ParseMode = tgbotapi.ModeMarkdown

				sendMessage(msg, bot)
			}
		}
	}
}
