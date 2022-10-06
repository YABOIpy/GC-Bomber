package main

import (
	"bytes"
	"encoding/json"	
	"github.com/andybalholm/brotli"
	"io/ioutil"
	"bufio"
	"os"
	"sync"
	"fmt"
	"log"
	"net/http"
)


type structs struct {
	Dcfd  string
	Sdcfd string
	ID string `json:"id"`
}




func Create_gc(token string) {
	payload := map[string]string{
		"recipients": "",
	}
	Client := http.Client{}
	xp,_ := json.Marshal(payload)
	Cookie := Build_cookie()
	Cookies := "__dcfduid=" + Cookie.Dcfd + "; " + "__sdcfduid=" + Cookie.Sdcfd + "; "
	req,_ := http.NewRequest("POST", "https://discord.com/api/v9/users/@me/channels", bytes.NewBuffer(xp))
	for x,o := range map[string]string{
		"accept": "*/*",
		"accept-encoding": "gzip, deflate, br",
		"accept-language": "en-US,en-NL;q=0.9,en-GB;q=0.8",
		"authorization": token,
		"content-type": "application/json",
		"cookie":Cookies,
		"origin": "https://discord.com",
		"referer": "https://discord.com/channels/@me/",
		"sec-fetch-dest": "empty",
		"sec-fetch-mode": "cors",
		"sec-fetch-site": "same-origin",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
		"x-context-properties": "eyJsb2NhdGlvbiI6Ik5ldyBHcm91cCBETSJ9",
		"x-debug-options": "bugReporterEnabled",
		"x-discord-locale": "en-US",
		"x-super-properties": "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDA2Iiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTUwNzQ4LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==",
	} {
		req.Header.Set(x,o)
	}
	resp, err := Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		p, m := DecodeBr(body)
		if m != nil {
			log.Fatal(m)
		}
		var ResponseBody structs
		json.Unmarshal(p, &ResponseBody)
		fmt.Print("[>] Made GC | [ID]: ", ResponseBody.ID)
	} else if resp.StatusCode == 429{
		fmt.Println("[/] Rate Limmit | [TIME]: ", resp)//ratetime
	} else {
		fmt.Print("[ERROR] : ", resp)
	}

}


func Build_cookie() structs {
	log.SetOutput(ioutil.Discard)
	req, err := http.Get("https://discord.com")
	if err != nil {
		log.Fatal(err)
		CookieNil := structs{}
		return CookieNil
	}
	defer req.Body.Close()

	Cookie := structs{}
	if req.Cookies() != nil {
		for _, cookie := range req.Cookies() {
			if cookie.Name == "__dcfduid" {
				Cookie.Dcfd = cookie.Value
			}
			if cookie.Name == "__sdcfduid" {
				Cookie.Sdcfd = cookie.Value
			}
		}
	}
	return Cookie
}
func spam_gc() {

}

func test(token string) {
	fmt.Println(token)
}




func DecodeBr(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	br := brotli.NewReader(r)
	return ioutil.ReadAll(br)
}

func readLines() ([]string, error) {
	file, err := os.Open("tokens.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}





func main() {
	logo := "nigger \n nugger"
	fmt.Println(logo)
	lines, err := readLines()
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(lines))
	for i := 0; i < len(lines); i++ {
		Create_gc(lines[i])
	}
}
