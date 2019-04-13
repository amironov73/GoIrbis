### Структура Connection

Структура `Connection` - "рабочая лошадка". Она осуществляет связь с сервером и всю необходимую перепаковку данных из клиентского представления в сетевое.

Экземпляр клиента создается конструктором:

```go
connection := irbis.NewConnection()
```

При создании клиента можно указать (некоторые) настройки:

```go
client := irbis.NewConnection()
client.Host = "irbis.rsl.ru"
client.Port = 5555 // нестандартный порт!
client.Username = "ninja"
client.Password = "i_am_invisible"
```

Поле|Тип|Назначение|Значение по умолчанию
----|---|----------|---------------------
Host        |string  | Адрес сервера|"127.0.0.1"
Port        |int     | Порт|6666
Username    |string  | Имя (логин) пользователя|пустая строка
Password    |string  | Пароль пользователя|пустая строка
Database    |string  | Имя базы данных|"IBIS"
Workstation |string  | Тип АРМа (см. таблицу ниже)| "C"

Типы АРМов

Обозначение|Тип
-----------|---
"R" | Читатель
"C" | Каталогизатор
"M" | Комплектатор
"B" | Книговыдача
"K" | Книгообеспеченность
"A" | Администратор

Можно использовать мнемонические константы, определённые в файле `Constants.go`:

```go
const ADMINISTRATOR = "A" // Адмнистратор
const CATALOGER     = "C" // Каталогизатор
const ACQUSITIONS   = "M" // Комплектатор
const READER        = "R" // Читатель
const CIRCULATION   = "B" // Книговыдача
const BOOKLAND      = "B" // Книговыдача
const PROVISITON    = "K" // Книгообеспеченность
```

Обратите внимание, что адрес сервера задается строкой, так что может принимать как значения вроде `192.168.1.1`, так и `irbis.yourlib.com`.

Если какой-либо из вышеперечисленных параметров не задан явно, используется значение по умолчанию.

#### Подключение к серверу и отключение от него

Только что созданный клиент еще не подключен к серверу. Подключаться необходимо явно с помощью метода `Connect`, при этом можно указать параметры подключения:

```go
client := irbis.NewConnection()
client.Host = "myhost.com"
if !client.Connect() {
	log.Fatal("Не удалось подключиться!")
}
```

Отключаться от сервера необходимо с помощью метода `Disconnect`, желательно помещать вызов в блок defer сразу после подключения (чтобы не забыть):

```go
client.Connect()
defer client.Disconnect()
```

При подключении клиент получает с сервера INI-файл с настройками, которые могут понадобиться в процессе работы:

```go
client.Connect()
defer client.Disconnect()
// Получаем имя MNU-файла, хранящего перечень форматов
formatMenuName := client.Ini.GetValue("Main", "FmtMnu", "FMT31.MNU")
```

Полученный с сервера INI-файл хранится в поле `Ini`.

Повторная попытка подключения с помощью того же экземпляра `Connection` игнорируется. При необходимости можно создать другой экземпляр и подключиться с его помощью (если позволяют клиентские лицензии). Аналогично игнорируются повторные попытки отключения от сервера.

Проверить статус "клиент подключен или нет" можно с помощью поля `Connected`:

```go
if client.Connected {
    // В настоящее время мы не подключены к серверу
}
```

Вместо индивидуального задания каждого из полей `Host`, `Port`, `Username`, `Password` и `Database`, можно использовать метод `ParseConnectionString`:

```go
client.ParseConnectionString("host=192.168.1.4;port=5555;" +
         "username=itsme;password=secret;")
client.Connect()
``` 

#### Многопоточность

Клиент написан в наивном однопоточном стиле, поэтому не поддерживает одновременный вызов методов из разных потоков.

Для одновременной отсылки на сервер нескольких команд необходимо создать соответствующее количество экземпляров подключений (если подобное позволяет лицензия сервера).

#### Подтверждение подключения

`GoIrbis` самостоятельно не посылает на сервер подтверждений того, что клиент все еще подключен. Этим должно заниматься приложение, например, по таймеру. 

Подтверждение посылается серверу методом `NoOp`:
 
```go
client.NoOp()        
```

#### Чтение записей с сервера

```go
mfn := 123
record := client.ReadRecord(mfn)        
```

Можно прочитать несколько записей сразу:

```go
mfns := []int{12, 34, 56}
records := client.ReadRecords(mfns)
```

Можно прочитать определенную версию записи

```go
mfn := 123
version := 3
record := client.ReadRecordVersion(mfn, version)
```

#### Сохранение записи на сервере

```go
// Любым образом создаём в памяти клиента
// или получаем с сервера запись.
record := client.ReadRecord(123)

// Производим какие-то манипуляции над записью
record.Add(999, "123")

// Отсылаем запись на сервер
newMaxMfn := client.WriteRecord(record)
println("New Max MFN:", newMaxMfn)
```

Сохранение нескольких записей (возможно, из разных баз данных):

```go
records := make([]MarcRecord,10)
...
if !client.WriteRecords(records) {
    log.Fatal("Failure!")
}
```

#### Удаление записи на сервере

```go
mfn := 123
client.deleteRecord(mfn)
```

Восстановление записи:

```go
mfn := 123
record := client.UndeleteRecord(mfn)
```

#### Поиск записей

```go
found := client.Search(`"A=ПУШКИН$"`)
println("Найдено записей:", len(found))
```

Обратите внимание, что поисковый запрос заключен в дополнительные кавычки. Эти кавычки явлются элементом синтаксиса поисковых запросов ИРБИС64, и лучше их не опускать.

Вышеприведённый запрос вернёт не более 32 тыс. найденных записей. Сервер ИРБИС64 за одно обращение к нему может выдать не более 32 тыс. записей. Чтобы получить все записи, используйте метод SearchAll (см. ниже), он выполнит столько обращений к серверу, сколько нужно.

Поиск с одновременной загрузкой записей:

```go
records := client.SearchRead(`"A=ПУШКИН$"`, 50)
println("Найдено записей:", len(records))
```

Поиск и загрузка единственной записи:

```go
record := client.SearchSingleRecord(`"I=65.304.13-772296"`)
if record == nil {
    println("Не нашли!")
}
```

Количество записей, соответствующих поисковому выражению:

```go
expression := `"A=ПУШКИН$"`
count := client.SearchCount(expression)
```

Расширенный поиск: можно задать не только количество возвращаемых записей, но и расформатировать их.

```go
parameters := NewSearchParameters()
parameters.Expression = `"A=ПУШКИН$"`
parameters.Format = BRIEF_FORMAT
parameters.NumberOfRecords = 5
found := client.SearchEx(parameters)
if len(found) == 0 {
	println("Не нашли")
} else {
    // в found находится слайс структур FoundLine
    first := found[0]
    fmt.Println("MFN:", first.Mfn, "DESCRIPTION:", first.Description)
}
```

Поиск всех записей (даже если их окажется больше 32 тыс.):

```go
found := client.SearchAll(`"A=ПУШКИН$"`)
println("Найдено записей:", len(count))
```

Подобные запросы следует использовать с осторожностью, т. к. они, во-первых, создают повышенную нагрузку на сервер, и во-вторых, потребляют очень много памяти на клиенте. Некоторые запросы (например, "I=$") могут вернуть все записи в базе данных, а их там может быть десятки миллионов.

#### Форматирование записей

```go
mfn := 123
format := BRIEF_FORMAT
text := client.FormatRecord(format, mfn)
println("Результат форматирования:", text)
```

При необходимости можно использовать в формате все символы UNICODE:

```go
mfn := 123
format := "'Ἀριστοτέλης: ', v200^a"
text := client.FormatRecord(format, mfn)
println("Результат форматирования:", text)
```

Форматирование нескольких записей:

```go
mfns := []int {12, 34, 56}
format := BRIEF_FORMAT
lines := client.FormatRecords(format, mfns)
fmt.Println("Результаты:", lines)
```

#### Печать таблиц

```go
table := new(TableDefinition)
table.Database = "IBIS"
table.Table = "@tabf1w"
table.SearchQuery = `"T=A$"`
text := client.PrintTable(table)
```

#### Работа с контекстом

Функция | Назначение
--------|-----------
ListFiles | Получение списка файлов на сервере
ReadIniFile | Получение INI-файла с сервера
ReadMenuFile | Получение MNU-файла с сервера
ReadSearchScenario | Загрузка сценариев поиска с сервера
ReadTextFile | Получение текстового файла с сервера
ReadTextLines | Получение текстового файла в виде массива строк
ReadTreeFile | Получение TRE-файла с сервера
UpdateIniFile | Обновление строк серверного INI-файла
WriteTextFile | Сохранение текстового файла на сервере

#### Работа с мастер-файлом

Функция | Назначение
--------|-----------
ReadRawRecord | Чтение указанной записи в "сыром" виде
WriteRawRecord | Сохранение на сервере "сырой" записи

#### Работа со словарем

Функция | Назначение
--------|-----------
ListTerms | Получение списка терминов с указанным префиксом
ReadPostings | Чтение постингов поискового словаря
ReadTerms | Чтение терминов поискового словаря
ReadTermsEx | Расширенное чтение терминов

#### Информационные функции

Функция | Назначение
--------|-----------
GetDatabaseInfo | Получение информации о базе данных
GetMaxMfn | Получение максимального MFN для указанной базы данных
GetServerVersion | Получение версии сервера
ListDatabases | Получение списка баз данных с сервера
ToConnectionString | Получение строки подключения

#### Администраторские функции

Нижеперечисленные записи доступны лишь из АРМ "Администратор", поэтому подключаться к серверу необходимо так:

```go
client = new IrbisConnection()
client.Username = 'librarian'
client.Password = 'secret'
client.Workstation = ADMINISTRATOR
if !client.Connect() {
	log.Fatal("Не удалось подключиться")
}
```

Функция | Назначение
--------|-----------
ActualizeDatabase | Актуализация базы данных
ActualizeRecord | Актуализация записи
CreateDatabase | Создание базы данных
CreateDictionary | Создание словаря
DeleteDatabase | Удаление базы данных
DeleteFile | Удаление файла на сервере 
GetServerStat | Получение статистики с сервера
GetUserList | Получение списка пользователей с сервера
ListProcesses | Получение списка серверных процессов
ReloadDictionary | Пересоздание словаря
ReloadMasterFile | Пересоздание мастер-файла
RestartServer | Перезапуск сервера
TruncateDatabase | Опустошение базы данных
UnlockDatabase | Разблокирование базы данных
UnlockRecords | Разблокирование записей
UpdateUserList | Обновление списка пользователей на сервере

#### Глобальная корректировка

```go
settings := new(GblSettings)
settings.Database = "IBIS"
settings.MfnList = []int{1, 2, 3}
settings->statements = []GblStatement {
  GblStatement{ADD_FIELD, "3000", "XXXXXXXXX", "'Hello'"}
}
result = connection.GlobalCorrection(settings)
for line := range result {
	println(line)
}
```

#### Расширение функциональности

**ExecuteAnyCommand(string $command, array $params)** -- выполнение произвольной команды с параметрами в кодировке ANSI.


[Предыдущая глава](chapter1.md) [Следующая глава](chapter3.md)
