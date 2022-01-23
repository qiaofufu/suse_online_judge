package util

import "fmt"

func CodeKey(email string) string {
	return fmt.Sprintf("code:%s",email)
}