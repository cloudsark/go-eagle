package alerts

import (
	"github.com/cloudsark/go-eagle/config"
	c "github.com/cloudsark/go-eagle/constants"
	"github.com/cloudsark/go-eagle/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/slack-go/slack"
)

const (
	PingUp          = " is up"
	PingDown        = " is down"
	ValidSsl        = " SSL certificate is valid"
	SslExpiredDate1 = " will be expired in "
	SslExpired      = " is expired"
	SslExpiredDate2 = " Days"
	CheckPort       = "Port:"
	CheckPortUp     = " is up on "
	CheckPortDown   = " is down on "
	LoadAvgMsg1     = "Load on "
	LoadAvgMsg2     = " is "
	DiskCritical    = " space is critical on "
	DiskNormal      = " space is normal on "
)

var snsARN = c.OSEnv("SNS_ARN")
var slackToken = c.OSEnv("SLACK_TOKEN")
var slackChannel = c.OSEnv("SLACK_CHANNEL")

/*
sendslack function posts message to a slack channel
Input: Bot User OAuth Access Token
Output: A success/error message
todo: add timestamp to attachment
*/
func sendslack(token, channel, text,
	color, priority string) {
	api := slack.New(token)
	attachment := slack.Attachment{
		Text:  text,
		Color: color,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Priority",
				Value: priority,
				Short: false,
			},
		},
	}

	channelID, timestamp, err := api.PostMessage(
		channel,
		slack.MsgOptionText("Alert", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		logger.ErrorLogger.Fatalf("%s\n", err)
		return
	}
	logger.GeneralLogger.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

}

func slackalerts(slackToken, slackChannel,
	domain, message, status string) {
	switch {
	case status == "PingUp":
		{
			sendslack(slackToken, slackChannel, message,
				"#008000", "Normal")
		}
	case status == "PingDown":
		{
			sendslack(slackToken, slackChannel, message,
				"#FF0000", "High")
		}
	case status == "SslValid":
		{
			sendslack(slackToken, slackChannel, message,
				"#008000", "Normal")
		}
	case status == "SslNotValidWarn":
		{
			sendslack(slackToken, slackChannel, message,
				"#FFA500", "Medium")
		}
	case status == "SslNotValidCrit":
		{
			sendslack(slackToken, slackChannel, message,
				"#FF0000", "High")
		}
	case status == "PortUp":
		{
			sendslack(slackToken, slackChannel, message,
				"#008000", "Normal")
		}
	case status == "PortDown":
		{
			sendslack(slackToken, slackChannel, message,
				"#FF0000", "High")
		}
	case status == "AvgLoadHigh":
		{
			sendslack(slackToken, slackChannel, message,
				"#FF0000", "High")
		}
	case status == "AvgLoadNormal":
		{
			sendslack(slackToken, slackChannel, message,
				"#008000", "Normal")
		}
	case status == "DiskCritical":
		{
			sendslack(slackToken, slackChannel, message,
				"#FF0000", "High")
		}
	case status == "DiskNormal":
		{
			sendslack(slackToken, slackChannel, message,
				"#008000", "Normal")
		}
	}
}

func sendmail(message string) {
	sess, err := session.NewSession(aws.NewConfig().WithRegion("eu-west-1"))

	if err != nil {
		logger.ErrorLogger.Println("NewSession error:", err)
		return
	}

	client := sns.New(sess)
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(snsARN),
	}

	result, err := client.Publish(input)
	if err != nil {
		logger.ErrorLogger.Println("Publish error:", err)
		return
	}
	logger.GeneralLogger.Printf("Notification successfully sent with an ID %s ", *result.MessageId)
}

func sendmailalert(message string) {
	sendmail(message)
}

// Alerter sends alerts if alert type is set to "on" in main.yaml file
func Alerter(slackToken, slackChannel,
	domain, message, status string) {
	if config.AlertStruct("Slack") == "true" {
		slackalerts(slackToken, slackChannel,
			domain, message, status)
	}
	if config.AlertStruct("Email") == "true" {
		sendmailalert(message)
	}
}
