### Структуры MarcRecord, RecordField и SubField

#### MarcRecord

Каждый экземпляр `MarcRecord` соответствует одной записи в базе данных ИРБИС. Он содержит следующие поля:

Поле     |Тип                |Назначение
---------|-------------------|----------
Database | string            | Имя базы данных, из которой загружена данная запись. Для вновь созданных записей пустая строка.
Mfn      | int               | Номер записи в мастер-файле. Для вновь созданных записей 0.
Status   | int               | Статус записи: логически удалена, отсутствует (см. ниже).
Version  | int               | Номер версии записи.
Fields   | \[\]\*RecordField | Слайс указателей на поля записи

Статус записи: набор флагов (определены как константы в `Constants.go`)

Имя                     |Число|Значение
------------------------|-----|--------
LOGICALLY_DELETED       | 1   | Логически удалена (может быть восстановлена).
PHYSICALLY_DELETED      | 2   | Физически удалена (не может быть восстановлена).
ABSENT                  | 4   | Отсутствует.
NON_ACTUALIZED          | 8   | Не актуализирована.
NEW_RECORD              | 16  | Первый экземпляр записи (флаг фактически не используется).
LAST                    | 32  | Последняя версия записи.
LOCKED                  | 64  | Запись заблокирована на ввод.
AUTOIN_ERROR            | 128 | Ошибка в Autoin.gbl.
FULLTEXT_NOT_ACTUALIZED | 256 | Полный текст не актуализирован. 

**func NewMarcRecord() \*MarcRecord** -- конструктор, создаёт новый экземпляр записи в памяти клиента.

**func (record \*MarcRecord) Add(tag int, value string) \*RecordField** -- добавляет в конец записи поле с указанными меткой и значением. Возвращает добавленное поле, поэтому может использоваться для "цепочечных" вызовов методов, добавляющих подполя в это поле (см. пример ниже).

**func (record \*MarcRecord) Clear()** -- очищает запись (удаляет все поля).

**func (record \*MarcRecord) Clone() \*MarcRecord** -- клонирует запись со всеми полями.

**func (record \*MarcRecord) Decode(lines []string)** -- декодирование записи из протокольного представления.

**func (record \*MarcRecord) Encode(delimiter string) string** -- кодирование записи в протокольное представление.

**func (record \*MarcRecord) FM(tag int) string** -- получение значения поля с указанной меткой. Если поле не найдено, возвращается пустая строка.

**func (record \*MarcRecord) FSM(tag int, code rune) string** -- получение значения подполя с указанными меткой и кодом. Если подполе не найдено, возвращается пустая строка.

**func (record \*MarcRecord) FMA(tag int) []string** -- получение слайса со значениями полей с указанной меткой. Если поля не найдены, возвращается слайс нулевой длины.

**func (record \*MarcRecord) FSMA(tag int, code rune) []string** -- получение слайса со значениями подполей с указанными меткой и кодом. Если подполя не найдены, возвращается слайс нулевой длины.

**func (record \*MarcRecord) GetField(tag, occurrence int) \*RecordField** -- получение указанного повторения поля с указанной меткой. Если поле не найдено, возвращает `nil`.

**func (record \*MarcRecord) GetFields(tag int) []\*RecordField** -- получение слайса полей с указанной меткой. Если поля не найдены, возвращается пустой массив.

**func (record \*MarcRecord) GetFirstField(tag int) \*RecordField** -- получение первого повторения поля с указанной меткой. Если поле не найдено, возвращается `nil`.

**func (record \*MarcRecord) HaveField(tag int) bool** -- выясняет, есть ли в записи поле с указанной меткой.

**func (record \*MarcRecord) InsertAt(i int, tag int, value string) \*RecordField** -- вставляет поле в указанную позицию.

**func (record \*MarcRecord) IsDeleted() bool** -- проверка статуса, не удалена ли запись.

**func (record \*MarcRecord) RemoveAt(i int) \*MarcRecord** -- удаляет поле в указанной позиции.

**func (record \*MarcRecord) RemoveField(tag int) \*MarcRecord** -- удаляет все поля с указанной меткой.

**func (record \*MarcRecord) Reset()** -- сброс состояния записи, отвязка её от базы данных. Поля данных остаются при этом нетронутыми.

**func (record \*MarcRecord) SetField(tag int, value string) \*MarcRecord** -- устанавливает значение первого повторения поля с указанной меткой. Если такого поля нет, оно создаётся.

**func (record \*MarcRecord) SetSubfield(tag int, code rune, value string) \*MarcRecord** -- устанавливает значение подполя первого повторения поля с указанной меткой. Если необходимые поля или подполе отсутствуют, они создаются.

**func (record \*MarcRecord) String() string** -- выдаёт строковое представление записи.

```go
record := NewMarcRecord()
record.Add(700, "").
    Add('a', "Миронов").
    Add('b', "А. В.").
    Add('g', "Алексей Владимирович")
record.Add(200, "").
    Add('a', "Заглавие книги").
	Add('e', "подзаголовочные сведения")
```

#### RecordField

Поле      |Тип             |Назначение
----------|----------------|----------
Tag       | int            | Метка поля.
Value     | string         | Значение поля до первого разделителя.
Subfields | \[\]\*SubField | Слайс указателей на подполя.

**func NewRecordField(tag int, value string) \*RecordField** -- конструктор поля, создаёт поле с указанными меткой и значением.

**func (field \*RecordField) Add(code rune, value string) \*RecordField** -- добавляет подполе с указанными кодом и значением к полю. Возвращает `field`, так что может испольоваться для "цепочечных" вызовов.

**func (field \*RecordField) AddNonEmpty(code rune, value string) \*RecordField** -- добавляет подполе, при условии, что его значение не пустое.

**func (field \*RecordField) Clear() \*RecordField** -- очищает поле (удаляет значение и все подполя). Метка поля остаётся нетронутой. Возвращает `field`.

**func (field \*RecordField) Clone() \*RecordField** -- клонирует поле со всеми подполями.

**func (field \*RecordField) DecodeBody(body string)** -- декодирует только текст поля и подполей (без метки).

**func (field \*RecordField) Decode(text string)** -- декодирует поле из протокольного представления (метку, значение и подполя).

**func (field \*RecordField) Encode() string** -- кодирует поле в протокольное представление (метку, значение и подполя).

**func (field \*RecordField) EncodeBody() string** -- кодирует поле в протокольное представление (только значение и подполя).

**func (field \*RecordField) GetEmbeddedFields() []\*RecordField** -- получает слайс встроенных полей из данного поля.

**func (field \*RecordField) GetFirstSubField(code rune) \*SubField** -- возвращает первое вхождение подполя с указанным кодом или `nil`.

**func (field \*RecordField) GetFirstSubFieldValue(code rune) string** -- возвращает значение первого вхождения подполя с указанным кодом или пустую строку.

**func (field \*RecordField) GetValueOrFirstSubField() string** -- выдаёт значение для ^*.

**func (field \*RecordField) HaveSubField(code rune) bool** -- выясняет, есть ли подполе с указанным кодом.

**func (field \*RecordField) InsertAt(i int, code rune, value string)** -- вставляет подполе в указанную позицию.

**func (field \*RecordField) RemoveAt(i int)** -- удаляет подполе в указанной позиции.

**func (field \*RecordField) RemoveSubfield(code rune)** --удаляет все подполя с указанным кодом.

**func (field \*RecordField) ReplaceSubfield(code rune, oldValue, newValue string) \*RecordField** -- заменяет значение  подполя.

**func (field \*RecordField) SetSubfield(code rune, value string) \*RecordField** -- устанавливает значение первого повторения подполя с указанным кодом. Если value==nil, подполе удаляется.

**func (field \*RecordField) String() string** -- возвращает строковое представление данного поля.

**func (field \*RecordField) Verify() bool** -- проверяет, правильно ли сформировано поле (и все его подполя).

```go
field := NewRecordField(700, "")
field.Add('a', "Миронов").
    Add('b', "А. В.").
    Add('g', "Алексей Владимирович")
```

#### SubField

Поле|Тип|Назначение
----|---|----------
Code  | rune   | Код подполя
Value | string | Значение подполя

**func NewSubField(code rune, value string) \*SubField** -- конструктор подполя, создаёт подполе с указанными кодом и значением.

**func (subfield \*SubField) Clone() \*SubField** -- клонирует подполе.

**func (subfield \*SubField) Decode(text string)** -- декодирует подполе из протокольного представления.

**func (subfield \*SubField) Encode() string** -- кодирует подполе в протокольное представление.

**func (subfield \*SubField) String() string** -- выдаёт текстовое представление подполя.

**func (subfield \*SubField) Verify() bool** -- проверяет, правильно ли сформировано подполе.

```go
subfield := NewSubField('a', "Подполе A")
fmt.Println(subfield.String())
```

#### RawRecord

Запись с нераскодированными полями/подполями.

Поле|Тип|Назначение
---------|------------|----------
Database | string     | Имя базы данных, из которой загружена данная запись. Для вновь созданных записей пустая строка.
Mfn      | int        | Номер записи в мастер-файле. Для вновь созданных записей 0.
Status   | int        | Статус записи: логически удалена, отсутствует (аналогично `MarcRecord`).
Version  | int        | Номер версии записи.
Fields   | \[\]string | Слайс полей записи в "сыром" виде.

**func NewRawRecord() \*RawRecord** -- конструктор, создаёт новый экземпляр записи в памяти клиента.

**func (record \*RawRecord) Decode(lines []string)** -- декодирует запись из протокольного представления.

**func (record \*RawRecord) Encode(delimiter string) string** -- кодирует запись в протокольное представление.

**func (record \*RawRecord) IsDeleted() bool** -- проверка статуса, не удалена ли запись.

**func (record \*RawRecord) Reset()** -- сброс состояния записи, отвязка её от базы данных. Поля данных остаются при этом нетронутыми.

**func (record \*RawRecord) String() string** -- выдаёт строковое представление записи.

[Предыдущая глава](chapter2.md) [Следующая глава](chapter4.md)
