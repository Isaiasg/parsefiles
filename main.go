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

	keys := []string{"User", "ID", "Price"}

	csvStructure := map[string]string{
		"User":  "user ([a-z]+) ",
		"ID":    " id=([0-9]+),",
		"Price": " price=([0-9]+\\.?[0-9]+)",
	}

	regexpCollection := make(map[string]*regexp.Regexp, len(keys))
	for _, key := range keys {
		fmt.Println(key, csvStructure[key])
		compiledRegexp, _ := regexp.Compile(csvStructure[key])
		regexpCollection[key] = compiledRegexp
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	csvRecords := [][]string{}
	csvRecords = append(csvRecords, keys)

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

			line := make([]string, len(keys))
			for index, key := range keys {
				line[index] = regexpCollection[key].FindStringSubmatch(text)[1]
			}

			csvRecords = append(csvRecords, line)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	createOutputFile(csvRecords, *outputFileNamePtr)
}

func createOutputFile(records [][]string, fileName string) {
	outputFile, err1 := os.Create(fileName)
	if err1 != nil {
		log.Fatal("Cannot create file", err1)
	}

	w := csv.NewWriter(outputFile)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
