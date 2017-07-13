package slack

import (
	"github.com/nlopes/slack"
	"github.com/olivoil/standup-lambda/modules/standup"
	"github.com/olivoil/standup-parser"
)

// PostStatusToSlack posts a status to a channel in Slack
func PostStatusToSlack(token string, channel string, status *standup.Status) error {
	client := slack.New(token)

	attachments, err := ConvertStatementToAttachments(status.Statement)
	if err != nil {
		return err
	}

	params := slack.PostMessageParameters{
		AsUser:      true,
		Attachments: attachments,
		Parse:       "full",
	}

	if _, _, err := client.PostMessage(channel, "", params); err != nil {
		return err
	}

	return nil
}

// ConvertStatementToAttachments converts a Statement to slack message attachments
// so that it can be displayed within slack's UI
func ConvertStatementToAttachments(stmt *parser.Statement) (attachments []slack.Attachment, err error) {
	if stmt.Yesterday.Valid || stmt.Today.Valid {
		fields := []slack.AttachmentField{}

		if stmt.Yesterday.Valid {
			key := stmt.Yesterday.Key
			if key == "" {
				key = "Yesterday"
			}
			fields = append(fields, slack.AttachmentField{
				Title: key,
				Value: stmt.Yesterday.Val,
				Short: true,
			})
		}

		if stmt.Today.Valid {
			key := stmt.Today.Key
			if key == "" {
				key = "Today"
			}
			fields = append(fields, slack.AttachmentField{
				Title: key,
				Value: stmt.Today.Val,
				Short: false,
			})
		}

		attachments = append(attachments, slack.Attachment{
			Color:  "#5e8eb7",
			Fields: fields,
		})
	}

	if stmt.Meetings.Valid || stmt.Blockers.Valid {
		fields := []slack.AttachmentField{}

		if stmt.Meetings.Valid {
			fields = append(fields, slack.AttachmentField{
				Title: "Meetings",
				Value: stmt.Meetings.Val,
				Short: true,
			})
		}

		if stmt.Blockers.Valid {
			fields = append(fields, slack.AttachmentField{
				Title: "Blockers",
				Value: stmt.Blockers.Val,
				Short: true,
			})
		}

		attachments = append(attachments, slack.Attachment{
			Color:  "#6c6c6c",
			Fields: fields,
		})
	}

	if stmt.LP.Valid || stmt.Jira.Valid {
		fields := []slack.AttachmentField{}

		if stmt.LP.Valid {
			fields = append(fields, slack.AttachmentField{
				Title: "LP",
				Value: stmt.LP.Lit,
				Short: true,
			})
		}

		if stmt.Jira.Valid {
			fields = append(fields, slack.AttachmentField{
				Title: "Jira",
				Value: stmt.Jira.Lit,
				Short: true,
			})
		}

		attachments = append(attachments, slack.Attachment{
			Color:  "#549b57",
			Fields: fields,
		})
	}

	return
}
