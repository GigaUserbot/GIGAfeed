package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

const GIGA_FEED_CHAT_ID = -1001799797732

func main() {
	b, err := gotgbot.NewBot(BOT_TOKEN, &gotgbot.BotOpts{})
	if err != nil {
		panic("failed to create bot: " + err.Error())
	}
	webhookListener(b)
}

func webhookListener(b *gotgbot.Bot) {
	mux := http.NewServeMux()
	mux.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		type status struct {
			Alive       bool   `json:"alive"`
			BotId       int64  `json:"bot_id"`
			BotUsername string `json:"bot_username"`
		}
		jsonResp, _ := json.Marshal(
			status{
				Alive:       true,
				BotId:       b.Id,
				BotUsername: b.Username,
			},
		)
		w.Write(jsonResp)
	})
	mux.HandleFunc("/", processUpdate(b))
	server := &http.Server{
		Addr:        "0.0.0.0:" + PORT,
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
<b>Pusher's Email</b>: %s`

const ISSUE_TEMPL = `
<b><u><a href="github.com/gigauserbot">THE GIGA PROJECT</a></u></b>

<b><u>New Issue Update</u></b>
<b>Repository</b>: <code>%s</code>
<b>Action</b>: <code>%s</code>
<b>Issue</b>: <a href="%s">%s</a>
<b>By</b>: %s
`

const PR_TEMPL = `
<b><u><a href="github.com/gigauserbot">THE GIGA PROJECT</a></u></b>

<b><u>New PR Update</u></b>
<b>Repository</b>: <code>%s</code>
<b>Action</b>: <code>%s</code>
<b>Pull Request</b>: <a href="%s">%s</a>
<b>By</b>: %s
`

const DISC_TEMPL = `
<b><u><a href="github.com/gigauserbot">THE GIGA PROJECT</a></u></b>

<b><u>New Discussion Update</u></b>
<b>Repository</b>: <code>%s</code>
<b>Action</b>: <code>%s</code>
<b>Pull Request</b>: <a href="%s">%s</a>
<b>By</b>: %s
`

const COMNT_TEMPL = `
<b><u><a href="github.com/gigauserbot">THE GIGA PROJECT</a></u></b>

<b><u>New Comment Update</u></b>
<b>Repository</b>: <code>%s</code>
<b>Type</b>: <code>%s</code>
<b>Comment Link</b>: <a href="%s">click here</a>
<b>Action</b>: <code>%s</code>
<b>By</b>: %s
`

const REPO_TEMPL = `
<b><u><a href="github.com/gigauserbot">THE GIGA PROJECT</a></u></b>

<b><u>New Repository Update</u></b>
<b>Repository</b>: <code>%s</code>
<b>Action</b>: <code>%s</code>
<b>By</b>: %s
`

func handleUpdate(b *gotgbot.Bot, event *Event) {
	if event.Repository.Private {
		// Don't log info about private repos
		return
	}
	switch {
	case event.Comment.Url != "":
		var ctype string
		if event.Issue.Number != 0 {
			ctype = "Issue Comment"
		} else if event.Discussion.Number != 0 {
			ctype = "Discussion Comment"
		} else {
			ctype = "Commit Comment"
		}
		send(b,
			fmt.Sprintf(
				COMNT_TEMPL,
				event.Repository.Name,
				ctype,
				event.Comment.Url,
				event.Action,
				event.Sender.Name,
			),
			&gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Repository", Url: event.Repository.Url}},
				},
			},
		)
	case event.Ref != "":
		text := fmt.Sprintf(
			PUSH_TEMPL,
			event.Repository.Name,
			event.Ref,
			event.Compare,
			event.Pusher.Name,
			event.Pusher.Email,
		)
		if len(event.Commits) > 0 {
			text += "\n<b>Commits</b>:"
			for _, c := range event.Commits {
				text += fmt.Sprintf("\n  • <a href=\"%s\">%s</a>", c.Url, strings.Split(c.Message, "\n\n")[0])
			}
		}
		send(b,
			text,
			&gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Repository", Url: event.Repository.Url}},
				},
			},
		)
	case event.Issue.Number != 0:
		send(b,
			fmt.Sprintf(
				ISSUE_TEMPL,
				event.Repository.Name,
				event.Action,
				event.Issue.Url,
				event.Issue.Title,
				event.Sender.Name,
			),
			&gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Repository", Url: event.Repository.Url}},
				},
			},
		)
	case event.PullRequest.Number != 0:
		send(b,
			fmt.Sprintf(
				PR_TEMPL,
				event.Repository.Name,
				event.Action,
				event.PullRequest.Url,
				event.PullRequest.Title,
				event.Sender.Name,
			),
			&gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Repository", Url: event.Repository.Url}},
				},
			},
		)
	case event.Discussion.Number != 0:
		send(b,
			fmt.Sprintf(
				DISC_TEMPL,
				event.Repository.Name,
				event.Action,
				event.Discussion.Url,
				event.Discussion.Title,
				event.Sender.Name,
			),
			&gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Repository", Url: event.Repository.Url}},
				},
			},
		)
	case event.Repository.Name != "":
		send(b,
			fmt.Sprintf(
				REPO_TEMPL,
				event.Repository.Name,
				event.Action,
				event.Sender.Name,
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
		fmt.Println("failed to send message:", err.Error())
	}
}
