package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func help() {
	fmt.Println("WhatsApp group stats extractor")
	fmt.Println("Usage:")
	fmt.Printf("%s <<WhatsApp export file>>\n", os.Args[0])
}

func getstats(filename string) (map[string]int, map[string]int, string, string) {
	validline := regexp.MustCompile(`^\d.*\s-\s.*:`)
	date := regexp.MustCompile(`^\d.-\d.-\d.`)
	user := regexp.MustCompile(`-\s(.*?):`)
	datedb := make(map[string]int)
	userdb := make(map[string]int)
	bytesfile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	filecontent := string(bytesfile)
	lines := strings.Split(filecontent, "\n")
	for _, line := range lines {
		if validline.MatchString(line) {
			date := date.FindStringSubmatch(line)[0]
			username := user.FindStringSubmatch(line)[1]
			userdb[username]++
			datedb[date]++
		}
	}
	startdate := date.FindStringSubmatch(lines[0])[0]
	enddate := date.FindStringSubmatch(lines[len(lines)-2])[0]
	return userdb, datedb, startdate, enddate
}

func avgpostperday(mapname map[string]int) int {
	avg := 0
	for _, v := range mapname {
		avg = avg + v
	}
	avg = avg / len(mapname)
	return avg
}

func totalposts(mapname map[string]int) int {
	total := 0
	for _, v := range mapname {
		total = total + v
	}
	return total
}

type datevalue struct {
	date  string
	count int
}

func top3days(mapname map[string]int) []datevalue {
	tmp := make([]datevalue, 3)
	for i := 0; i < 3; i++ {
		tmpholder := ""
		for k, v := range mapname {
			if v > mapname[tmpholder] {
				tmpholder = k
			}
		}
		tmp[i].date = tmpholder
		tmp[i].count = mapname[tmpholder]
		delete(mapname, tmpholder)
	}
	return tmp
}

func main() {
	if len(os.Args) == 1 {
		help()
	} else {
		userdb, datedb, startdate, enddate := getstats(os.Args[1])
		fmt.Printf("Analyzing: \t%s\n", os.Args[1])
		fmt.Printf("Start date: \t%s\n", startdate)
		fmt.Printf("End date: \t%s\n", enddate)
		fmt.Printf("Total msg: \t%d\n", totalposts(userdb))
		fmt.Printf("Avg per day: \t%d\n", avgpostperday(datedb))
		fmt.Println("Total messages per user:")
		for k, v := range userdb {
			fmt.Println("  ", k, ":\t", v)
		}
		fmt.Println("Average messages per user/day:")
		for k, v := range userdb {
			fmt.Println("  ", k, ":\t", v/len(datedb))
		}
		fmt.Println("Top 3 days:")
		for _, record := range top3days(datedb) {
			fmt.Printf("   %s :\t %d\n", record.date, record.count)
		}
	}
}
