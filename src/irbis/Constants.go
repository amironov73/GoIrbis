package irbis

const (

	// Статус записи

	LOGICALLY_DELETED       = 1   // Запись логически удалена
	PHYSICALLY_DELETED      = 2   // Запись физически удалена
	ABSENT                  = 4   // Запись отсутствует
	NON_ACTUALIZED          = 8   // Запись не актуализирована
	NEW_RECORD              = 16  // Первый экземпляр записи
	LAST_VERSION            = 32  // Последняя версия записи
	LOCKED_RECORD           = 64  // Запись заблокирована на ввод
	AUTOIN_ERROR            = 128 // Ошибка в Autoin.gbl
	FULLTEXT_NOT_ACTUALIZED = 256 // Полный текст не актуализирован

	// Распространённые форматы

	ALL_FORMAT       = "&uf('+0')" // Полные данные по полям
	BRIEF_FORMAT     = "@brief"    // Краткое библиографическое описание
	IBIS_FORMAT      = "@ibiskw_h" // Формат IBIS (старый)
	INFO_FORMAT      = "@info_w"   // Информационный формат
	OPTIMIZED_FORMAT = "@"         // Оптимизированный формат

	// Распространённые поиски

	KEYWORD_PREFIX    = "K="  // Ключевые слова
	AUTHOR_PREFIX     = "A="  // Индивидуальный автор, редактор, составитель
	COLLECTIVE_PREFIX = "M="  // Коллектив или мероприятие
	TITLE_PREFIX      = "T="  // Заглавие
	INVENTORY_PREFIX  = "IN=" // Инвентарный номер, штрих-код или радиометка
	INDEX_PREFIX      = "I="  // Шифр документа в базе

	// Логические операторы для поиска

	LOGIC_OR                = 0 // Только ИЛИ
	LOGIC_OR_AND            = 1 // ИЛИ и И
	LOGIC_OR_AND_NOT        = 2 // ИЛИ, И, НЕТ (по умолчанию)
	LOGIC_OR_AND_NOT_FIELD  = 3 // ИЛИ, И, НЕТ, И (в поле)
	LOGIC_OR_AND_NOT_PHRASE = 4 // ИЛИ, И, НЕТ, И (в поле), И (фраза)

	// Коды АРМ

	ADMINISTRATOR = "A" // Адмнистратор
	CATALOGER     = "C" // Каталогизатор
	ACQUSITIONS   = "M" // Комплектатор
	READER        = "R" // Читатель
	CIRCULATION   = "B" // Книговыдача
	BOOKLAND      = "B" // Книговыдача
	PROVISITON    = "K" // Книгообеспеченность

	// Команды глобальной корректировки

	ADD_FIELD        = "ADD"
	DELETE_FIELD     = "DEL"
	REPLACE_FIELD    = "REP"
	CHANGE_FIELD     = "CHA"
	CHANGE_WITH_CASE = "CHAC"
	DELETE_RECORD    = "DELR"
	UNDELETE_RECORD  = "UNDELR"
	CORRECT_RECORD   = "CORREC"
	CREATE_RECORD    = "NEWMFN"
	EMPTY_RECORD     = "EMPTY"
	UNDO_RECORD      = "UNDOR"
	GBL_END          = "END"
	GBL_IF           = "IF"
	GBL_FI           = "FI"
	GBL_ALL          = "ALL"
	GBL_REPEAT       = "REPEAT"
	GBL_UNTIL        = "UNTIL"
	PUTLOG           = "PUTLOG"

	IRBIS_DELIMITER = "\x1F\x1E" // Разделитель строк в ИРБИС
	SHORT_DELIMITER = "\x1E"     // Короткая версия разделителя строк

)
