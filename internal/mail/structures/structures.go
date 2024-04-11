package structures

type EmailRecipient struct {
	Email    string
	Subject  string
	Messages MessagesData
}

type MessagesData struct {
	PreHeader string
	Messages  []string
}
