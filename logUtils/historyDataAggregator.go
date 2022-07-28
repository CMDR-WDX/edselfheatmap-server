package logUtils

import (
	"EDSelfHeatmap/data"
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func getRelevantFiles(journalDir string) []string {
	relevantFiles := make([]string, 0)

	files, err := ioutil.ReadDir(journalDir)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if strings.HasPrefix(fileName, "Journal.") && strings.HasSuffix(fileName, ".log") {
			relevantFiles = append(relevantFiles, path.Join(journalDir, file.Name()))
		}
	}
	return relevantFiles
}

type journalJumpEvent struct {
	Event      string    `json:"event,omitempty"`
	StarSystem string    `json:"StarSystem,omitempty"`
	StarPos    []float32 `json:"StarPos,omitempty"`
}

func extractJumpEventsFromJournalFile(filepath string) []data.RequestBody {
	openedFile, err := os.Open(filepath)
	defer openedFile.Close()
	if err != nil {
		log.Println("[ERR] Failed to parse file")
		log.Println(err)
		return make([]data.RequestBody, 0)
	}
	fScanner := bufio.NewScanner(openedFile)
	fScanner.Split(bufio.ScanLines)

	returnData := make([]data.RequestBody, 0)

	for fScanner.Scan() {
		line := fScanner.Text()
		if !strings.Contains(line, "\"event\":\"FSDJump\"") {
			continue
		}

		var asJson journalJumpEvent
		err := json.Unmarshal([]byte(line), &asJson)
		if err != nil {
			log.Println("[ERR] Failed to parse entry. Skipping")
			log.Println(err)
		}

		x := asJson.StarPos[0]
		y := asJson.StarPos[1]
		z := asJson.StarPos[2]

		reqBody := data.RequestBody{
			SystemName: &asJson.StarSystem,
			X:          &x,
			Y:          &y,
			Z:          &z,
		}

		returnData = append(returnData, reqBody)

	}
	return returnData
}

func GetJournalHistoryJumpEvents() []data.RequestBody {
	journalDir := strings.Trim(os.Getenv("LOG_DIR"), " ")

	val, err := os.Stat(journalDir)
	if err != nil {
		log.Println("[WARNING] Skipped journal history aggregation because of an error")
		log.Println(err)
		return make([]data.RequestBody, 0)
	}

	if !val.IsDir() {
		log.Println("[WARNING] Skipped journal history aggregation because provided path is a file, not a directory")
		return make([]data.RequestBody, 0)
	}

	relevantFiles := getRelevantFiles(journalDir)
	fileCount := len(relevantFiles)
	log.Printf("Found %d relevant files.\n", fileCount)
	jumpEvents := make([]data.RequestBody, 0)
	for i, file := range relevantFiles {
		log.Printf("Parsing file %d / %d", i+1, fileCount)
		events := extractJumpEventsFromJournalFile(file)
		jumpEvents = append(jumpEvents, events...)
	}
	return jumpEvents

}
