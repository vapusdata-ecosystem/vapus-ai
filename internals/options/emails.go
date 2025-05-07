package options

type SendEmailRequest struct {
	From             string            `json:"from"`
	SenderName       string            `json:"senderName"`
	To               []string          `json:"to"`
	Subject          string            `json:"subject"`
	Body             string            `json:"body"`
	Attachments      []*Attachment     `json:"attachments"`
	HtmlTemplateBody string            `json:"htmlTemplateBody"`
	RawMessage       []byte            `json:"rawMessage"`
	Tags             map[string]string `json:"tags"`
	ID               string            `json:"id"`
	CCList           []string          `json:"ccList"`
	BCCList          []string          `json:"bccList"`
}

type SendEmailResponse struct {
	Recievers   []string `json:"recievers"`
	Status      string   `json:"status"`
	Attachments []string `json:"attachments"`
	Message     string   `json:"message"`
}

type Attachment struct {
	FileName    string `json:"fileName"`
	Data        []byte `json:"data"`
	Format      string `json:"format"`
	ContentType string `json:"contentType"`
}

type ReadEmailRequest struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty" toml:"id,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
}

type ListEmailsRequest struct {
	RecievedOn     string   `json:"recievedOn,omitempty" yaml:"recievedOn,omitempty" toml:"recievedOn,omitempty"`
	SenderEmails   []string `json:"senderEmails,omitempty" yaml:"senderEmails,omitempty" toml:"senderEmails,omitempty"`
	Subject        string   `json:"subject,omitempty" yaml:"subject,omitempty" toml:"subject,omitempty"`
	Body           string   `json:"body,omitempty" yaml:"body,omitempty" toml:"body,omitempty"`
	Attachments    []string `json:"attachments,omitempty" yaml:"attachments,omitempty" toml:"attachments,omitempty"`
	SearchKeywords []string `json:"searchKeywords,omitempty" yaml:"searchKeywords,omitempty" toml:"searchKeywords,omitempty"`
	SortBy         string   `json:"sortBy,omitempty" yaml:"sortBy,omitempty" toml:"sortBy,omitempty"`
	SortOrder      string   `json:"sortOrder,omitempty" yaml:"sortOrder,omitempty" toml:"sortOrder,omitempty"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Labels         []string `json:"labels,omitempty" yaml:"labels,omitempty" toml:"labels,omitempty"`
}

type DeleteEmailRequest struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty" toml:"id,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
}
