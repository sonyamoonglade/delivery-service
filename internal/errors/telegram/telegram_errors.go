package tg_errors

type TelegramError struct {
	Message string
	Err     error
}

func (e TelegramError) Error() string {
	return e.Message
}
