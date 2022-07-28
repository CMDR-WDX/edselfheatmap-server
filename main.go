package main

import (
	data2 "EDSelfHeatmap/data"
	"EDSelfHeatmap/database"
	"EDSelfHeatmap/img"
	"EDSelfHeatmap/logUtils"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to find .env-File. This isnt necessarily an issue if you pass the Env Args on startup.")
	} else {
		log.Println(".env Found")
	}

	PortAddr, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Printf("Failed to parse Port. Either none provided, or wrong. defaulting to 8080. Reason: %s\n", err)
		PortAddr = 8080
	}
	log.Printf("Port is %d\n", PortAddr)


	shouldDoFullLogScan := database.Init()
	defer database.DB.Close()


	if shouldDoFullLogScan {
		log.Println("FIRST TIME RUN: GETTING HISTORICAL DATA")
		jumpEvents := logUtils.GetJournalHistoryJumpEvents()
		count := len(jumpEvents)
		for i, entry := range jumpEvents {
			if i % 10 == 0 {
				log.Printf("Placing in Database: %d / %d", i+1, count)
			}
			pxEntry := data2.MakePixelEntry(entry)
			err := database.NotifyAboutNewSystemNoEmit(pxEntry)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	initState, err := database.GetDBPixelMap()
	if err != nil {
		log.Fatalln(err)
	}

	img.InitFromEnv(initState)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "image/png")
			imgResult := img.MakeImage()
			w.Write(imgResult.Bytes())
		}
		if r.Method == "POST" {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("ERR: ", err)
				w.WriteHeader(400)
				fmt.Fprintf(w, err.Error())
			}
			var asJson data2.RequestBody
			err = json.Unmarshal(data, &asJson)
			if err != nil {
				log.Println("ERR: ", err)
				w.WriteHeader(400)
				fmt.Fprintf(w, err.Error())
			}
			if asJson.Z == nil || asJson.Y == nil || asJson.X == nil || asJson.SystemName == nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "Wrong Structure. Json with x,y,z and systemName required")
			}

			// Convert to Pixel info
			asPixel := data2.MakePixelEntry(asJson)

			database.NotifyAboutNewSystem(asPixel)

		}
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PortAddr), nil))
}
