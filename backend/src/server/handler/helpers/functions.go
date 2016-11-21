package helpers

func Error(error string) string {
	return `{"error":"` + error + `"}`
}
