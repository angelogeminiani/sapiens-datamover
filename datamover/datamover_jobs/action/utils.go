package action

func IsRecordNotFoundError(err error) bool {
	return nil != err && err.Error() == "record not found"
}
