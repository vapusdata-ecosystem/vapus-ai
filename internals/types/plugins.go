package types

type PluginServices string

func (x PluginServices) String() string {
	return string(x)
}

const (
	SENDGRID          PluginServices = "SENDGRID"
	MAILGUN           PluginServices = "MAILGUN"
	SMTP              PluginServices = "SMTP"
	TWILIO            PluginServices = "TWILIO"
	AWS_SES           PluginServices = "AWS_SES"
	AZURE_EMAIL       PluginServices = "AZURE_EMAIL"
	AZURE_SMS         PluginServices = "AZURE_SMS"
	GMAIL             PluginServices = "GMAIL"
	GOOGLE_DRIVE      PluginServices = "GOOGLE_DRIVE"
	SLACK             PluginServices = "SLACK"
	SERPSEARCH        PluginServices = "SERP"
	Bedrock_Guardrail PluginServices = "AWS_Berock"
	Pangea_Guardrail  PluginServices = "Pangea"
	Mistral_Guardrail PluginServices = "Mistral"
)
