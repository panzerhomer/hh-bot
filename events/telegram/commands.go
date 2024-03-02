package telegram

import (
	"HHBot/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	SearchCmd   = "/search"
	SettingsCmd = "/settings"
	HelpCmd     = "/help"
	StartCmd    = "/start"
	MeCmd       = "/me"
)

const (
	exp1 = "Нет опыта"
	exp2 = "От 1 года до 3 лет"
	exp3 = "От 3 года до 6 лет"
	exp4 = "Более 6 лет"
)

const (
	sal1 = "з/п не указана"
)

var (
	firstMenu = "menu with a shiny inline button."

	nextButton  = "Next"
	backButton  = "Back"
	closeButton = "Close"

	replyMarkupInlineKeyboard = map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": backButton, "callback_data": backButton},
				{"text": nextButton, "callback_data": nextButton},
			},
			{
				{"text": closeButton, "callback_data": closeButton},
			},
		},
		"resize_keyboard": true,
	}

	replyMarkupKeyboardExp = map[string]interface{}{
		"keyboard": [][]map[string]interface{}{
			{
				{"text": exp1, "callback_data": exp1},
				{"text": exp2, "callback_data": exp2},
			},
			{
				{"text": exp3, "callback_data": exp3},
				{"text": exp4, "callback_data": exp4},
			},
		},
		"resize_keyboard":   true,
		"one_time_keyboard": true,
	}

	emptyBoard = map[string]interface{}{
		"keyboard": [][]map[string]interface{}{},
	}
)

func createVacancy(v *models.Vacancy) string {
	text := fmt.Sprintf(`
<strong>Вакансия:</strong> %s  
<strong>Компания:</strong> %s  
<strong>Город:</strong> %s 
<strong>Зарплата:</strong> %s  
<strong>Опыт:</strong> %s 
<strong>Навыки:</strong> %s
<a href="%s">Cсылка на вакансию</a>
`,
		v.Title, v.Company, v.City, v.Salary, v.Experience, v.Skills, v.URL)
	return text
}

func createMenu(res []models.Vacancy) string {
	vacancies := ""
	for i := 0; i < len(res); i++ {
		vacancies += createVacancy(&res[i])
	}
	return vacancies
}

func createMe(userID int, f *models.Filter) string {
	id := strconv.Itoa(userID)
	textMe := fmt.Sprintf(`
<strong>My settings ⚙️</strong> 

<strong>ID:</strong> %s 
<strong>City:</strong> %s 
<strong>Salary:</strong> %s 
<strong>Experience:</strong> %s 
	`, id, f.City, f.Salary, f.Experience)
	return textMe
}

func (p *Processor) doCmd(text string, chatID int, userID int, username string, callbackID string, data string) error {
	command := strings.TrimSpace(text)
	p.mutex.Lock()
	defer p.mutex.Unlock()

	log.Printf("got new command '%s' from '%s", command, username)
	log.Printf("got new click '%s' from '%s", data, callbackID)

	_, ok := p.users[userID]
	if !ok {
		p.users[userID] = &UserData{StateDefault, []string{}}
	}

	switch command {
	case MeCmd:
		filter := p.sendMe(chatID, userID)
		if filter != nil {
			text := createMe(userID, filter)
			return p.sendMessage(chatID, text)
		}
		return p.sendMessage(chatID, "error")
	case StartCmd:
		p.users[userID].State = StateDefault
		return p.sendMessage(chatID, msgHello)
	case SettingsCmd:
		p.users[userID].State = StateSettingsQuestion1
		return p.sendMessage(chatID, msgCity)
	case SearchCmd:
		p.users[userID].State = StateSearch
		return p.sendMessage(chatID, msgSearch)
	default:
		return p.handleMessage(chatID, userID, text, callbackID, data)
	}
}

func (p *Processor) handleMessage(chatID int, userID int, text string, callbackID string, data string) error {
	currentUser := p.users[userID]

	switch currentUser.State {
	case StateSettingsQuestion1:
		currentUser.State = StateSettingsQuestion2
		currentUser.Answers = append(currentUser.Answers, text)
		return p.sendMessage(chatID, msgSalary)
	case StateSettingsQuestion2:
		currentUser.State = StateSettingsQuestion3
		currentUser.Answers = append(currentUser.Answers, text)
		return p.sendKeyboard(chatID, msgExperience, replyMarkupKeyboardExp)
	case StateSettingsQuestion3:
		currentUser.State = StateSearch
		currentUser.Answers = append(currentUser.Answers, text)
		f := models.Filter{}
		f.UserID = userID
		f.City = currentUser.Answers[0]
		f.Salary = currentUser.Answers[1]
		f.Experience = currentUser.Answers[2]
		p.storage.SetSettings(context.Background(), &f)
		return p.sendKeyboard(chatID, msgSettingsSuccess, emptyBoard)
	case StateSearch:
		f := models.Filter{}
		if len(currentUser.Answers) > 0 {
			f = createFilter(userID, currentUser.Answers)
		} else {
			f, _ = p.storage.GetSettings(context.Background(), userID)
		}
		return p.sendVacancies(chatID, userID, text, f, replyMarkupInlineKeyboard)
	default:
		return p.sendMessage(chatID, msgUnknownCommand)
	}
}

func createFilter(userID int, answers []string) models.Filter {
	filter := models.Filter{}
	filter.UserID = userID
	filter.City = answers[0]
	filter.Salary = answers[1]
	filter.Experience = answers[2]
	return filter
}

func (p *Processor) sendMe(chatID int, userID int) *models.Filter {
	filter, _ := p.storage.GetSettings(context.Background(), userID)
	return &filter
}

func (p *Processor) sendKeyboard(chatID int, text string, board any) error {
	replyMarkupJSON, err := json.Marshal(board)
	if err != nil {
		log.Fatal(err)
	}

	return p.tg.SendMessage(chatID, text, string(replyMarkupJSON))
}

func (p *Processor) sendVacancies(chatID int, userID int, text string, f models.Filter, keyboard any) error {
	search := strings.Trim(text, " ")
	res, err := p.storage.GetVacancies(context.Background(), &f, search)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		p.sendMessage(chatID, msgNoVacancies)
	}

	menu := createMenu(res)

	keyBoardJSON, _ := json.Marshal(keyboard)

	return p.tg.SendMessage(chatID, menu, string(keyBoardJSON))
}

func (p *Processor) sendMessage(chatID int, text string) error {
	return p.tg.SendMessage(chatID, text, "")
}
