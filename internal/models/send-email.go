package models

import (
	"fmt"
)

type SendEmailActivity struct{}

func (s *SendEmailActivity) Execute(task Task) error {
	to := task.Config["to"].(string)
	subject := task.Config["subject"].(string)
	body := task.Config["body"].(string)

	fmt.Printf("Sending email to %s\nSubject: %s\nBody: %s\n", to, subject, body)
	return nil
}
