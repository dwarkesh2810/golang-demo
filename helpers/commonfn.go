package helpers

func SecondsToDayHourMinAndSeconds(seconds int) (int64, int64, int64, int64) {
	days := seconds / 86400
	hour := (seconds % 86400) / 3600
	minute := (seconds % 3600) / 60
	second := seconds % 60
	return int64(days), int64(hour), int64(minute), int64(second)
}
