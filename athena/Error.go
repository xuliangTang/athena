package athena

func Error(err error, msg ...string) {
	if err != nil {
		errMsg := err.Error()
		if len(msg) > 0 {
			errMsg = msg[0]
		}
		panic(errMsg)
	}
}

func Unwrap(result any, err error) any {
	if err != nil {
		panic(err.Error())
	}

	return result
}

func UnwrapOrEmpty(result string, err error) string {
	if err != nil {
		return ""
	}

	return result
}
