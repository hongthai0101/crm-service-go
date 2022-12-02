package pkg

func GetLimit(limit int64) int64 {
	if limit == 0 {
		return 10
	}
	return limit
}

func GetSkip(skip int64) int64 {
	if skip == 0 {
		return 0
	}
	return skip
}
