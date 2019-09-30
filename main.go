package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/denisbrodbeck/machineid"
	"github.com/gempir/go-twitch-irc"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows/registry"
)

var words []string
var patchf = ""
var amass []string
var mykey = ""
var msgcount = 0

type MyMainWindow struct {
	*walk.MainWindow
	edit    *walk.TextEdit
	account *walk.TextEdit
	count   *walk.TextEdit
	sherlok *walk.TextEdit
	path    string
}
type joinU struct {
	*walk.MainWindow
	edit       *walk.TextEdit
	path       string
	delay      *walk.TextEdit
	onemessage *walk.TextEdit
}

type masschat struct {
	*walk.MainWindow
	edit *walk.TextEdit
	chat *walk.TextEdit
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("405")
	}
	var wgg sync.WaitGroup
	mw := &MyMainWindow{}
	MW := MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "BeilusChat",
		MinSize:  Size{300, 20},
		Size:     Size{300, 50},
		MaxSize:  Size{400, 100},
		Layout:   VBox{},
		Children: []Widget{
			TextEdit{
				Text:     "Информационная консоль",
				MinSize:  Size{200, 40},
				AssignTo: &mw.edit, ReadOnly: true,
			},
			HSplitter{
				Children: []Widget{
					TextEdit{
						Text:    "Имя канала",
						MinSize: Size{100, 20},
						MaxSize: Size{100, 20},
						//				Size:     Size{200, 50},
						AssignTo: &mw.account, ReadOnly: false,
					},
					TextEdit{
						Text:    "кол-во ботов",
						MinSize: Size{100, 20},
						MaxSize: Size{100, 20},
						//				Size:     Size{200, 50},
						AssignTo: &mw.count, ReadOnly: false,
					},
					TextEdit{
						Text:    "Имя нужного бота",
						MinSize: Size{100, 20},
						MaxSize: Size{100, 20},
						//				Size:     Size{200, 50},
						AssignTo: &mw.sherlok, ReadOnly: false,
					},
				},
			},
			HSplitter{
				Children: []Widget{
					PushButton{
						MinSize:   Size{100, 30},
						MaxSize:   Size{100, 30},
						Text:      "Аккаунты",
						OnClicked: mw.pbClicked,
					},
					PushButton{
						MinSize: Size{100, 30},
						MaxSize: Size{100, 30},
						Text:    "Открыть папку",
						OnClicked: func() {
							err = exec.Command("explorer", dir).Start()
						},
					},
					PushButton{
						MinSize: Size{100, 30},
						MaxSize: Size{100, 30},
						Text:    "авточаты",
						OnClicked: func() {
							GetKey()
							if parseUrl() == 1 {
								cc, err := strconv.Atoi(mw.count.Text())
								if err != nil {
									mw.edit.SetText("Введи верное количество ботов")
								} else {
									if len(patchf) < 1 {
										mw.edit.SetText("Выбери путь до аккаунтов")
									} else {
										wgg.Add(1)
										go func(room string, cc int, patchz string, vacc []string) {
											defer wgg.Done()
											juser(room, cc, patchz, vacc)
										}(mw.account.Text(), cc, patchf, amass)
									}
								}
							} else {
								mw.edit.SetText("For activate, send you key to starchenkoleo@gmail.com " + mykey)
							}
						},
					},
					PushButton{
						Text: "Чат(до 20 чел)",
						OnClicked: func() {
							GetKey()
							if parseUrl() == 1 {
								cc, err := strconv.Atoi(mw.count.Text())
								if err != nil {
									mw.edit.SetText("Введи верное количество ботов")
								} else {
									if cc > 20 {
										cc = 20
									}
									if len(patchf) < 1 {
										mw.edit.SetText("Выбери путь до аккаунтов")
									} else {
										wgg.Add(1)
										go func(room string, cc int, patchz string, vacc []string) {
											defer wgg.Done()
											massbot(room, cc, patchz, vacc)
										}(mw.account.Text(), cc, patchf, amass)
									}
								}
							} else {
								mw.edit.SetText("Для активации, отправь ключ starchenkoleo@gmail.com " + mykey)
							}
						},
					},
					PushButton{
						MinSize: Size{100, 30},
						MaxSize: Size{100, 30},
						Text:    "Найти бота",
						OnClicked: func() {
							if len(patchf) < 1 {
								mw.edit.SetText("Выбери путь до аккаунтов")
							} else {
								wgg.Add(1)
								go func(room string, name string, patchf string) {
									defer wgg.Done()
									var lostmans []string
									zz, err := ioutil.ReadFile(patchf) // just pass the file name
									if err != nil {
										fmt.Print(err)
									}
									for _, mmsg := range strings.Split(string(zz), "\n") {
										if strings.Contains(strings.ToLower(mmsg), strings.ToLower(name)) {
											lostmans = append(lostmans, mmsg)
										}
									}
									if len(lostmans) > 0 {
										if len(lostmans) > 36 {
											mw.edit.SetText("Нашлось слишком много ботов")
										} else {
											fmt.Println(room, name, patchf)
											tryFound(room, name, lostmans)
										}
									} else {
										mw.edit.SetText("Не удалось найти бота")
									}
								}(mw.account.Text(), mw.sherlok.Text(), patchf)
							}
						},
					},
				},
			},
		},
	}
	if _, err := MW.Run(); err != nil {
		os.Exit(1)
	}
}

func tryFound(room string, name string, vmass []string) {
	myfont := new(Font)
	myfont.Bold = true
	myfont.PointSize = 10
	var wg sync.WaitGroup
	mw := &masschat{}
	zzz := []Widget{}
	zz1 := []Widget{}
	zz2 := []Widget{}
	for i, x := range vmass {
		sender := x
		zzz = append(zzz, PushButton{
			Text: strings.Split(vmass[i], ":")[0] + strconv.Itoa(i),
			OnClicked: func() {
				wg.Add(1)
				go func(value string, room string, rmsg string) {
					defer wg.Done()
					fmt.Println(value)
					chat(room, value, rmsg)
				}(sender, room, mw.edit.Text())
			},
		})
	}
	if len(zzz) > 4 {
		zz1 = zzz[:len(zzz)/2]
		zz2 = zzz[len(zzz)/2:]
	} else {
		zz1 = zzz
	}
	MW := MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    room + " MassChat",
		MinSize:  Size{500, 500},
		Size:     Size{1000, 1000},

		Layout: VBox{},

		Children: []Widget{
			TextEdit{
				MinSize: Size{300, 300},

				Font:     *myfont,
				AssignTo: &mw.chat, ReadOnly: true,
			},
			TextEdit{
				AssignTo: &mw.edit, ReadOnly: false,
			},
			PushButton{
				Text: "Очистить окно сообщений",
				OnClicked: func() {
					mw.edit.SetText("")
				},
			},
			PushButton{
				Text: "Отслеживать чат",
				OnClicked: func() {
					wg.Add(1)
					go func(room string, editt *walk.TextEdit) {
						defer wg.Done()
						spyRoom(room, editt)
					}(room, mw.chat)
				},
			},
			HSplitter{
				MinSize:  Size{80, 20},
				MaxSize:  Size{80, 20},
				Children: zz1,
			},
			HSplitter{
				MinSize:  Size{80, 20},
				MaxSize:  Size{80, 20},
				Children: zz2,
			},
		},
	}
	if _, err := MW.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func takeSender(sender string, vmass []string) string {
	for _, value := range vmass {
		fmt.Println(sender, value)
		if sender == (strings.Split(value, ":"))[0] {
			fmt.Println("RETURN", value)
			return value
		}
	}
	return ""
}

func massbot(room string, count int, patchz string, vmass []string) {
	myfont := new(Font)
	myfont.Bold = true
	myfont.PointSize = 10
	var wg sync.WaitGroup
	mw := &masschat{}
	zzz := []Widget{}
	zz1 := []Widget{}
	zz2 := []Widget{}
	for i, x := range vmass {
		if i >= count {
			break
		}
		sender := x
		zzz = append(zzz, PushButton{

			Text: strings.Split(vmass[i], ":")[0] + strconv.Itoa(i),
			OnClicked: func() {
				wg.Add(1)
				go func(value string, room string, rmsg string) {
					defer wg.Done()
					fmt.Println(value)
					chat(room, value, rmsg)
				}(sender, room, mw.edit.Text())
			},
		})

	}
	zz1 = zzz[:len(zzz)/2]
	zz2 = zzz[len(zzz)/2:]

	MW := MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    room + " MassChat",
		MinSize:  Size{500, 50},
		Size:     Size{500, 50},
		Layout:   VBox{},

		Children: []Widget{
			TextEdit{
				MinSize:  Size{700, 500},
				Font:     *myfont,
				AssignTo: &mw.chat, ReadOnly: true,
			},
			TextEdit{
				AssignTo: &mw.edit, ReadOnly: false,
			},
			PushButton{
				Text: "Очистить чат",
				OnClicked: func() {
					mw.edit.SetText("")
				},
			},
			PushButton{
				Text: "Обновлять чат",
				OnClicked: func() {
					wg.Add(1)
					go func(room string, editt *walk.TextEdit) {
						defer wg.Done()
						spyRoom(room, editt)
					}(room, mw.chat)
				},
			},
			HSplitter{
				MinSize:  Size{50, 20},
				MaxSize:  Size{50, 20},
				Children: zz1,
			},
			HSplitter{
				MinSize:  Size{50, 20},
				MaxSize:  Size{50, 20},
				Children: zz2,
			},
		},
	}

	if _, err := MW.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		//os.Exit(1)
	}

}

func juser(room string, count int, patchz string, vmass []string) {
	var flag = 1
	var wg sync.WaitGroup
	mww := &joinU{}
	MW := MainWindow{
		AssignTo: &mww.MainWindow,
		Title:    room,
		MinSize:  Size{300, 20},
		Size:     Size{300, 50},
		Layout:   VBox{},
		Children: []Widget{
			TextEdit{
				Text:    "Информационное окно",
				MinSize: Size{200, 40},

				AssignTo: &mww.edit, ReadOnly: true,
			},
			PushButton{
				Text: "Остановить ботов",
				OnClicked: func() {
					flag = 0
					mww.edit.SetText("Wait 1 min for full stop please.")
				},
			},
			PushButton{
				Text: "Вкл список чата",
				OnClicked: func() {
					if count > len(vmass)-1 {
						count = len(vmass) - 1
					}
					flag = 1
					wg.Add(1)
					go func(room string, flag *int) {

						defer wg.Done()
						for _, value := range vmass[:count] {
							time.Sleep(100 * time.Millisecond)
							wg.Add(1)
							go func(value string, room string, flag *int) {
								defer wg.Done()
								conchat(value, room, flag)
							}(value, room, flag)
						}
					}(room, &flag)
					//

					cc := strconv.Itoa(count)

					mww.edit.SetText(cc + " Боты работают")
				},
			},
			TextEdit{
				Text:    "Задержка. 1=секунда",
				MinSize: Size{20, 20},
				MaxSize: Size{20, 20},
				//				Size:     Size{200, 50},
				AssignTo: &mww.delay, ReadOnly: false,
			},
			TextEdit{
				Text:    "Одно сообщение",
				MinSize: Size{20, 20},
				MaxSize: Size{20, 20},
				//				Size:     Size{200, 50},
				AssignTo: &mww.onemessage, ReadOnly: false,
			},
			PushButton{
				Text: "Рандомные сообщения",
				OnClicked: func() {
					cc, err := strconv.Atoi(mww.delay.Text())
					if err != nil {
						mww.edit.SetText("Bad count delay!!!")
					} else {
						if cc == 0 {
							cc = 1
						}
						rand.Seed(time.Now().UTC().UnixNano())

						//if count > len(vmass)-1 {
						count = len(vmass) - 1
						//	}

						dlg := new(walk.FileDialog)

						dlg.FilePath = mww.path
						dlg.Title = "Select File"
						dlg.Filter = "Exe files (*.txt)|*.txt|All files (*.*)|*.*"

						if ok, err := dlg.ShowOpen(mww); err != nil {
							mww.edit.AppendText("Error : File Open\r\n")
							return
						} else if !ok {
							mww.edit.AppendText("Cancel\r\n")
							return
						}
						mww.path = dlg.FilePath
						patchz := dlg.FilePath
						bots := strconv.Itoa(cc)
						mww.edit.SetText(bots + " bots work now")
						wg.Add(1)
						flag = 1
						go func(cc int, patchz string, vmass []string, room string, count int, flag *int) {
							defer wg.Done()
							flagchat(cc, patchz, vmass, room, count, flag)
						}(cc, patchz, vmass, room, count, &flag)
					}
				},
			},
			PushButton{
				Text: "Сообщения по порядку",
				OnClicked: func() {
					cc, err := strconv.Atoi(mww.delay.Text())
					if err != nil {
						mww.edit.SetText("Bad count delay!!!")
					} else {
						if cc == 0 {
							cc = 1
						}
						rand.Seed(time.Now().UTC().UnixNano())

						if count > len(vmass)-1 {
							count = len(vmass) - 1
						}

						dlg := new(walk.FileDialog)

						dlg.FilePath = mww.path
						dlg.Title = "Select File"
						dlg.Filter = "Exe files (*.txt)|*.txt|All files (*.*)|*.*"

						if ok, err := dlg.ShowOpen(mww); err != nil {
							mww.edit.AppendText("Error : File Open\r\n")
							return
						} else if !ok {
							mww.edit.AppendText("Cancel\r\n")
							return
						}
						mww.path = dlg.FilePath
						patchz := dlg.FilePath
						bots := strconv.Itoa(cc)
						mww.edit.SetText(bots + " Боты работают")
						wg.Add(1)
						flag = 1
						go func(cc int, patchz string, vmass []string, room string, count int, flag *int) {
							defer wg.Done()
							flagchat2(cc, patchz, vmass, room, count, flag)
						}(cc, patchz, vmass, room, count, &flag)
					}
				},
			},
			PushButton{
				Text: "Спам одного сообщения",
				OnClicked: func() {
					cc, err := strconv.Atoi(mww.delay.Text())
					if err != nil {
						mww.edit.SetText("Неправильное значение задержки")
					} else {
						if cc == 0 {
							cc = 1
						}

						mww.edit.SetText("Бот работает")
						wg.Add(1)
						flag = 1

						go func(cc int, msgg string, vmass []string, room string, count int, flag *int) {
							defer wg.Done()
							onemsg(cc, msgg, vmass, room, count, flag)
						}(cc, mww.onemessage.Text(), vmass, room, count, &flag)
					}
				},
			},
		},
	}

	if _, err := MW.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag = 0
		//os.Exit(1)
	}

	wg.Wait()
}

//func msgcountFunc()

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

func onemsg(cc int, msg string, vmass []string, room string, count int, flag *int) {

	var wg sync.WaitGroup
	ccc := 0

	if len(msg) > 0 {

		for _, acc := range vmass {
			if ccc == count {
				ccc = 0
			}
			if *flag == 0 {
				return
			}
			if *flag == 0 {
				return
			}
			wg.Add(1)
			time.Sleep(time.Duration(cc*1000) * time.Millisecond)
			go func(value string, room string, rmsg string) {
				defer wg.Done()
				chat(room, value, rmsg)
			}(acc, room, msg)
			ccc++
		}
	} else {
		return
	}
}

func flagchat2(cc int, patchz string, vmass []string, room string, count int, flag *int) {
	var msg []string
	var wg sync.WaitGroup
	zz, err := ioutil.ReadFile(patchz) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	msg = msg[:0]
	ccc := 0
	for _, mmsg := range strings.Split(string(zz), "\n") {
		msg = append(msg, mmsg)
	}
	if len(msg) > 0 {
		for {
			for _, rmsg := range msg {
				if ccc == count {
					ccc = 0
				}
				if *flag == 0 {
					return
				}
				if *flag == 0 {
					return
				}
				wg.Add(1)
				time.Sleep(time.Duration(cc*1000) * time.Millisecond)
				go func(value string, room string, rmsg string) {
					defer wg.Done()
					chat(room, value, rmsg)
				}(vmass[ccc], room, rmsg)
				ccc++
			}
		}
	} else {
		return
	}

}
func flagchat(cc int, patchz string, vmass []string, room string, count int, flag *int) {
	var msg []string
	var wg sync.WaitGroup
	zz, err := ioutil.ReadFile(patchz) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	msg = msg[:0]
	for _, mmsg := range strings.Split(string(zz), "\n") {
		msg = append(msg, mmsg)
	}
	r := 0
	r2 := 0
	if count == 0 {
		count = 1
	}
	cc = cc + randInt(0, 6)
	if len(msg) > 0 {
		for {
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
		time.Sleep(time.Duration(900) * time.Millisecond)
		msgcount++
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

func rrchat() {
	client := twitch.NewClient("beilus", "oauth:duz4vc6rztpiwf2aa9j51nlqaixaus")
	client.Join("beilus")
	client.OnConnect(func() {
		tsong := ""
		for {
			b, err := ioutil.ReadFile("C:/Users/Beilus/Desktop/stream/SongName.txt") // just pass the file name
			if err != nil {
				fmt.Print(err)
			}
			if tsong != string(b) {
				tsong = string(b)
				client.Say("beilus", "Now "+tsong)
				time.Sleep(2000 * time.Millisecond)
				client.Say("beilus", "Best casino. 100 free spins https://is.gd/AuZvtC")
			}
			time.Sleep(5000 * time.Millisecond)
		}

	})
	err := client.Connect()
	if err != nil {
		fmt.Println("i`m die...", err)
		return
	}
}

func (mw *MyMainWindow) pbClicked() {
	GetKey()
	dlg := new(walk.FileDialog)

	dlg.FilePath = mw.path
	dlg.Title = "Select File"
	dlg.Filter = "Exe files (*.txt)|*.txt|All files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.edit.AppendText("Error : File Open\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("Cancel\r\n")
		return
	}
	mw.path = dlg.FilePath
	patchf = dlg.FilePath

	zz, err := ioutil.ReadFile(patchf) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	amass = amass[:0]
	for _, value := range strings.Split(string(zz), "\n") {
		if len(strings.Split(string(value), ":")) == 3 {
			amass = append(amass, value)
		}
	}
	mw.edit.SetText(strconv.Itoa(len(amass)) + " Аккаунтов в наличии")
}

func SetKey() string {
	file, err := os.Create("key.txt")
	if err != nil {
		return "bad"
	}
	defer file.Close()
	rand.Seed(time.Now().UTC().UnixNano())
	r := randomString(30)
	file.WriteString(r)
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err := k.SetStringValue("Chatbot", r); err != nil {
		fmt.Println(err)
	}
	z, _, err := k.GetStringValue("Chatbot")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(z)
	if err := k.Close(); err != nil {
		fmt.Println(err)
	}
	return (z)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func _check(err error) {
	if err != nil {
		panic(err)
	}
}

func spyRoom(room string, editt *walk.TextEdit) {
	/*	var msgmass []string
		client := twitch.NewClient("Meemal6412", "oauth:3e139brdfy7jxpki47ptdpcguzpn5o")

		client.Join(room)
		err := client.Connect()
		if err != nil {
			panic(err)
		}
	*/
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
