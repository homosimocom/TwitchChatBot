package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/sqweek/dialog"
)

var fbot string = ""
var vmass []string
var vchat []string

func main() {
	var flag = 1

	var wg sync.WaitGroup
	app := app.New()
	gname := &widget.Entry{}
	gname = widget.NewEntry()

	lcount := &widget.Entry{}
	lcount = widget.NewEntry()
	cdelay := &widget.Entry{}
	cdelay = widget.NewEntry()
	gbot := &widget.Button{}
	gchat := &widget.Button{}
	connbots := &widget.Button{}
	startchat := &widget.Button{}
	stopchat := &widget.Button{}

	gbot = widget.NewButton("select a list of bots(wait 1-20 sec)", func() {
		filename, err := dialog.File().Filter("txt", "txt").Load()
		if filename == "" {
			return
		}
		_check(err)
		file, err := os.Open(filename)
		_check(err)
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		fmt.Print(b)
		vmass = []string{}
		for _, value := range strings.Split(string(b), "\n") {
			fmt.Println(value)
			if len(value) > 0 {
				vmass = append(vmass, value)
			}

		}
		if len(vmass) > 0 {
			gchat.Show()
			connbots.Show()
			//	startchat.Show()
		} else {
			gchat.Hide()
			connbots.Hide()
			//	startchat.Hide()
		}
	})
	gchat = widget.NewButton("select a list of chat", func() {
		filename, err := dialog.File().Filter("txt", "txt").Load()
		if filename == "" {
			return
		}

		fmt.Println(filename)
		_check(err)

		file, err := os.Open(filename)
		if err != nil {

			log.Fatal(err)
		}
		defer file.Close()
		if filename == "" {
			return
		}
		b, err := ioutil.ReadAll(file)
		fmt.Print(b)
		vchat = []string{}
		for _, value := range strings.Split(string(b), "\n") {
			fmt.Println(value)
			if len(value) > 0 {
				vchat = append(vchat, value)
			}

		}
		if len(vchat) > 0 {
			startchat.Show()
		} else {
			startchat.Hide()
		}
	})
	connbots = widget.NewButton("Connect chat list", func() {
		if len(gname.Text) == 0 {
			return
		}
		count, err := strconv.Atoi(lcount.Text)
		if err != nil {
			lcount.SetText("")
			fmt.Println("error bots maximum")
			return
		}
		if len(vmass) == 0 {
			return
		}
		flag = 1
		go func(room string, flag *int) {
			//count := 10
			defer wg.Done()
			if count > len(vmass)-1 {
				count = len(vmass) - 1
			}
			for _, value := range vmass[:count] {
				time.Sleep(100 * time.Millisecond)
				wg.Add(1)
				go func(value string, room string, flag *int) {
					defer wg.Done()
					conchat(value, room, flag)
				}(value, room, flag)
			}
		}(gname.Text, &flag)
		fmt.Println("Connect bots process...")

	})
	startchat = widget.NewButton("Start chat bots", func() {
		if len(gname.Text) == 0 {
			return
		}
		count, err := strconv.Atoi(lcount.Text)
		if err != nil {
			lcount.SetText("")
			fmt.Println("error bots maximum")
			return
		}
		delay, err := strconv.Atoi(cdelay.Text)
		if err != nil {
			cdelay.SetText("")
			fmt.Println("error bots delay")
			return
		}
		if delay == 0 {
			return
		}
		if len(vmass) == 0 {
			return
		}
		if len(vchat) == 0 {
			return
		}
		flag = 1
		fmt.Println(count, delay)
		startchat.Hide()
		//stopchat.Show()
		gname.ReadOnly = true
		fmt.Println("Connect bots process...")
		wg.Add(1)
		go func(delay int, room string, count int, flag *int) {
			defer wg.Done()
			flagchat(delay, room, count, flag)
		}(delay, gname.Text, count, &flag)

	})

	stopchat = widget.NewButton("Stop all bots", func() {
		startchat.Show()
		//stopchat.Hide()
		gname.ReadOnly = false
		flag = 0
		fmt.Println("Connect bots process...")
	})
	startchat.Show()
	//stopchat.Hide()
	gname.SetPlaceHolder("Enter channel name")
	lcount.SetPlaceHolder("Enter bots count")
	cdelay.SetPlaceHolder("Enter delay for bots")
	w := app.NewWindow("LetsWin")
	qq := widget.NewVBox(widget.NewVBox(
		widget.NewLabel("LetsWin - beilus chat bots"),
		gname,
		lcount,
		cdelay,

		gbot, gchat, connbots, startchat, stopchat,

		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))
	bbb := widget.NewEntry()
	bbb.SetText(GetKey() + " Need Activate @flexcat telegram")
	bad := widget.NewVBox(widget.NewVBox(
		bbb,
	))

	qq.Show()
	bbb.Hide()
	gchat.Hide()
	connbots.Hide()
	startchat.Hide()
	w.SetContent(widget.NewVBox(
		qq, bad,
	))

	w.ShowAndRun()
}
