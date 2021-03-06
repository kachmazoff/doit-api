package controller

func createMessage(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
	}
}
