package helpers

import "time"

func TheTime() time.Time {
	return time.Now().UTC()
}
