package frontend

func getErrorTitle(statusCode int) string {
	switch statusCode {
	default:
		return "UNKNOWN"
	case 404:
		return "NOT FOUND"
	case 500:
		return "Internal error"
	}
}
