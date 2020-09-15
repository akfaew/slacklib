package slacklib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	. "github.com/akfaew/aeutils"
	"github.com/akfaew/utils"
	. "github.com/akfaew/webhandler"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"google.golang.org/appengine"
)

func InteractiveHandler(interactiveHandler func(context.Context, slack.InteractionCallback) (*slack.Msg, error)) APIHandler {
	return func(w http.ResponseWriter, r *http.Request) *APIError {
		ctx := appengine.NewContext(r)

		if r.Body == nil {
			return NewAPIError(http.StatusNotAcceptable, utils.Errorfc("Empty body"))
		}
		defer r.Body.Close()

		err := r.ParseForm()
		if err != nil {
			return NewAPIError(http.StatusNoContent, utils.Errorc(err))
		}

		payloadPost := r.PostFormValue("payload")
		if len(payloadPost) == 0 {
			return NewAPIError(http.StatusNoContent, utils.Errorfc("No Payload"))
		}

		var payload slack.InteractionCallback
		err = json.NewDecoder(strings.NewReader(payloadPost)).Decode(&payload)
		if err != nil {
			return NewAPIError(http.StatusNoContent, utils.Errorfc("json.Decode(): err=%v", err))
		}

		reply, err := interactiveHandler(ctx, payload)
		if err != nil {
			return NewAPIError(http.StatusInternalServerError, utils.Errorfc("actionDispatcher(): err=%v", err))
		}
		if reply != nil {
			b, err := json.Marshal(reply)
			if err != nil {
				return NewAPIError(http.StatusInternalServerError, utils.Errorfc("json.Marshal(): err=%v", err))
			}
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(b); err != nil {
				// let ServeHTTP() at least log this failure
				return NewAPIError(http.StatusInternalServerError, fmt.Errorf("Failed to write response: %v", err))
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte{}); err != nil {
				// let ServeHTTP() at least log this failure
				return NewAPIError(http.StatusInternalServerError, fmt.Errorf("Failed to write response: %v", err))
			}
		}

		return nil
	}
}

func CommandHandler(conf *Conf, commandHandler func(context.Context, *slack.SlashCommand) (*slack.Msg, error)) APIHandler {
	return func(w http.ResponseWriter, r *http.Request) *APIError {
		ctx := appengine.NewContext(r)

		s, err := slack.SlashCommandParse(r)
		if err != nil {
			return NewAPIError(http.StatusBadRequest, utils.Errorc(err))
		}

		if !s.ValidateToken(conf.VerificationToken) {
			return NewAPIError(http.StatusUnauthorized, utils.Errorfc("s.Token=%v", s.Token))
		}

		reply, err := commandHandler(ctx, &s)
		if err != nil {
			return NewAPIError(http.StatusInternalServerError, utils.Errorfc("err=%v", err))
		} else if reply == nil {
			return nil
		}

		b, err := json.Marshal(reply)
		if err != nil {
			return NewAPIError(http.StatusInternalServerError, utils.Errorfc("json.Marshal(): err=%v", err))
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			return NewAPIError(http.StatusInternalServerError, utils.Errorfc("w.Write(): err=%v", err))
		}

		return nil
	}
}

func EventsHandler(conf *Conf, messageEvent func(context.Context, *slackevents.MessageEvent, *slackevents.EventsAPICallbackEvent) *APIError) APIHandler {
	return func(w http.ResponseWriter, r *http.Request) *APIError {
		ctx := appengine.NewContext(r)

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r.Body); err != nil {
			return NewAPIError(http.StatusInternalServerError, utils.Errorfc("buf.ReadFrom(): err=%v", err))
		}
		body := buf.String()

		token := struct {
			Token string `json:"token"`
		}{}
		if err := json.Unmarshal(buf.Bytes(), &token); err != nil {
			return NewAPIError(http.StatusUnauthorized, utils.Errorfc("json.Unmarshal(%v): err=%v", buf.String(), err))
		}

		if conf.VerificationToken != token.Token {
			return NewAPIError(http.StatusUnauthorized, utils.Errorfc("cannot verify token %v", token.Token))
		}
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body),
			slackevents.OptionVerifyToken(&slackevents.TokenComparator{
				VerificationToken: token.Token, // We verified this manually above
			}))
		if err != nil {
			return NewAPIError(http.StatusInternalServerError, utils.Errorfc("slackevents.ParseEvent(%+v): err=%v", body, err))
		}

		switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				return NewAPIError(http.StatusInternalServerError, utils.Errorfc("json.Unmarshal(%+v): err=%v", body, err))
			}
			w.Header().Set("Content-Type", "text")
			if _, err := w.Write([]byte(r.Challenge)); err != nil {
				return NewAPIError(http.StatusInternalServerError, utils.Errorfc("w.Write(): err=%v", err))
			}
		case slackevents.AppRateLimited:
			LogCriticalfd(ctx, "Rate limiting")
			return NewAPIError(http.StatusTooManyRequests, utils.Errorfc("Rate limiting"))
		case slackevents.CallbackEvent:
			switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				ce := eventsAPIEvent.Data.(*slackevents.EventsAPICallbackEvent)

				return messageEvent(ctx, ev, ce)
			}
		}

		return nil
	}
}
