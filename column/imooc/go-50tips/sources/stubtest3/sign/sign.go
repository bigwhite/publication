package sign

import "time"

func Get(sender string) string {
	now := time.Now().Format(time.RFC1123)
	return now
}
