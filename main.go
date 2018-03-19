package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"
)

func main() {
	rest := false
	nocount := false
	reset := false
	doneperiods := ""
	bonuses := ""
	n, a := 25, 25
	t := 5
	c := n * 60
	tick := time.Tick(time.Duration(t) * 60 * 1000 * time.Millisecond)
	boom := time.After(time.Duration(a) * 60 * 1000 * time.Millisecond)

	if len(os.Args) > 1 {
		if os.Args[1] == "rest" {
			rest = true
		} else if os.Args[1] == "nocount" {
			nocount = true
		} else if os.Args[1] == "reset" {
			reset = true
		}
	}

	// If the file doesn't exist, create it, or append to the file
	usr, _ := user.Current()
	dir := usr.HomeDir
	path := fmt.Sprintf("%s/gopomidoro.log", dir)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	b1 := make([]byte, 10)
	_, err = f.Read(b1)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	reslist := strings.Split(string(b1), " ")
	if len(reslist) > 1 {
		if reset {
			if _, err := f.WriteAt([]byte("0 done 0 "), 0); err != nil {
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
			return
		} else {
			doneperiods = reslist[0]
			bonuses = reslist[2]
		}
	} else {
		doneperiods = "0"
		bonuses = "0"
		if _, err := f.WriteAt([]byte("0 done 0 "), 0); err != nil {
			log.Fatal(err)
		}
	}

	dc, err := strconv.Atoi(doneperiods)
	if err != nil {
		fmt.Println(err)
	}
	bc, err := strconv.Atoi(bonuses)
	if err != nil {
		fmt.Println(err)
	}

	if rest && bc == 0 {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		panic("You dont have bonuses to rest, work more!")
	}

	for {
		select {
		case <-tick:
			n = n - 5
			fmt.Println("tick.")
			str := fmt.Sprintf("залишилося - %d хв.", n)
			u := ""

			if n == 5 {
				u = "critical"
			} else {
				u = "low"
			}

			cmd := exec.Command("/usr/bin/notify-send", "-u", u, str)
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		case <-boom:
			fmt.Println("BOOM!")
			cmd := exec.Command("/usr/bin/notify-send", "Час вийшов, відпочинь")
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%d\n", bc)
			if !nocount {
				if !rest {
					dc = dc + 1
					if dc%2 == 0 {
						bc = bc + 1
					}
				} else {
					bc = bc - 1
				}
			}
			msg := fmt.Sprintf("%d done %d ", dc, bc)
			if _, err := f.WriteAt([]byte(msg), 0); err != nil {
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
			return
		default:
			res := ""
			c = c - 1
			str := fmt.Sprintf("залишилося - %d сек.", c)
			fmt.Println(str)
			if c < 100 {
				res = fmt.Sprintf("%s   %d %s ", doneperiods, c, bonuses)
			} else if c < 1000 {
				res = fmt.Sprintf("%s  %d %s ", doneperiods, c, bonuses)
			} else {
				res = fmt.Sprintf("%s %d %s ", doneperiods, c, bonuses)
			}
			if _, err := f.WriteAt([]byte(res), 0); err != nil {
				log.Fatal(err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
