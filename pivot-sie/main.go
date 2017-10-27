package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"log"

	"github.com/djimenez/iconv-go"
)

type Ver struct {
	Type  string
	ID    int
	Date  int
	Descr string
	Trans []Trans
}

type Trans struct {
	Account int
	Amount  int
}

type VerList []Ver

func (s VerList) Len() int {
	return len(s)
}

func (s VerList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s VerList) Less(i, j int) bool {
	if s[i].Date == s[j].Date {
		return s[i].ID < s[j].ID
	}
	return s[i].Date < s[j].Date
}

func main() {
	if len(os.Args) != 2 {
		println("clean-sie -- Reorder and clean SIE files")
		println("usage: clean-sie <input-file>")
		os.Exit(2)
	}

	fname := os.Args[1]

	inf, e := os.Open(fname)
	if e != nil {
		panic(e)
	}

	convInf, e := iconv.NewReader(inf, "cp850", "utf-8")
	if e != nil {
		panic(e)
	}
	scanner := bufio.NewScanner(convInf)

	var vers VerList
	var curVer Ver

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#VER") {
			fields := strings.SplitN(line, " ", 5)
			typ := fields[1]
			id, err := strconv.Atoi(fields[2])
			if err != nil {
				fmt.Println(line)
				log.Fatal(err)
			}
			date, err := strconv.Atoi(fields[3])
			if err != nil {
				fmt.Println(line)
				log.Fatal(err)
			}
			descr := fields[4]
			curVer = Ver{Type: typ, ID: id, Date: date, Descr: descr}
		} else if strings.HasPrefix(line, "#TRANS") {
			fields := strings.SplitN(line, " ", -1)
			account, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println(line)
				log.Fatal(err)
			}
			amount, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				fmt.Println(line)
				log.Fatal(err)
			}
			trans := Trans{Account: account, Amount: int(amount * 100)}
			curVer.Trans = append(curVer.Trans, trans)
		} else if line == "{" || strings.HasPrefix(line, "#RTRANS") || strings.HasPrefix(line, "#BTRANS") {
			// Nothing
		} else if line == "}" {
			vers = append(vers, curVer)
		}
	}

	sort.Sort(vers)

	accounts := make(map[int]map[int]int) // month => account => value
	accountUsed := make(map[int]bool)

	for _, ver := range vers {
		month := ver.Date / 100
		accs, ok := accounts[month]
		if !ok {
			accs = make(map[int]int)
			accounts[month] = accs
		}
		for _, tran := range ver.Trans {
			accs[tran.Account] += tran.Amount
			accountUsed[tran.Account] = true
		}
	}

	var months []int
	var accountIDs []int
	for month := range accounts {
		months = append(months, month)
	}
	for acc := range accountUsed {
		if acc/1000 < 3 {
			continue
		}
		accountIDs = append(accountIDs, acc)
	}
	sort.Ints(months)
	sort.Ints(accountIDs)

	fmt.Printf(`"Month"`)
	for _, acc := range accountIDs {
		fmt.Printf(`;"%d"`, acc)
	}
	fmt.Println()

	for _, month := range months {
		fmt.Printf(`"%d-%d"`, month/100, month%100)
		for _, acc := range accountIDs {
			v := fmt.Sprintf("%.02f", float64(accounts[month][acc])/100)
			fmt.Printf(";%s", strings.Replace(v, ".", ",", 1))
		}
		fmt.Println()
	}
}
