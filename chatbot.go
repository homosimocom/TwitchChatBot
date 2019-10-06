package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/denisbrodbeck/machineid"
	"github.com/gempir/go-twitch-irc"
)

func conchat(value string, room string, flag *int) {
	data := strings.Split(value, ":")
	client := twitch.NewClient(data[0], "oauth:"+data[2])
	for {
		if *flag == 0 {
			return
		}
		client.Join(room)
		client.OnConnect(func() {
			fmt.Println("Connected ", data)
		})
		err := client.Connect()
		if err != nil {

			fmt.Println("i`m die...", err)
			return
		}
		if *flag == 0 {
			return
		}
		time.Sleep(1800000 * time.Millisecond)
	}
}

func _check(err error) {

	if err != nil {
		panic(err)
	}

}

func flagchat(ccc int, room string, count int, flag *int) {
	var msg []string
	var wg sync.WaitGroup

	msg = vchat

	r := 0
	r2 := 0
	r3 := 0
	if count == 0 {
		count = 1
	}

	if len(msg) > 0 {
		for {
			cc := ccc + randInt(0, 2)
			r3 = randInt(0, 30)
			if r3 < 5 {
				cc = randInt(0, 10)
			}
			r2 = randInt(0, count)
			if *flag == 0 {
				return
			}
			value := vmass[r2]
			if len(msg) == 1 {
				r = 0
			} else {
				r = randInt(0, len(msg)-1)
			}
			rmsg := msg[r]
			if *flag == 0 {
				return
			}
			wg.Add(1)
			time.Sleep(time.Duration(cc*1000) * time.Millisecond)
			go func(value string, room string, rmsg string) {
				defer wg.Done()
				chat(room, value, rmsg)
			}(value, room, rmsg)

		}
	} else {
		return
	}

}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func chat(room string, data string, msg string) {
	if len(data) < 1 {
		fmt.Println(room, data, msg)
		return
	}
	tdata := strings.Split(data, ":")
	client := twitch.NewClient(tdata[0], "oauth:"+tdata[2])
	client.Join(room)
	client.OnConnect(func() {
		client.Say(room, msg)
		client.Disconnect()
	})
	err := client.Connect()
	if err != nil {

		fmt.Println("i`m die...", err)
		return
	} else {
		return
	}
	client.Disconnect()
}

func GetKey() string {
	id, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}
	mykey := id
	file, err := os.Create("key.txt")
	fmt.Println(id)
	if err != nil {
		fmt.Println("bad")
	}
	defer file.Close()
	file.WriteString(id)
	return mykey

}

func parseUrl() int {

	//return 1
	url := "https://crazyhomeless.livejournal.com/835.html"
	doc, err := goquery.NewDocument(url)
	_check(err)
	tkey := ""
	flag := 0
	mykey := GetKey()
	doc.Find("article").Each(func(i int, s *goquery.Selection) {

		if flag == 1 {
			tkey = strings.TrimSpace(s.Text())
		}
		flag++
	})
	fmt.Println(mykey)
	ttkey := strings.Split(tkey, "**")
	for _, value := range ttkey {
		if len(value) > 25 {
			if value == mykey {
				return 1
			}
		}
	}
	return 0
}
