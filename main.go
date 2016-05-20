package main

import (
    "log"
    "fmt"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Cities struct {
  Status string `json:"status"`
  Items int `json:"items"`
  PageID int `json:"page_id"`
  PageLimit int `json:"page_limit"`
  Response[] City `json:"response"`
}

type City struct {
  ID string `json:"id"`
  Title string `json:"title"`
  Icon string `json:"icon"`
  Locale string `json:"locale"`
  Lat string `json:"lat"`
  Lon string `json:"lon"`
}

type Restaurants struct {
  Status      string `json:"status"`
  Items      int `json:"items"`
  PageID   int `json:"page_id"`
  PageLimit  int `json:"page_limit"`
  Response[]   Restaurant `json:"response"`
}
type Restaurant struct{
  ID         string    `json:"id"`
  LogoSquare     string `json:"logo_square"`
  CoverInner   string `json:"cover_inner"`
  Title       string `json:"city_id"`
  CityID     string  `json:"city_id"`
  Logo       string  `json:"logo"`
  Cover       string `json:"cover"`
  Description   string `json:"description"`
  Images[]     Images  `json:"images"`
  News       News `json:"news"`
  MenuFiles     MenuFiles `json:"menu_files"`
}

type MenuFiles struct {
  ID       int  `json:"id"`
  File     string  `json:"file"`
  Title     string  `json:"title"`
  Description string  `json:"description"`
  Size    string  `json:"size"`
}

type News struct {
  ID int `json:"id"`
  PlaceID int `json:"place_id"`
  Title string `json:"title"`
  Text string `json:"text"`
  Cover string `json:"cover"`
  CoverInner string `json:"cover_inner"`
}

type Images struct {
  ID string `json:"id"`
  Path string `json:"path"`
}


const token = "213609888:AAFYxyhXs5-u62eiNDFLF0exgG2l87xJQjs"

var cities = Cities {
    Response: []City{
        City{Title: "Москва", ID: "1"},
        City{Title: "Санкт-Петербург", ID: "2"},
    } }
var restaurants = []string{"Пивасик у Игоря", "Вкусная еда от тети Любы", "Перекус"}

func main() {
    bot, err := tgbotapi.NewBotAPI(token)
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
                for i, e := range cities.Response {
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%d. %s",i, e.Title))
                    bot.Send(msg)
                } 
                // msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список городов")
                // msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
                //     tgbotapi.NewInlineKeyboardRow(
                //         tgbotapi.NewInlineKeyboardButtonData("1", "eat"),
                //         tgbotapi.NewInlineKeyboardButtonData("1", "eat"),
                //         tgbotapi.NewInlineKeyboardButtonData("2", "preorder"),
                //         tgbotapi.NewInlineKeyboardButtonData("3", "menu"),
                //         tgbotapi.NewInlineKeyboardButtonData("37", "eat") ))
                
            case text == RestaurantMenuStartMenuItem:
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список меню")
                bot.Send(msg)
            case containsInCities(cities.Response, text):
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список ресторанов")
                bot.Send(msg)
            case contains(restaurants, text):
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список меню")
                bot.Send(msg)
            default:
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я Вас не понял")
                bot.Send(msg)
        }
    }
}

// WhereToEatStartMenuItem const Где поесть.
const WhereToEatStartMenuItem = "Где поесть"
// PreOrderStartMenuItem const Предзаказ.
const PreOrderStartMenuItem = "Предзаказ"
// RestaurantMenuStartMenuItem const Меню.
const RestaurantMenuStartMenuItem = "Меню"

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