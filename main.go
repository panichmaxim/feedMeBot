package main

import (
    "fmt"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/ivahaev/russian-time"
	"encoding/json"
	"strconv"
	"strings"
	"regexp"
	"time"
)

type CallbackQueryPageData struct {
    Title string `json:"title"`
    Page int `json:"page"`
}

var Bot *tgbotapi.BotAPI

var restaurants = Restaurants{}

func main() {

    Bot, _ = tgbotapi.NewBotAPI(boToken)

    Bot.Debug = true

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    // Restaraunt initializaion
    id := "1"
    get(URLResta, id, &restaurants)

    updates, _ := Bot.GetUpdatesChan(u)

    for update := range updates {
        switch {
            case update.Message != nil:
                text := strings.ToLower(update.Message.Text)
                msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Oo")
                if (isSelectMenu(text)) {
                    menu(text, update.Message.Chat.ID)
                }
                if (isDate(text)) { // Через 2 часа  || 12.03 15:00
                    date(text, update.Message.Chat.ID)
                } 
                if (isYesOrNo(text)) {
                    responseYes(update.Message.Chat.ID)
                }
                if (isNo(text)) {
                        responseNo(update.Message.Chat.ID)
                }
                if (isThx(text)){
                    s:= "Пожалуйста, " + update.Message.From.FirstName + ". Мне нравится работать с Вами!"
                    msg2 = tgbotapi.NewMessage(update.Message.Chat.ID,s)
                    Bot.Send(msg2)
                }
                if (isOficiant(text)) {
                    msg2 = tgbotapi.NewMessage(update.Message.Chat.ID, "К Вам сейчас подойдут.")
                    Bot.Send(msg2)
                }
                if (isCardBank(text)) {
                    msg2 = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш заказ успешно забронирован. Ваш столик 51. Ждем Вас в нашем ресторане :)")
                    Bot.Send(msg2)
                }
                text = update.Message.Text
                switch {
                    case text == "/start start" || text == "/start":
                        firstName := update.Message.From.FirstName
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s!", firstName)) 
                        msg.ReplyMarkup = StartMenu()
                        Bot.Send(msg)
                    case text == WhereToEatStartMenuItem:
                        // Выберите город
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите город") 
                        msg.ReplyMarkup = CitiesMenu()
                        Bot.Send(msg)
                    case text == PreOrderStartMenuItem:
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите город") 
                        msg.ReplyMarkup = CitiesMenu()
                        Bot.Send(msg)
                    case text == RestaurantMenuStartMenuItem:
                        str := " \n" + "Блюдо: " + Menu["0"] + " \n"
                        str += "Цена: " + Price["0"] + " \n" + " \n"
                        str += "Для выбора блюд из меню используй команду \"Заказать\" и перечесление номеров блюд через пробел" + " \n"
                        str += ImageMenu["0"]
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
                        msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
                            tgbotapi.NewInlineKeyboardRow(
                                tgbotapi.NewInlineKeyboardButtonData("<", "{ \"title\":\"menu\", \"page\":0}"),
                                tgbotapi.NewInlineKeyboardButtonData("№ 0", "{ \"title\":\"menu\", \"page\":0}"),
                                tgbotapi.NewInlineKeyboardButtonData(">", "{ \"title\":\"menu\", \"page\":1}") ))
                        Bot.Send(msg)
                    case containsInCities(cities.Response, text):
                        // Отображаем рестораны с пагинацией
                        // for i, e := range restaurants.Response {
                        //     msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%d. %s",i, e.Title))
                        //     bot.Send(msg)
                        // } 
                        str := restaurants.Response[0].Title + " \n"
                        str += restaurants.Response[0].DescriptionShort + " \n"
                        str += "Адрес: " + restaurants.Response[0].ContactsAddress + " \n"
                        str += "Телефон: " + restaurants.Response[0].Telephone + " \n"
                        str += "Часы работы: " + restaurants.Response[0].Work + " \n" + " \n"
                        str +=  URLImages + restaurants.Response[0].Cover
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
                        msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
                            tgbotapi.NewInlineKeyboardRow(
                                tgbotapi.NewInlineKeyboardButtonData("<", "{ \"title\":\"restaraunt\", \"page\":0}"),
                                tgbotapi.NewInlineKeyboardButtonData("Выбрать", "{ \"title\":\"restaraunt\", \"page\":0}"),
                                tgbotapi.NewInlineKeyboardButtonData(">", "{ \"title\":\"restaraunt\", \"page\":1}") ))
                        Bot.Send(msg)
                    case containsInRestaurants(restaurants.Response, text):
                        // Переделать на выбор меню
                        
                }
            case update.CallbackQuery != nil:
                fmt.Println("CallbackQuery")
                var callBack CallbackQueryPageData
                fmt.Println(update.CallbackQuery.Data)
                json.Unmarshal([]byte(update.CallbackQuery.Data), &callBack)
                fmt.Println(callBack.Page)
                switch {
                    case callBack.Title == "restaraunt":
                        prevPage := 0
                        nextPage := 0
                        if callBack.Page == 0 {
                            prevPage = callBack.Page
                            nextPage = callBack.Page+1
                        } 
                        if callBack.Page != 0 && callBack.Page != len(restaurants.Response) {
                            prevPage = callBack.Page-1
                            nextPage = callBack.Page+1
                        } 
                        if callBack.Page == len(restaurants.Response) {
                            prevPage = callBack.Page-1
                            nextPage = callBack.Page+1
                        } 
                        str := restaurants.Response[callBack.Page].Title + " \n"
                        str += restaurants.Response[callBack.Page].DescriptionShort + " \n"
                        str += "Адрес: " + restaurants.Response[callBack.Page].ContactsAddress + " \n"
                        str += "Телефон: " + restaurants.Response[callBack.Page].Telephone + " \n"
                        str += "Часы работы: " + restaurants.Response[callBack.Page].Work + " \n" + " \n"
                        str += URLImages + restaurants.Response[callBack.Page].Cover 
                        msg := tgbotapi.NewEditMessageText(int64(update.CallbackQuery.From.ID), update.CallbackQuery.Message.MessageID, str)
                        kb := tgbotapi.NewInlineKeyboardMarkup(
                            tgbotapi.NewInlineKeyboardRow(
                                tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("{ \"title\":\"restaraunt\", \"page\":%d}", prevPage)),
                                tgbotapi.NewInlineKeyboardButtonData("Выбрать", fmt.Sprintf("{ \"title\":\"menu_choose\", \"page\":%d}", callBack.Page)),
                                tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("{ \"title\":\"restaraunt\", \"page\":%d}", nextPage)) ))
                        msg.ReplyMarkup = &kb
                        Bot.Send(msg)
                    case callBack.Title == "menu":
                        prevPage := 0
                        nextPage := 0
                        if callBack.Page == 0 {
                            prevPage = callBack.Page
                            nextPage = callBack.Page+1
                        } 
                        if callBack.Page != 0 && callBack.Page != len(restaurants.Response) {
                            prevPage = callBack.Page-1
                            nextPage = callBack.Page+1
                        } 
                        if callBack.Page == len(restaurants.Response) {
                            prevPage = callBack.Page-1
                            nextPage = callBack.Page+1
                        } 
                        str := " \n" + "Блюдо: " + Menu[strconv.FormatInt(int64(callBack.Page), 10)] + " \n"
                        str += "Цена: " + Price[strconv.FormatInt(int64(callBack.Page), 10)] + " \n" + " \n"
                        str += "Для выбора блюд из меню используй команду \"Заказать\" и перечесление номеров блюд через пробел" + " \n"
                        str += ImageMenu[strconv.FormatInt(int64(callBack.Page), 10)]
                        msg := tgbotapi.NewEditMessageText(int64(update.CallbackQuery.From.ID), update.CallbackQuery.Message.MessageID, str)
                        kb := tgbotapi.NewInlineKeyboardMarkup(
                            tgbotapi.NewInlineKeyboardRow(
                                tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("{ \"title\":\"menu\", \"page\":%d}", prevPage)),
                                tgbotapi.NewInlineKeyboardButtonData("№ " + strconv.Itoa(callBack.Page), fmt.Sprintf("{ \"title\":\"menu\", \"page\":%d}", callBack.Page)),
                                tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("{ \"title\":\"menu\", \"page\":%d}", nextPage)) ))
                        msg.ReplyMarkup = &kb
                        Bot.Send(msg)
                    case callBack.Title == "menu_choose":
                        str := " \n" + "Блюдо: " + Menu["0"] + " \n"
                        str += "Цена: " + Price["0"] + " \n" + " \n" 
                        str += "Для выбора блюд из меню используй команду \"Заказать\" и перечесление номеров блюд через пробел" + " \n"
                        str += ImageMenu["0"]
                        msg := tgbotapi.NewMessage(int64(update.CallbackQuery.From.ID), str)
                        msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
                            tgbotapi.NewInlineKeyboardRow(
                                tgbotapi.NewInlineKeyboardButtonData("<", "{ \"title\":\"menu\", \"page\":0}"),
                                tgbotapi.NewInlineKeyboardButtonData("№ 0", "{ \"title\":\"menu\", \"page\":0}"),
                                tgbotapi.NewInlineKeyboardButtonData(">", "{ \"title\":\"menu\", \"page\":1}") ))
                        Bot.Send(msg)
                }
        }
    }
}

// StartMenu func return start menu of chatbot.
func StartMenu() tgbotapi.ReplyKeyboardMarkup {  
    return tgbotapi.NewReplyKeyboard(
            []tgbotapi.KeyboardButton {
                tgbotapi.NewKeyboardButton(WhereToEatStartMenuItem),
            }, 
            []tgbotapi.KeyboardButton {
                tgbotapi.NewKeyboardButton(PreOrderStartMenuItem),
            },
            []tgbotapi.KeyboardButton {
                tgbotapi.NewKeyboardButton(RestaurantMenuStartMenuItem),
            } )
}

// CitiesMenu func return cities menu of chatbot.
func CitiesMenu() tgbotapi.ReplyKeyboardMarkup { 
    return NewReplyHideKeyboard(
            []tgbotapi.KeyboardButton {
                tgbotapi.NewKeyboardButton("Москва"),
            }, 
            []tgbotapi.KeyboardButton {
                tgbotapi.NewKeyboardButton("Санкт-Петербург"),
            })
}

func containsInCities(cities []City, name string) bool {
    for _, city := range cities {
        if city.Title == name {
            return true
        }
    }
    return false
}

func containsInRestaurants(restaurants []Restaurant, name string) bool {
    for _, restaurant := range restaurants {
        if restaurant.Title == name {
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

func NewReplyHideKeyboard(rows ...[]tgbotapi.KeyboardButton) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	keyboard = append(keyboard, rows...)

	return tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
        OneTimeKeyboard: true,
	}
}

// Обработка текстовых команд

func isOficiant(text string) bool {
    return strings.Contains(text,"официан")
}
func isThx(text string) bool{
    return strings.Contains(text,"=)")
}
func isNo(text string) bool{
    return strings.Contains(text, "нет") && len(text)<8 // 2 bytes one symbol
}
func isCardBank(number string ) bool{
    str := strings.Replace(number, " " ,"",-1)
    fmt.Println(str)
    match, _ := regexp.MatchString("([a-z]+)", str)
    fmt.Println(match)
    if  !match  && len(str)>10 && strings.Contains(str,"1234") {
          return true
    } else{
           return false
    }
}
func responseNo(id int64) { 
    str := " \n" + "Оо А что не так? Попробуйте еще раз с выбора блюд.\n(Пример: Заказ [номер])" + " \n" + " \n"
    str += "Блюдо: " + Menu["0"] + " \n"
    str += "Цена: " + Price["0"] + " \n"
    str += ImageMenu["0"]
    msg := tgbotapi.NewMessage(id, str)
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
    tgbotapi.NewInlineKeyboardButtonData("<", "{ \"title\":\"menu\", \"page\":0}"),
    tgbotapi.NewInlineKeyboardButtonData("№ 0", "{ \"title\":\"menu\", \"page\":0}"),
    tgbotapi.NewInlineKeyboardButtonData(">", "{ \"title\":\"menu\", \"page\":1}") ))
    Bot.Send(msg)
}
func responseYes(id int64) {
    msg := tgbotapi.NewMessage(id, "Введите номер банковской карты. После мы сможем забронировать Вам столик ^^")
    Bot.Send(msg)
}
func isYesOrNo(str string) bool{
    strings.Replace(str," ","",-1)
    fmt.Println(str)
    return strings.Contains(str,"yes") || (strings.Contains(str, "да") && len(str) >3)
}
func date(date string, id int64 ) {
    if (strings.Contains(date,"через")) {
        across(date, id)
    } else {
        dateTime(date, id)
    }
}
func across(dateTime string, id int64) {
    re := regexp.MustCompile("[0-9]+")
    time := re.FindAllString(dateTime, -1)
    var str string
    if (strings.Contains(dateTime,"минут")) {
        str = "Через " + time[0] + " минут мы забронируем Вам столик. Все верно?"
    } 
    if (strings.Contains(dateTime, "час")) {
        str = "Через " + time[0] + " часов мы забронируем Вам столик. Все точно?"
    }
    msg := tgbotapi.NewMessage(id, str)
    Bot.Send(msg)
}
func dateTime(date string, id int64 ) {
    t := rtime.Now()
    // Or if you are using time.Time object:
    standardTime := time.Now()
    t = rtime.Time(standardTime)
    msg := tgbotapi.NewMessage(id,"Ваш заказ будет доступен через "+ t.TimeString() + ". Все верно?")
    Bot.Send(msg)
}

func isDate(date string) bool {
    return strings.Contains(date,"через") || (strings.Contains(date,".") && strings.Contains(date,":")) 
} 

func isSelectMenu(str string ) bool {
    return strings.Contains(str,"заказ") && strings.Contains(str, " ")
}

func menu(data string, id int64) {
    data =  strings.Replace(data,"заказ","",-1)
    var result[] string
    if (strings.Contains(data,","))  {
        result = strings.Split(data, ",");
    } else {
        result = strings.Fields(data)
    }
    responseMsg :="Ваш итоговый заказ составил:\n\n"
    for i := range result {
        num,_:=strconv.Atoi(result[i]);
        responseMsg += Menu[strconv.FormatInt(int64(num+1),10)] + " "+ Price[strconv.FormatInt(int64(num+1),10)] + " рублей \n"
    }
    responseMsg +="\nКогда хотите нас посетить? \n (Пример:\"Через 2 часа\" или \" 30.04 15:00\")"
    msg := tgbotapi.NewMessage(id, responseMsg)
    Bot.Send(msg)
}




