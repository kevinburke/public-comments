package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var bayCounties = map[string]bool{
	"San Mateo County":     true,
	"Contra Costa County":  true,
	"San Francisco County": true,
	"Santa Clara County":   true,
	"Alameda County":       true,
	"Marin County":         true,
}

func main() {
	f, err := os.Open("county-to-county.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	bs := bufio.NewScanner(f)
	in := make(map[string]map[int]int)
	out := make(map[string]map[int]int)
	count := 0
	//in := make(map[int64]int)
	//out := make(map[int64]int)
	// parsing instructions found here:
	// https://www.census.gov/content/dam/Census/topics/population/migration/guidance-for-data-users/acs-migration-tutorial/2011-2015%20Migration%20Flows%20Documentation.pdf
	// note field names in the document are 1-indexed.
	for bs.Scan() {
		count++
		text := bs.Text()
		if text == "" {
			continue
		}
		if len(text) < 15 {
			log.Fatalf("too short row: %q\n", text)
		}
		currentCounty := strings.TrimSpace(text[45:79])
		otherCounty := strings.TrimSpace(text[225:259])
		if otherCounty == "-" || otherCounty == "" {
			continue
		}
		if bayCounties[currentCounty] == false && bayCounties[otherCounty] == false {
			continue
		}
		if _, ok := in[currentCounty]; bayCounties[currentCounty] && !ok {
			in[currentCounty] = make(map[int]int)
		}
		if _, ok := out[currentCounty]; bayCounties[currentCounty] && !ok {
			out[currentCounty] = make(map[int]int)
		}
		if _, ok := in[otherCounty]; bayCounties[otherCounty] && !ok {
			in[otherCounty] = make(map[int]int)
		}
		if _, ok := out[otherCounty]; bayCounties[otherCounty] && !ok {
			out[otherCounty] = make(map[int]int)
		}
		code64, err := strconv.ParseInt(text[13:15], 10, 8)
		if err != nil {
			log.Fatal(err)
		}
		code := int(code64)
		if code <= 0 || code > 4 {
			log.Fatal("invalid race code", code)
		}
		if len(text) < 382 {
			log.Printf("row does not contain net flow: %q", text)
			continue
		}
		flow, err := strconv.ParseInt(strings.TrimSpace(text[375:382]), 10, 64)
		if err != nil {
			log.Fatalf("error parsing %q: %v", text, err)
		}
		if bayCounties[currentCounty] {
			// FROM other county TO current county
			if tot, ok := in[currentCounty][code]; ok {
				in[currentCounty][code] = tot + int(flow)
			} else {
				in[currentCounty][code] = int(flow)
			}
		}
		if bayCounties[otherCounty] {
			// FROM current county TO other county
			if tot, ok := out[otherCounty][code]; ok {
				out[otherCounty][code] = tot + int(flow)
			} else {
				out[otherCounty][code] = int(flow)
			}
		}
		//diffCountyForCurrentCounty, err := strconv.ParseInt(strings.TrimSpace(text[147:153]), 10, 64)
		//if err != nil {
		//log.Fatal(err)
		//}
		//diffStateForCurrentCounty, err := strconv.ParseInt(strings.TrimSpace(text[163:169]), 10, 64)
		//if err != nil {
		//log.Fatal(err)
		//}
		//currentCountyToDiffState, err := strconv.ParseInt(strings.TrimSpace(text[163:169]), 10, 64)
		//if err != nil {
		//log.Fatal(err)
		//}
		//fmt.Println("code", code, "diff county for same county", diffCountyForCurrentCounty, "diff state", diffStateForCurrentCounty)
	}
	if err := bs.Err(); err != nil {
		log.Fatal(err)
	}
	counties := make([]string, len(bayCounties))
	i := 0
	for county := range bayCounties {
		counties[i] = county
		i++
	}
	sort.Strings(counties)
	fmt.Printf("scanned %d rows\n", count)
	for i := range counties {
		totalMigrants := 0
		for j := 1; j <= 4; j++ {
			totalMigrants += in[counties[i]][j]
		}
		fmt.Println(counties[i], "Migration In:")
		fmt.Printf("White: %d (%.1f%%)\n", in[counties[i]][1], float64(100*in[counties[i]][1])/float64(totalMigrants))
		fmt.Printf("Black: %d (%.1f%%)\n", in[counties[i]][2], float64(100*in[counties[i]][2])/float64(totalMigrants))
		fmt.Printf("Asian: %d (%.1f%%)\n", in[counties[i]][3], float64(100*in[counties[i]][3])/float64(totalMigrants))
		fmt.Printf("Other: %d (%.1f%%)\n", in[counties[i]][4], float64(100*in[counties[i]][4])/float64(totalMigrants))

		totalMigrants = 0
		for j := 1; j <= 4; j++ {
			totalMigrants += out[counties[i]][j]
		}
		fmt.Println(counties[i], "Migration Out:")
		fmt.Printf("White: %d (%.1f%%)\n", out[counties[i]][1], float64(100*out[counties[i]][1])/float64(totalMigrants))
		fmt.Printf("Black: %d (%.1f%%)\n", out[counties[i]][2], float64(100*out[counties[i]][2])/float64(totalMigrants))
		fmt.Printf("Asian: %d (%.1f%%)\n", out[counties[i]][3], float64(100*out[counties[i]][3])/float64(totalMigrants))
		fmt.Printf("Other: %d (%.1f%%)\n", out[counties[i]][4], float64(100*out[counties[i]][4])/float64(totalMigrants))
		fmt.Println("")
	}
}
