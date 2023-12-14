package telegram

import (
	"HHBot/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	// Menu texts
	firstMenu = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	// secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	nextButton  = "Next"
	backButton  = "Back"
	closeButton = "Close"

	// replyMarkupInlineKeyboard = map[string]interface{}{
	// 	"inline_keyboard": [][]map[string]interface{}{
	// 		{
	// 			{"text": backButton, "callback_data": backButton},
	// 			{"text": nextButton, "callback_data": nextButton},
	// 		},
	// 		{
	// 			{"text": closeButton, "callback_data": closeButton},
	// 		},
	// 	},
	// 	"resize_keyboard": true,
	// }
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

func createMenu(res []models.Vacancy) string {
	v := res[0]
	textMenu := fmt.Sprintf("<strong>Вакансия:</strong> %s <br/> <strong>Компания:</strong> %s <br/> <strong>Город:</strong> %s <br/> <strong>Зарплата:</strong> %s <br/> <strong>Опыт:</strong> %s <br/> <strong>Ссылка:</strong> %s", v.Title, v.Company, v.City, v.Salary, v.Experience, v.URL)

	return textMenu
}

func (p *Processor) doCmd(text string, chatID int, userID int, username string, callbackID string, data string) error {
	command := strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", command, username)
	log.Printf("got new click '%s' from '%s", data, callbackID)
	_, ok := p.users[userID]
	if !ok {
		p.users[userID] = &UserData{StateDefault, []string{}}
	}

	log.Print("user state", p.users[userID])

	switch command {
	case MeCmd:
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if len(p.users[userID].Answers) == 0 {
			return p.sendMessage(chatID, username+"\n\nempty")
		}
		return p.sendMessage(chatID, strings.Join(p.users[userID].Answers, " "))
	case StartCmd:
		return p.sendMessage(chatID, msgHello)
	case SettingsCmd:
		p.users[userID].State = StateSettingsQuestion1
		return p.sendMessage(chatID, msgCity)
	case SearchCmd:
		p.users[userID].State = StateSearch
		return p.sendVacancies(chatID, text)
	default:
		return p.handleMessage(chatID, userID, text, callbackID, data)
	}
}

func (p *Processor) handleMessage(chatID int, userID int, text string, callbackID string, data string) error {
	switch p.users[userID].State {
	case StateSettingsQuestion1:
		p.users[userID].State = StateSettingsQuestion2
		p.users[userID].Answers = append(p.users[userID].Answers, text)
		return p.sendMessage(chatID, msgSalary)
	case StateSettingsQuestion2:
		p.users[userID].State = StateSettingsQuestion3
		p.users[userID].Answers = append(p.users[userID].Answers, text)
		return p.sendKeyboard(chatID, msgExperience)
	case StateSettingsQuestion3:
		p.users[userID].State = StateDefault
		p.users[userID].Answers = append(p.users[userID].Answers, text)
		return p.tg.SendMessage(chatID, msgSettingsSuccess, "")
	default:
		return p.sendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) sendKeyboard(chatID int, text string) error {
	replyMarkupJSON, err := json.Marshal(replyMarkupKeyboardExp)
	if err != nil {
		log.Fatal(err)
	}

	return p.tg.SendMessage(chatID, text, string(replyMarkupJSON))
}

func (p *Processor) sendMenu(chatID int, userID int, callbackID string, data string, replyMarkup any) error {
	replyMarkupJSON, err := json.Marshal(replyMarkup)
	if err != nil {
		log.Fatal(err)
	}

	return p.tg.SendMessage(chatID, firstMenu, string(replyMarkupJSON))
}

func (p *Processor) sendVacancies(chatID int, text string) error {
	filter := models.Filter{}
	filter.City = "Москва"
	search := strings.Trim(text, " ")
	res, err := p.storage.GetVacancies(context.Background(), filter, search)
	if err != nil {
		return err
	}
	log.Print("[sendVacancies]", res, search)
	menu := createMenu(res)
	return p.tg.SendMessage(chatID, menu, "")
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp, "")
}

func (p *Processor) sendMessage(chatID int, text string) error {
	return p.tg.SendMessage(chatID, text, "")
}
