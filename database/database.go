package database

import (
	"EDSelfHeatmap/data"
	"EDSelfHeatmap/img"
	"database/sql"
	"errors"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"os"
)

var DB *sql.DB

func Init() bool {
	parseLogHistory := false
	_, err := os.Stat("db.sqlite")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			parseLogHistory = true
		} else {
			log.Fatalln(err)
		}
	}

	d, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	DB = d
	res, err := d.Exec("CREATE TABLE IF NOT EXISTS systems (system_name varchar PRIMARY KEY, x integer NOT NULL, y integer NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	return parseLogHistory

}

func notifyAboutNewSystemNoEmit(data data.PixelEntry) (error, bool) {
	// Check if system already in DB
	exec, err := DB.Query("select count (*) from systems where system_name = ?", data.SystemName)
	if err != nil {
		log.Println(err)
		return err, false
	}
	var count int

	for exec.Next() {
		err = exec.Scan(&count)
		if err != nil {
			log.Println(err)
			exec.Close()
			return err, false
		}
	}
	exec.Close()
	if count == 0 {
		// Put in new System
		_, err := DB.Exec("insert into systems (system_name, x, y) values (?, ?, ?)", data.SystemName, data.X, data.Y)
		if err != nil {
			log.Println(err)
			return err, false
		}
		return nil, true

	} else {
		// Already present
		return nil, false
	}
	return nil, false
}

func NotifyAboutNewSystemNoEmit(data data.PixelEntry) error {
	err, _ := notifyAboutNewSystemNoEmit(data)
	return err
}

func NotifyAboutNewSystem(data data.PixelEntry) error {

	err, isNew := notifyAboutNewSystemNoEmit(data)
	if err != nil {
		return err
	}
	if isNew {
		img.Increment(data.X, data.Y)
	}
	return nil
}

func GetDBPixelMap() (map[data.IntPoint]int, error) {
	res, err := DB.Query("SELECT * FROM systems")
	if err != nil {
		log.Fatal(err)
	}

	lookup := make(map[data.IntPoint]int)

	{
		defer res.Close()
		for res.Next() {
			var name string
			var x int
			var y int
			err = res.Scan(&name, &x, &y)
			if err != nil {
				log.Println(err)
				return map[data.IntPoint]int{}, err
			}
			asPoint := data.IntPoint{x, y}
			value, ok := lookup[asPoint]
			if !ok {
				value = 0
			}
			value += 1
			lookup[asPoint] = value
		}
	}

	return lookup, nil
}
