package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	filterPtr := flag.String("filter", ".txt", "file extension to filter")
	outputFileNamePtr := flag.String("fileName", "output.csv", "Output file name")

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	type userdata struct {
		User  string
		ID    string
		Price string
	}

	stats := make([]userdata, 0, 3)

	for _, fileInfo := range files {

		fileName := fileInfo.Name()
		if !strings.Contains(fileName, *filterPtr) {
			continue
		}

		fmt.Println("Processing", fileName)

		file, err := os.Open(fileName)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()

			userExp, _ := regexp.Compile("user ([a-z]+) ")
			idExp, _ := regexp.Compile(" id=([0-9]+),")
			priceExp, _ := regexp.Compile(" price=([0-9]+\\.?[0-9]+)")

			userName := userExp.FindStringSubmatch(text)
			id := idExp.FindStringSubmatch(text)
			price := priceExp.FindStringSubmatch(text)

			userStats := userdata{
				User:  userName[1],
				ID:    id[1],
				Price: price[1],
			}

			stats = append(stats, userStats)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(stats)

	csvRecords := [][]string{
		{"user", "id", "price"},
	}

	for _, data := range stats {
		csvRecords = append(csvRecords, []string{data.User, data.ID, data.Price})
	}

	outputFile, err1 := os.Create(*outputFileNamePtr)
	if err1 != nil {
		log.Fatal("Cannot create file", err1)
	}

	w := csv.NewWriter(outputFile)

	for _, record := range csvRecords {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
