package types

type ServiceErrors map[string]string

func (s ServiceErrors) Error() string {
	return s["error"]
}

func SvcError(msg string, errs ...error) ServiceErrors {
	if len(errs) == 0 {
		return ServiceErrors{"error": msg}
	}

	return ServiceErrors{"error": msg, "details": errs[0].Error()}
}
