package slacklib

import (
	"fmt"

	"github.com/slack-go/slack"
)

var (
	ColorSuccess = "#00BB00" // An operation succeded
	ColorError   = "#F35A00" // An operation failed
	ColorInfo    = "#000000" // Regular message
	ColorList    = "#999999" // An element in a list of items
)

type AttachmentT struct{}
type MsgT struct{}

var Attachment AttachmentT
var Msg MsgT

func (_ AttachmentT) Success(text string, a ...interface{}) slack.Attachment {
	return slack.Attachment{
		Color: ColorSuccess,
		Text:  fmt.Sprintf(text, a...),
	}
}

func (_ AttachmentT) Successt(title, text string, a ...interface{}) slack.Attachment {
	return slack.Attachment{
		Color: ColorSuccess,
		Title: title,
		Text:  fmt.Sprintf(text, a...),
	}
}

func (_ AttachmentT) Info(text string, a ...interface{}) slack.Attachment {
	return slack.Attachment{
		Color: ColorInfo,
		Text:  fmt.Sprintf(text, a...),
	}
}

func (_ AttachmentT) Infot(title, text string, a ...interface{}) slack.Attachment {
	return slack.Attachment{
		Color: ColorInfo,
		Title: title,
		Text:  fmt.Sprintf(text, a...),
	}
}

func (_ AttachmentT) Error(text string, a ...interface{}) slack.Attachment {
	return slack.Attachment{
		Color: ColorError,
		Text:  fmt.Sprintf(text, a...),
	}
}

func (_ AttachmentT) Errort(title, text string, a ...interface{}) slack.Attachment {
	return slack.Attachment{
		Color: ColorError,
		Title: title,
		Text:  fmt.Sprintf(text, a...),
	}
}

func (_ MsgT) MakeMsg(attachments ...slack.Attachment) *slack.Msg {
	return &slack.Msg{
		Attachments:     attachments,
		ReplaceOriginal: true, // this should be default, it's a slack-go/slack bug
	}
}

func (_ MsgT) Success(text string, a ...interface{}) *slack.Msg {
	return Msg.MakeMsg(slack.Attachment{
		Color: ColorSuccess,
		Text:  fmt.Sprintf(text, a...),
	})
}

func (_ MsgT) Info(text string, a ...interface{}) *slack.Msg {
	return Msg.MakeMsg(slack.Attachment{
		Color: ColorInfo,
		Text:  fmt.Sprintf(text, a...),
	})
}

func (_ MsgT) Error(text string, a ...interface{}) *slack.Msg {
	return Msg.MakeMsg(slack.Attachment{
		Color: ColorError,
		Text:  fmt.Sprintf(text, a...),
	})
}
