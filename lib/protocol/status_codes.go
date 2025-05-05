package protocol

const (
	STATUS_OK = iota
	STATUS_NOT_FOUND
	STATUS_BAD_REQUEST
	STATUS_INTERNAL_SERVER_ERROR
)

func TranslateStatus(status int) string {
	switch status {
	case STATUS_OK:
		return "OK"
	case STATUS_NOT_FOUND:
		return "Not Found"
	case STATUS_BAD_REQUEST:
		return "Bad Request"
	case STATUS_INTERNAL_SERVER_ERROR:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}
