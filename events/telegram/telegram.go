package telegram

import (
	"HHBot/clients/telegram"
	"HHBot/events"
	"HHBot/storage"
	"HHBot/utils"
	"errors"
	"log"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserID   int
	Username string
}

type MetaCallback struct {
	ChatID     int
	UserID     int
	Username   string
	CallbackID string
	Data       string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, utils.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.CallbackQuery:
		return p.processCallback(event)
	default:
		return utils.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processCallback(event events.Event) error {
	meta, err := metaCallback(event)
	if err != nil {
		return utils.Wrap("can't process callback", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserID, meta.Username, meta.CallbackID, meta.Data); err != nil {
		return utils.Wrap("can't process callback", err)
	}

	return nil
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return utils.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserID, meta.Username, "", ""); err != nil {
		return utils.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, utils.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func metaCallback(event events.Event) (MetaCallback, error) {
	res, ok := event.Meta.(MetaCallback)
	if !ok {
		return MetaCallback{}, utils.Wrap("can't get metaCallback", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserID:   upd.Message.From.ID,
			Username: upd.Message.From.Username,
		}
		log.Println("[event upd]", upd.Message, updType)
	}

	if updType == events.CallbackQuery {
		res.Meta = MetaCallback{
			ChatID:     upd.CallbackQuery.Message.Chat.ID,
			UserID:     upd.CallbackQuery.From.ID,
			Username:   upd.CallbackQuery.From.Username,
			CallbackID: upd.CallbackQuery.ID,
			Data:       upd.CallbackQuery.Data,
		}

		// res, _ := json.Marshal(upd)
		// log.Println("[event upd]", string(res))
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}

	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
