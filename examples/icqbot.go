package main

import (
    "../src/irbis"
    "context"
    "encoding/json"
    "fmt"
    "github.com/mail-ru-im/bot-golang"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
)

type BotConfig struct {
	Token string
	Host string
	Port int
	Database string
	User string
	Password string
}

func readConfig() BotConfig {
	bytes, err := ioutil.ReadFile("bot-config.json")
	if err != nil {
		log.Fatal(err)
	}
	var result BotConfig
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func getConnection() *irbis.IrbisConnection {
	config := readConfig()
	result := irbis.IrbisConnection{}
	result.Host = config.Host
	result.Port = config.Port
	result.Database = config.Database
	result.Username = config.User
	result.Password = config.Password
	result.Workstation = "C"

	return &result
}

type Announce struct {
	Id string
	Name string
	Preview_Text string
	Detail_Picture string
	Property_Date_S1_Value string
	Property_Category_Value string
	Detail_Page_Url string
}

func getUrlText (url string) (content []byte, err error) {
	response, err2 := http.Get(url)
	if err2 != nil {
		err = err2
		log.Printf("Error getting url #{err2}")
		return
	}
	defer func () { _ = response.Body.Close() } ()
	content, err = ioutil.ReadAll(response.Body)
	return
}

func getAnnounces() (result []Announce, err error) {
	url := "http://irklib.ru/api"
	var content []byte
	content, err = getUrlText(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &result)
	return
}

func addButtons (message *botgolang.Message) {
	buttons := [][] botgolang.Button {
		{
			botgolang.NewCallbackButton("Анонсы", "Анонсы"),
			botgolang.NewCallbackButton("Контакты", "Контакты"),
		},
		{
			botgolang.NewCallbackButton("Режим работы", "Режим работы"),
			botgolang.NewCallbackButton("Помощь", "Помощь"),
		},
	}
	message.InlineKeyboard = buttons
}

func choose(first string, second string) string {
	if first == "" {
		return second
	}
	return first
}

func doAnnounces (bot *botgolang.Bot, chatId string, userId string) {
	announces, err := getAnnounces()

	if len(announces) == 0 {
		message := bot.NewTextMessage(chatId,
			"К сожалению, анонсов пока нет")
		_ = message.Send()
		return
	}


	if err == nil {
		for i := range announces {
			announce := bot.NewTextMessage(choose(chatId, userId),
				announces[i].Preview_Text)
			if i == len(announces) - 1 {
				addButtons(announce)
			}
			_ = announce.Send()
		}
	}
}

func doStart (bot *botgolang.Bot, chatId string, userId string) {
	text := `Я сижу в подвале среди миллионов книг.
Могу найти книжку, могу не найти 😁`
	start := bot.NewTextMessage(choose(chatId, userId), text)
	addButtons(start)
	_ = start.Send()
}

func doContacts (bot *botgolang.Bot, chatId string, userId string) {
	text := `Почтовый адрес: 664033, г. Иркутск, ул. Лермонтова, 253
Электронная почта: library@irklib.ru
Многоканальный телефон: (3952) 48-66-80
Добавочный номер приемной: 705`
	contacts := bot.NewTextMessage(choose(chatId, userId), text)
	addButtons(contacts)
	_ = contacts.Send()
}

func doRegime (bot *botgolang.Bot, chatId string, userId string) {
	text := `Режим работы:

ВТ-ВС 11.00-20.00 (до 22.00 в режиме читального зала)
ПН - выходной,
последняя пятница месяца - санитарный день`
	regime := bot.NewTextMessage(choose(chatId, userId), text)
	addButtons(regime)
	_ = regime.Send()
}

func doHelp (bot *botgolang.Bot, chatId string, userId string) {
	text := `Бот показывает анонсы мероприятий библиотеки, её контакты и режим работы.
Кроме того, бот ищет книги или статьи в электронном каталоге. Для поиска введите ключевое слово (например, черемша), заглавие книги (например, Голодные игры) или фамилию автора (например, Акунин)
Для расширения поиска используйте усечение окончаний слов (черемша → черемш).`
	help := bot.NewTextMessage(choose(chatId, userId), text)
	addButtons(help)
	_ = help.Send()
}

func checkIrbisConnection() {
	connection := getConnection()
	if !connection.Connect() {
		log.Fatal("Can't connect")
	}
	log.Println("Подключились к ИРБИС64")
	log.Println("\tбаза данных=", connection.Database)
	log.Println("\tMax MFN=", connection.GetMaxMfn(connection.Database))
	_ = connection.Disconnect()
}

func doSearch(bot *botgolang.Bot, message *botgolang.Message) {
	query := message.Text
	if strings.EqualFold(query, "/start") {
		doStart(bot, message.Chat.ID, "")
		return
	}
	if strings.EqualFold(query, "анонсы") ||
		strings.EqualFold(query, "/announces") {
		doAnnounces(bot, message.Chat.ID, "")
		return
	}
	if strings.EqualFold(query, "контакты") ||
		strings.EqualFold(query, "/contacts") {
		doContacts(bot, message.Chat.ID, "")
		return
	}
	if strings.EqualFold(query, "режим работы") ||
		strings.EqualFold(query, "/regime") {
		doRegime(bot, message.Chat.ID, "")
		return
	}
	if strings.EqualFold(query, "помощь") ||
		strings.EqualFold(query, "/help") {
		doHelp(bot, message.Chat.ID, "")
		return
	}

	_ = message.Reply("Ищу книги и статьи...")

	connection := getConnection()
	if !connection.Connect() {
		_ = message.Reply("Ошибка связи с ИРБИС")
		return
	}

	parameters := irbis.SearchParameters {
		Database: connection.Database,
		Expression: "\"K=" + query + "$\"",
		Format: "@brief",
		FirstRecord: 1,
		NumberOfRecords: 10,
	}

	found := connection.SearchEx(&parameters)
	if len(found) == 0 {
		reply := bot.NewTextMessage(message.Chat.ID,
			"К сожалению, ничего не найдено")
		addButtons(reply)
		_ = reply.Send()
	}
	for i := range found {
		book := bot.NewTextMessage(message.Chat.ID,
			found[i].Description)
		if i == len(found) - 1 {
			addButtons(book)
		}
		_ = book.Send()
	}

	_ = connection.Disconnect()
}

func doCallback(bot *botgolang.Bot, userId string, callback *botgolang.ButtonResponse) {
	callbackData := callback.CallbackData
	response := bot.NewButtonResponse(callback.QueryID, "", callbackData, false)
	_ = response.Send()

	if strings.EqualFold(callbackData, "/start") {
		doStart(bot, "", userId)
		return
	}
	if strings.EqualFold(callbackData, "анонсы") ||
		strings.EqualFold(callbackData, "/announces") {
		doAnnounces(bot, "", userId)
		return
	}
	if strings.EqualFold(callbackData, "контакты") ||
		strings.EqualFold(callbackData, "/contacts") {
		doContacts(bot, "", userId)
		return
	}
	if strings.EqualFold(callbackData, "режим работы") ||
		strings.EqualFold(callbackData, "/regime") {
		doRegime(bot, "", userId)
		return
	}

	doHelp(bot, "", userId)
}

func runBot() {
	config := readConfig()
	bot, err := botgolang.NewBot(config.Token)
	if err != nil {
		log.Fatalf("cannot connect to bot: %s", err)
	}

	log.Println(bot.Info)

	ctx := context.Background()
	for {
		updates := bot.GetUpdatesChannel(ctx)
		for update := range updates {
			fmt.Println(update.Type, update.Payload)

			switch update.Type {
			case botgolang.NEW_MESSAGE:
				message := update.Payload.Message()
				go doSearch(bot, message)

			case botgolang.CALLBACK_QUERY:
				callback := update.Payload.CallbackQuery()
				userId := update.Payload.From.UserID
				go doCallback(bot, userId, callback)

			default:
				fmt.Println("Unknown message type")
			}
		}
	}
}


func main() {
	checkIrbisConnection()

	runBot()
}
