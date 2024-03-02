package file

import (
	"fmt"
	"log"
	"os"
	"time"
)

type History struct {
	Name  string
	Url   string
	Start *time.Time
	End   *time.Time
}

func CreateFile(histores []History) {
	fmt.Printf("=>")
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Panic(err.Error())
	}

	fmt.Println("entries : ", entries)

	now := time.Now().Format("2006-01-02_15:04:05")

	fmt.Println("time : ", now)
	f, err := os.Create(fmt.Sprintf("./output/%s.report", now))

	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	for _, history := range histores {
		_, err := f.WriteString(fmt.Sprintf("name : %s\n url :%s  \n %s | %s\n", history.Name, history.Url, history.Start.Format("2006-01-02 15:04:05"), GetElapsedTime_ns(history)))
		if err != nil {
			log.Fatal(err)
		}

	}
}

func GetElapsedTime_ns(history History) string {
	start := history.Start
	end := history.End
	elapsed_ns := time.Duration(((*end).Unix() - (*start).Unix()) * 1000000000)
	return elapsed_ns.String()
}
