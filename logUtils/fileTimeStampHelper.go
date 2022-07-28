package logUtils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func isOdysseyFile(fileName string) bool {
	// Odyssey Example
	// Journal.2022-05-20T170535.01.log
	match, err := regexp.MatchString("\\bJournal\\.\\d{4}-\\d{2}-\\d{2}T\\d*\\.\\d*\\.log\\b", fileName)
	if err != nil {
		panic(err)
	}
	return match
}

func isHorizonsFile(fileName string) bool {
	// Horizons Example
	// Journal.220212014914.01.log
	match, err := regexp.MatchString("\\bJournal\\.\\d*\\.\\d*.log\\b", fileName)
	if err != nil {
		panic(err)
	}
	return match
}

func convertOddyToHorizonsName(name string) string {
	// if already horizons, return as is
	if isHorizonsFile(name) {
		return name
	}

	// 2022-05-20T170535
	// YYYY-MM-DDThhmmss
	splitString := strings.Split(name, ".")
	timeStampOddy := splitString[1]

	// 2022-05-20  170535
	dateAndTimeSplit := strings.Split(timeStampOddy, "T")
	//2022  05  20
	dates := strings.Split(dateAndTimeSplit[0], "-")
	// on the first one, remove the first two chars. turn 2022 into 22
	dates[0] = dates[0][2:]

	asHorizonsTimeStamp := strings.Join(dates, "") + dateAndTimeSplit[1]
	splitString[1] = asHorizonsTimeStamp
	return strings.Join(splitString, ".")
}


func getTimeStamp(name string) int64 {

	timeCodePart := strings.Split(name, ".")[1]
	year, _ := strconv.Atoi("20" + timeCodePart[0:2])
	month, _ := strconv.Atoi(timeCodePart[2:4])
	day, _ := strconv.Atoi(timeCodePart[4:6])
	hr, _ := strconv.Atoi(timeCodePart[6:8])
	min, _ := strconv.Atoi(timeCodePart[8:10])
	sec, _ := strconv.Atoi(timeCodePart[10:12])

	date := time.Date(year, time.Month(month), day, hr, min, sec, 0, time.UTC)

	return date.Unix()

}

func GetLogFileTimestamp(fileName string) int64 {
	if isOdysseyFile(fileName) {
		fileName = convertOddyToHorizonsName(fileName)
	}

	return getTimeStamp(fileName)
}
