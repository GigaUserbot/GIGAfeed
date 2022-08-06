package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

const GIGA_FEED_CHAT_ID = -1001799797732
const BOT_TOKEN = "hard code it here"

func main() {
	b, err := gotgbot.NewBot(BOT_TOKEN, &gotgbot.BotOpts{})
	if err != nil {
		panic("failed to create bot: " + err.Error())
	}
	webhookListener(b)
}

func webhookListener(b *gotgbot.Bot) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", processUpdate(b))
	server := &http.Server{
		Addr:        "0.0.0.0:3455",
		Handler:     mux,
		ReadTimeout: time.Second * 2,
	}
	server.ListenAndServe()
}

func processUpdate(bot *gotgbot.Bot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("en error occured:", err.Error())
			return
		}
		var event = Event{}
		err = json.Unmarshal(b, &event)
		if err != nil {
			fmt.Println("failed to process update:", err.Error())
			return
		}
		handleUpdate(bot, &event)
	}
}

const PUSH_TEMPL = `
<b><u><a href="github.com/gigauserbot">THE GIGA PROJECT</a></u></b>

<b><u>New Push</u></b>
<b>Repository</b>: <code>%s</code>
<b>Ref</b>: <code>%s</code>
<b>Changes</b>: <a href="%s">click here</a>
<b>Pusher's Name</b>: %s
<b>Pusher's Email</b>: %s
`

const ISSUE_TEMPL = `
<b><u><a href="github.com/gigauserbot">THE GIGA PROJECT</a></u></b>

<b><u>New Issue Update</u></b>
<b>Repository</b>: <code>%s</code>
<b>Action</b>: <code>%s</code>
<b>Issue</b>: <a href="%s">%s</a>
<b>By</b>: %s
`

func handleUpdate(b *gotgbot.Bot, event *Event) {
	fmt.Println("ok")
	if event.Repository.Private {
		// Don't log info about private repos
		return
	}
	switch {
	case event.Ref != "":
		send(b,
			fmt.Sprintf(
				PUSH_TEMPL,
				event.Repository.Name,
				event.Ref,
				event.Compare,
				event.Pusher.Name,
				event.Pusher.Email,
			),
			&gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Repository", Url: event.Repository.Url}},
				},
			},
		)
	case event.Issue.Number != 0:
		fmt.Println("issue number")
		send(b,
			fmt.Sprintf(
				ISSUE_TEMPL,
				event.Repository.Name,
				event.Action,
				event.Issue.Url,
				event.Issue.Title,
				event.Issue.User.Name,
			),
			&gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Repository", Url: event.Repository.Url}},
				},
			},
		)
	}
}

func send(b *gotgbot.Bot, text string, markup *gotgbot.InlineKeyboardMarkup) {
	_, err := b.SendMessage(GIGA_FEED_CHAT_ID, text, &gotgbot.SendMessageOpts{
		ParseMode:             "html",
		ReplyMarkup:           markup,
		DisableWebPagePreview: true,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
