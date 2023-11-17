package app

import (
	"sync"

	"github.com/MithraRa/hh-bot/config"
	repository "github.com/MithraRa/hh-bot/internal/repository/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

const menuText = `Название вакансии: /search\n
	Удалить историю поиска: /clear`

type App struct {
	config *config.Config
	repo   *repository.Repository
	Bot    *tgbotapi.BotAPI
	mt     *sync.RWMutex
	Req    map[int64]string
}

func NewApp(cfg *config.Config) *App {
	repo, _ := repository.New(cfg)
	return &App{
		config: cfg,
		repo:   repo,
		Bot:    createBot(cfg.TokenTG),
		mt:     &sync.RWMutex{},
		Req:    make(map[int64]string, 0),
	}
}

func createBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	bot.Debug = false
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Logged in as ", bot.Self.UserName)

	return bot
}

func (a *App) Run() {
	log.Infoln("The bot is running...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := a.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			a.Handler(update.Message.Chat.ID, update.Message.Text)
		}
	}
}

func (a *App) Handler(id int64, msg string) {

}
