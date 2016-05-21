package main

// Api examples

 // Cities example of possible cities array
var cities = Cities {
    Response: []City{
        City{Title: "Москва", ID: "61"},
        City{Title: "Санкт-Петербург", ID: "2"},
    } }
 
// Menu is an example of possible menu
var Menu = map[string]string {
        "0":  "Брускетта с томатами",
        "1":  "Брускетта с лососем",
        "2":   "Брускетта с крабом",
        "3":   "Ассорти колбас",
        "4":   "Ассорти европейских сыров",
        "5":   "Антипасти",
        "6":   "Оливье по-домашнему",
        "7":   "Овощной салат с грядки",
        "8":   "Тайский салат с цыпленком",
        "9":   "Салат с бакинскими овощами и яйцом-пашот",
        "10":   "Салат с треской",
        "12":   "Моцарелла с томатами и песто",
        "13":   "Теплый салат с курицей и копченым сулугуни",
        "14":   "Салат с бакинскими томатами с тархуном или щавелем",
        "15":   "Сугудай",
        "16":   "Цезарь с курицей",
        "17":   "Руккола с креветками", }

// Price is an example of possible menu prices
var Price = map[string]string {
        "0": "330",
        "1": "120",
        "2": "140",
        "3": "90",
        "4": "190",
        "5": "50", 
        "6": "290",
        "7": "390",
        "8": "490",
        "9": "590",
        "10": "690",
        "12": "790",
        "13": "900",
        "14": "200",
        "15": "300",
        "16": "440",
        "17": "230", }
    
// ImageMenu is an example of possible menu images
var ImageMenu = map[string]string {
        "0":"http://amam.ru/recepts_img/327/brusketta_s_tom-shag_6_amam.ru.jpg",
        "1":   "http://kakgotovit.com/wp-content/uploads/2014/10/Brusketta-s-lososem-i-avokado.jpg",
        "2":   "http://file.lavkalavka.com/recipe//1920/brusketta-s-kamchatskim-krabom.jpg",
        "3":   "http://leow.ru/wp-content/uploads/2012/01/bankoboev.ru_kolbasy_i_myaso.jpg",
        "4":   "http://wallpapers1.ru/fud/data/eda_0_000.jpg",
        "5":   "http://wallpapers1.ru/fud/data/eda_0_000.jpg",
        "6":   "http://wallpapers1.ru/fud/data/eda_0_000.jpg",
        "7":   "http://wallpapers1.ru/fud/data/eda_0_000.jpg",
        "8":   "http://wallpapers1.ru/fud/data/eda_0_000.jpg",
        "9":   "http://wallpapers1.ru/fud/data/eda_0_000.jpg",
        "10":   "http://wallpapers1.ru/fud/data/eda_0_000.jpg", }


// Bot constants

// PreOrderStartMenuItem constant Предзаказ.
const PreOrderStartMenuItem = "Предзаказ"