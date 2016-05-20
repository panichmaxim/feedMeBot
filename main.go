package main

import (
    "log"
    "fmt"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
    bot, err := tgbotapi.NewBotAPI(Token)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        text := update.Message.Text
        switch {
            case text == "/start start":
                firstName := update.Message.From.FirstName
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s!", firstName)) 
                msg.ReplyMarkup = StartMenu()
                bot.Send(msg)
            case text == WhereToEatStartMenuItem:
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "В разработке")
                bot.Send(msg)
            case text == PreOrderStartMenuItem:
                for i, e := range Cities.Response {
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%d. %s",i, e.Title))
                    bot.Send(msg)
                }                 
            case text == RestaurantMenuStartMenuItem:
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список меню")
                bot.Send(msg)
            case containsInCities(Cities.Response, text):
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список ресторанов")
                bot.Send(msg)
            case contains(Restaurants, text):
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список меню")
                bot.Send(msg)
            default:
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я Вас не понял")
                bot.Send(msg)
        }
    }
}


// StartMenu func return start menu of application.
func StartMenu() tgbotapi.ReplyKeyboardMarkup {  
    return tgbotapi.NewReplyKeyboard(
            []tgbotapi.KeyboardButton {
                tgbotapi.NewKeyboardButtonLocation(WhereToEatStartMenuItem),
            }, 
            []tgbotapi.KeyboardButton {
                tgbotapi.NewKeyboardButton(PreOrderStartMenuItem),
            },
            []tgbotapi.KeyboardButton {
                tgbotapi.NewKeyboardButton(RestaurantMenuStartMenuItem),
            } )
}

func containsInCities(cities []City, name string) bool {
    for _, city := range cities {
        if city.Title == name {
            return true
        }
    }
    return false
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}