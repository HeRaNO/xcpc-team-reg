package redis

import "fmt"

func makeEmailUserIDKey(email *string) string {
	return fmt.Sprintf("UID:%s", *email)
}

func makeEmailTokenKey(email *string) string {
	return fmt.Sprintf("EMAILTOKEN:%s", *email)
}

func makeEmailRequestKey(email *string) string {
	return fmt.Sprintf("EMAILREQ:%s", *email)
}

func makeEmailActionKey(email *string) string {
	return fmt.Sprintf("EMAILACTION:%s", *email)
}

func makeSessionKey(uuid *string) string {
	return fmt.Sprintf("SESSION:%s", *uuid)
}
