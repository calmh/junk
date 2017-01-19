package main

import (
	"bufio"
	"fmt"
	"io"
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

	var accounts map[int]int
	var currentMonth int

	fmt.Println("Month,Category,Value")

	for _, ver := range vers {
		month := ver.Date / 100
		if month != currentMonth {
			if currentMonth != 0 {
				getSummary(currentMonth, accounts).WriteTo(os.Stdout)
			}
			accounts = make(map[int]int)
			currentMonth = month
		}
		for _, tran := range ver.Trans {
			accounts[tran.Account] += tran.Amount
		}
	}

	getSummary(currentMonth, accounts).WriteTo(os.Stdout)
}

type summary struct {
	month     string
	income    int
	expenses  int
	employees int
}

func (s summary) WriteTo(w io.Writer) error {
	if _, err := fmt.Fprintf(w, `"%s","Income",%d`+"\n", s.month, s.income); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, `"%s","Expenses",%d`+"\n", s.month, s.expenses); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, `"%s","Employees",%d`+"\n", s.month, s.employees); err != nil {
		return err
	}
	return nil
}

func getSummary(month int, m map[int]int) (s summary) {
	s.month = fmt.Sprintf("%04d-%02d", month/100, month%100)
	for account, value := range m {
		switch account / 1000 {
		case 3:
			s.income -= value
		case 5, 6:
			s.expenses += value
		case 7:
			s.employees += value
		}
	}
	return
}
