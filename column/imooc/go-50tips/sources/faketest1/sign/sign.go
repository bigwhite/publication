package sign

import "time"

func Get() string {
	now := time.Now().Format(time.RFC1123)
	return now
}
