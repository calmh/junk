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
	Id    int
	Date  int
	Descr string
	Trans []Trans
}

type Trans struct {
	Account int
	Amount  float64
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
		return s[i].Id < s[j].Id
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
	fmt.Println(fname)

	inf, e := os.Open(fname)
	if e != nil {
		panic(e)
	}

	convInf, e := iconv.NewReader(inf, "cp850", "utf-8")
	if e != nil {
		panic(e)
	}
	scanner := bufio.NewScanner(convInf)

	outf, e := os.Create(fname + ".clean.se")
	if e != nil {
		panic(e)
	}

	convOutf, e := iconv.NewWriter(outf, "utf-8", "cp850")
	if e != nil {
		panic(e)
	}

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
			curVer = Ver{Type: typ, Id: id, Date: date, Descr: descr}
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
			trans := Trans{Account: account, Amount: amount}
			curVer.Trans = append(curVer.Trans, trans)
		} else if line == "{" || strings.HasPrefix(line, "#RTRANS") || strings.HasPrefix(line, "#BTRANS") {
			// Nothing
		} else if line == "}" {
			vers = append(vers, curVer)
		} else {
			fmt.Fprintln(convOutf, line)
		}
	}

	sort.Sort(vers)

	for i, ver := range vers {
		if len(ver.Trans) > 0 {
			fmt.Fprintf(convOutf, "#VER A %d %d %s\n", i+1, ver.Date, ver.Descr)
			fmt.Fprintln(convOutf, "{")
			for _, tran := range ver.Trans {
				fmt.Fprintf(convOutf, "#TRANS %d {} %.02f\n", tran.Account, tran.Amount)
			}
			fmt.Fprintln(convOutf, "}")
		}
	}
}
