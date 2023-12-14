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
)

var (
	// Menu texts
	firstMenu = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	// secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	nextButton  = "Next"
	backButton  = "Back"
	closeButton = "Close"

	replyMarkup = map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": backButton, "callback_data": backButton},
				{"text": nextButton, "callback_data": nextButton},
			},
			{
				{"text": closeButton, "callback_data": closeButton},
			},
		},
	}
)

func createMenu(res []models.Vacancy) string {
	v := res[0]
	textMenu := fmt.Sprintf("<strong>Вакансия:</strong> %s \n\n <strong>Компания:</strong> %s \n <strong>Город:</strong> %s \n <strong>Зарплата:</strong> %s \n <strong>Опыт:</strong> %s \n <strong>Ссылка:</strong> %s", v.Title, v.Company, v.City, v.Salary, v.Experience, v.URL)

	return textMenu
}

func (p *Processor) doCmd(text string, chatID int, userID int, username string, callbackID string, data string) error {
	// command := strings.TrimSpace(text)
	command := strings.Split(text, " ")

	log.Printf("got new command '%s' from '%s", command[0], username)
	log.Print(command[1])
	log.Printf("got new click '%s' from '%s", data, callbackID)

	switch command[0] {
	case StartCmd:
		return p.sendHello(chatID, msgHello)
	case HelpCmd:
		return p.sendMenu(chatID, username, callbackID, data, replyMarkup)
	case SearchCmd:
		return p.sendVacancies(chatID, text)

		// return p.sendHello(chatID, command)
	// case SettingsCmd:
	// 	return p.sendSettings(chatID)
	default:
		// return p.sendMenu(chatID, username, callbackID, data, replyMarkup)
		return p.sendHello(chatID, msgUnknownCommand)
	}
	// if strings.HasPrefix(text, "/") {
	// 	return p.processCommand(chatID, text, username)
	// } else if callbackID != "" {
	// 	return p.proccessButton(chatID, callbackID, data)
	// } else {
	// 	return p.tg.SendMessage(chatID, msgUnknownCommand, "")
	// }
}

func (p *Processor) proccessButton(chatID int, callbackID string, data string) error {
	return nil
}

// func (p *Processor) processCommand(chatID int, text string, username string) error {
// 	command := ""
// 	if strings.HasPrefix(text, SearchCmd) {
// 		command = SearchCmd
// 	} else if strings.HasPrefix(text, HelpCmd) {
// 		command = HelpCmd
// 	} else if strings.HasPrefix(text, SettingsCmd) {
// 		command = SettingsCmd
// 	}

// 	switch command {
// 	case StartCmd:
// 		return p.sendHello(chatID, msgHello)
// 	case HelpCmd:
// 		return p.sendMenu(chatID, username, callbackID, data, replyMarkup)
// 	case SearchCmd:
// 		// return p.sendVacancies(chatID, text)

// 		return p.sendHello(chatID, command)
// 	// case SettingsCmd:
// 	// 	return p.sendSettings(chatID)
// 	default:
// 		// return p.sendMenu(chatID, username, callbackID, data, replyMarkup)
// 		return p.sendHello(chatID, msgUnknownCommand)
// 	}
// }

func (p *Processor) sendMenu(chatID int, username string, callbackID string, data string, replyMarkup any) error {
	replyMarkupJSON, err := json.Marshal(replyMarkup)
	if err != nil {
		log.Fatal(err)
	}

	return p.tg.SendMessage(chatID, firstMenu, string(replyMarkupJSON))
}

func (p *Processor) sendVacancies(chatID int, text string) error {
	filter := models.Filter{}
	filter.City = "Москва"
	search := strings.Split(text, " ")
	res, err := p.storage.GetVacancies(context.Background(), filter, "аналитик")
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

func (p *Processor) sendHello(chatID int, text string) error {
	return p.tg.SendMessage(chatID, text, "")
}
