package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/ivahaev/russian-time"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CallbackQueryPageData struct {
	Title string `json:"title"`
	Page  int    `json:"page"`
}

var Bot *tgbotapi.BotAPI

var restaurants = Restaurants{}

var spellcheck *ToySpellcheck

func main() {
    
    spellcheck = &ToySpellcheck{}
	spellcheck.Train("./bigRus.txt")

	Bot, _ = tgbotapi.NewBotAPI(botToken)

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
            switch {
                case isSelectMenu(text):
                    menu(text, update.Message.Chat.ID)
                    break
                case isDate(text): // Через 2 часа  || 12.03 15:00
                    date(text, update.Message.Chat.ID)
                    break
                case isYesOrNo(text):
                    responseYes(update.Message.Chat.ID)
                    break
                case isNo(text):
                    responseNo(update.Message.Chat.ID)
                    break
                case isThx(text):
                    s := "Пожалуйста, " + update.Message.From.FirstName + ". Мне нравится работать с Вами!"
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, s)
                    Bot.Send(msg)
                    break
                case isOficiant(text):
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "К Вам сейчас подойдут.")
                    Bot.Send(msg)
                    break
                case isCardBank(text):
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш заказ успешно забронирован. Ваш столик 51. Ждем Вас в нашем ресторане :)")
                    Bot.Send(msg)
                    break
                case text == "/start start" || text == "/start":
                    firstName := update.Message.From.FirstName
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s!", firstName))
                    msg.ReplyMarkup = StartMenu()
                    Bot.Send(msg)
                    break
                case text == strings.ToLower(PreOrderStartMenuItem):
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите город")
                    msg.ReplyMarkup = CitiesMenu()
                    Bot.Send(msg)
                    break
                case containsInCities(cities.Response, text):
                    str, kb := RestarauntsPager(0)
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
                    msg.ReplyMarkup = kb
                    Bot.Send(msg)
                    break
			}
		case update.CallbackQuery != nil:
			var callBack CallbackQueryPageData
			json.Unmarshal([]byte(update.CallbackQuery.Data), &callBack)
			switch {
                case callBack.Title == RestarauntCallbackTitle:
                    str, kb := RestarauntsPager(callBack.Page)
                    msg := tgbotapi.NewEditMessageText(int64(update.CallbackQuery.From.ID), update.CallbackQuery.Message.MessageID, str)
                    msg.ReplyMarkup = &kb
                    Bot.Send(msg)
                    break
                case callBack.Title == MenuCallbackTitle:
                    str, kb := MenuPager(callBack.Page)
                    msg := tgbotapi.NewEditMessageText(int64(update.CallbackQuery.From.ID), update.CallbackQuery.Message.MessageID, str)
                    msg.ReplyMarkup = &kb
                    Bot.Send(msg)
                    break
                case callBack.Title == MenuChooseCallbackTitle:
                    str, kb := MenuPager(0)
                    msg := tgbotapi.NewMessage(int64(update.CallbackQuery.From.ID), str)
                    msg.ReplyMarkup = &kb
                    Bot.Send(msg)
                    break
			}
		}
	}
}

// StartMenu func returns start menu of chatbot.
func StartMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(PreOrderStartMenuItem),
		},)
}

// CitiesMenu func returns cities menu of chatbot.
func CitiesMenu() tgbotapi.ReplyKeyboardMarkup {
	return NewReplyHideKeyboard(
		[]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Москва"),
		},
		[]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Санкт-Петербург"),
		})
}

// MenuPager func returns menu with view pager.
func MenuPager(page int) (string, tgbotapi.InlineKeyboardMarkup) {
    prevPage := 0
	nextPage := 0
	if page == 0 {
		prevPage = page
		nextPage = page + 1
	}
	if page != 0 && page != len(restaurants.Response) {
        prevPage = page - 1
        nextPage = page + 1
	}
	if page == len(restaurants.Response) {
		prevPage = page - 1
		nextPage = page + 1
	}
	str := " \n" + "Блюдо: " + Menu[strconv.FormatInt(int64(page), 10)] + " \n"
	str += "Цена: " + Price[strconv.FormatInt(int64(page), 10)] + " \n" + " \n"
	str += "Для выбора блюд из меню используй команду \"Заказать\" и перечесление номеров блюд через пробел" + " \n"
	str += ImageMenu[strconv.FormatInt(int64(page), 10)]
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("{ \"title\":\"" + MenuCallbackTitle + "\", \"page\":%d}", prevPage)),
			tgbotapi.NewInlineKeyboardButtonData("№ "+strconv.Itoa(page+1), fmt.Sprintf("{ \"title\":\"" + MenuCallbackTitle + "\", \"page\":%d}", page)),
			tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("{ \"title\":\"" + MenuCallbackTitle + "\", \"page\":%d}", nextPage))))
    return str, kb
}

// RestarauntsPager func returns restaraunts with view pager.
func RestarauntsPager(page int) (string, tgbotapi.InlineKeyboardMarkup) {
    prevPage := 0
    nextPage := 0
    if page == 0 {
		prevPage = page
		nextPage = page + 1
	}
	if page != 0 && page != len(restaurants.Response) {
		prevPage = page - 1
		nextPage = page + 1
	}
	if page == len(restaurants.Response) {
		prevPage = page - 1
		nextPage = page + 1
	}
	str := restaurants.Response[page].Title + " \n"
	str += restaurants.Response[page].DescriptionShort + " \n"
	str += "Адрес: " + restaurants.Response[page].ContactsAddress + " \n"
	str += "Телефон: " + restaurants.Response[page].Telephone + " \n"
	str += "Часы работы: " + restaurants.Response[page].Work + " \n" + " \n"
	str += URLImages + restaurants.Response[page].Cover
	kb := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("{ \"title\":\"" + RestarauntCallbackTitle + "\", \"page\":%d}", prevPage)),
            tgbotapi.NewInlineKeyboardButtonData("Выбрать", fmt.Sprintf("{ \"title\":\"" + MenuChooseCallbackTitle + "\", \"page\":%d}", page)),
            tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("{ \"title\":\"" + RestarauntCallbackTitle + "\", \"page\":%d}", nextPage))))
	return str, kb
}

func containsInCities(cities []City, name string) bool {
	for _, city := range cities {
		if strings.ToLower(city.Title) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

// В теории это функция нужна если у нас много менюшек, на данный момент только 1 меню для всех ресторанов
func containsInRestaurants(restaurants []Restaurant, name string) bool {
	for _, restaurant := range restaurants {
		if strings.ToLower(restaurant.Title) == strings.ToLower(name) {
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
		ResizeKeyboard:  true,
		Keyboard:        keyboard,
		OneTimeKeyboard: true,
	}
}

// Обработка текстовых команд

func isOficiant(text string) bool {
	return strings.Contains(spellcheck.Correct(text), "официант")
}
func isThx(text string) bool {
	return strings.Contains(spellcheck.Correct(text), "спасибо")
}
func isNo(text string) bool {
	return strings.Contains(text, "нет") && len(text) < 8 // 2 bytes one symbol
}
func isCardBank(number string) bool {
	str := strings.Replace(number, " ", "", -1)
	fmt.Println(str)
	match, _ := regexp.MatchString("([a-z]+)", str)
	fmt.Println(match)
	if !match && len(str) > 10 && strings.Contains(str, "1234") {
		return true
	} else {
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
			tgbotapi.NewInlineKeyboardButtonData("<", "{ \"title\":\"" + MenuCallbackTitle + "\", \"page\":0}"),
			tgbotapi.NewInlineKeyboardButtonData("№ 0", "{ \"title\":\"" + MenuCallbackTitle + "\", \"page\":0}"),
			tgbotapi.NewInlineKeyboardButtonData(">", "{ \"title\":\"" + MenuCallbackTitle + "\", \"page\":1}")))
	Bot.Send(msg)
}
func responseYes(id int64) {
	msg := tgbotapi.NewMessage(id, "Введите номер банковской карты. После мы сможем забронировать Вам столик ^^")
	Bot.Send(msg)
}
func isYesOrNo(str string) bool {
	strings.Replace(str, " ", "", -1)
	fmt.Println(str)
	return strings.Contains(str, "yes") || (strings.Contains(str, "да") && len(str) > 3)
}
func date(date string, id int64) {
	if strings.Contains(date, "через") {
		across(date, id)
	} else {
		dateTime(date, id)
	}
}
func across(dateTime string, id int64) {
	re := regexp.MustCompile("[0-9]+")
	time := re.FindAllString(dateTime, -1)
	var str string
	if strings.Contains(dateTime, "минут") {
		str = "Через " + time[0] + " минут мы забронируем Вам столик. Все верно?"
	}
	if strings.Contains(dateTime, "час") {
		str = "Через " + time[0] + " часов мы забронируем Вам столик. Все точно?"
	}
	msg := tgbotapi.NewMessage(id, str)
	Bot.Send(msg)
}
func dateTime(date string, id int64) {
	t := rtime.Now()
	// Or if you are using time.Time object:
	standardTime := time.Now()
	t = rtime.Time(standardTime)
	msg := tgbotapi.NewMessage(id, "Ваш заказ будет доступен через "+t.TimeString()+". Все верно?")
	Bot.Send(msg)
}

func isDate(date string) bool {
	return strings.Contains(date, "через") || (strings.Contains(date, ".") && strings.Contains(date, ":"))
}

func isSelectMenu(str string) bool {
	return strings.Contains(str, "заказ") && strings.Contains(str, " ")
}

func menu(data string, id int64) {
	data = strings.Replace(data, "заказать", "", -1)
	var result []string
	if strings.Contains(data, ",") {
		result = strings.Split(data, ",")
	} else {
		result = strings.Fields(data)
	}
	responseMsg := "Ваш итоговый заказ составил:\n\n"
	for i := range result {
		num, _ := strconv.Atoi(result[i])
		responseMsg += Menu[strconv.FormatInt(int64(num-1), 10)] + " " + Price[strconv.FormatInt(int64(num-1), 10)] + " рублей \n"
	}
	responseMsg += "\nКогда хотите нас посетить? \n (Пример:\"Через 2 часа\" или \" 30.04 15:00\")"
	msg := tgbotapi.NewMessage(id, responseMsg)
	Bot.Send(msg)
}
