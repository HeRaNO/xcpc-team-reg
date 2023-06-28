package email

import "fmt"

func MakeStuEmail(stuid *string) string {
	return fmt.Sprintf("%s%s", *stuid, stuEmailSuffix)
}
