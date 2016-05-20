package main

import (
    "log"
    "github.com/Syfaro/telegram-bot-api"
)

const token = "213609888:AAFYxyhXs5-u62eiNDFLF0exgG2l87xJQjs"

func main() {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
        msg.ReplyToMessageID = update.Message.MessageID

        bot.Send(msg)
    }
}