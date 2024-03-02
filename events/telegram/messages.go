package telegram

const msgHelp = `Я могу показывать вакансии сайта hh.ru!
/start - начало работы с ботом
/search - для поиска вакансии 
/settings - для установки фильтров
/me - мои настройки`

const (
	msgHello           = "<strong>Привет!</strong> 👾\n\n" + msgHelp
	msgSearch          = "<strong>Введите название профессии 🤓</strong>"
	msgCity            = "<strong>Введите название города 🤓</strong>"
	msgSalary          = "<strong>Введите желаемую зарплату 🤓</strong>"
	msgExperience      = "<strong>Выберите опыт работы 🤓</strong>"
	msgSettingsSuccess = "<strong>Фильтры заданы 🤓</strong>"
	msgNoVacancies     = "<strong>Нет таких вакансий 😔</strong>"
	msgUnknownCommand  = "<strong>Неизвестная команда</strong> 🤔"
)
