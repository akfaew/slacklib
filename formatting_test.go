package slacklib

import (
	"testing"

	"github.com/akfaew/test"
)

func TestAttachments(t *testing.T) {
	test.FixtureExtra(t, "Attachment.Success", Attachment.Success("hello %s", "world"))
	test.FixtureExtra(t, "Attachment.Successt", Attachment.Successt("hello %s", "world"))
	test.FixtureExtra(t, "Attachment.Info", Attachment.Info("hello %s", "world"))
	test.FixtureExtra(t, "Attachment.Infot", Attachment.Infot("hello %s", "world"))
	test.FixtureExtra(t, "Attachment.Error", Attachment.Error("hello %s", "world"))
	test.FixtureExtra(t, "Attachment.Errort", Attachment.Errort("hello %s", "world"))

	test.FixtureExtra(t, "Msg.Success", Msg.Success("hello %s", "world"))
	test.FixtureExtra(t, "Msg.Info", Msg.Info("hello %s", "world"))
	test.FixtureExtra(t, "Msg.Error", Msg.Error("hello %s", "world"))
}
