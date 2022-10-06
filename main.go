package main

//imports
import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
)

// structs
type structs struct {
	Dcfd  string
	Sdcfd string
	ID    string  `json:"id"`
	Time  float64 `json:"retry_after"`
}

// Functions
func Create_gc(token string) {
	payload := map[string]string{
		"recipients": "",
	}
	Client := http.Client{
		//Transport: &http.Transport{
			//Proxy: http.ProxyURL(p),
		//},
	}
	xp, _ := json.Marshal(payload)
	Cookie := Build_cookie()
	Cookies := "__dcfduid=" + Cookie.Dcfd + "; " + "__sdcfduid=" + Cookie.Sdcfd + "; "
	req, _ := http.NewRequest("POST", "https://discord.com/api/v9/users/@me/channels", bytes.NewBuffer(xp))
	for x, o := range map[string]string{
		"accept":               "*/*",
		"accept-encoding":      "gzip, deflate, br",
		"accept-language":      "en-US,en-NL;q=0.9,en-GB;q=0.8",
		"authorization":        token,
		"content-type":         "application/json",
		"cookie":               Cookies,
		"origin":               "https://discord.com",
		"referer":              "https://discord.com/channels/@me/",
		"sec-fetch-dest":       "empty",
		"sec-fetch-mode":       "cors",
		"sec-fetch-site":       "same-origin",
		"user-agent":           "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
		"x-context-properties": "eyJsb2NhdGlvbiI6Ik5ldyBHcm91cCBETSJ9",
		"x-debug-options":      "bugReporterEnabled",
		"x-discord-locale":     "en-US",
		"x-super-properties":   "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDA2Iiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTUwNzQ4LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==",
	} {
		req.Header.Set(x, o)
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
		gc_id := ResponseBody.ID
		fmt.Println("(\033[32m+\033[39m) Made GC | [ID]: ", gc_id)
		f, err := os.OpenFile("gcs.txt", os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, ers := f.WriteString(gc_id + "\n")
		if ers != nil {
			log.Fatal(ers)
		}
		
		
	} else if resp.StatusCode == 429 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var ResponseBody structs
		json.Unmarshal(body, &ResponseBody)
		timeout := ResponseBody.Time
		fmt.Println("(\033[33m/\033[39m) Rate Limmit | [TIME]: ", timeout, "\033[u")
		time.Sleep(100 *time.Second)


	} else {
		fmt.Print("(\033[31mx\033[39m)[ERROR] : ", resp)
	}

}


func member(token string, user_id string, gc_id string, method string) {
	Client := &http.Client{}
	Cookie := Build_cookie()
	Cookies := "__dcfduid=" + Cookie.Dcfd + "; " + "__sdcfduid=" + Cookie.Sdcfd + "; "
	req, err := http.NewRequest(""+method+"", "https://discord.com/api/v9/channels/"+gc_id+"/recipients/"+user_id+"", nil)
	if err != nil {
		log.Fatal(err)
	}
	for x,o := range map[string]string{
		"accept":               "*/*",
		"accept-encoding":      "gzip, deflate, br",
		"accept-language":      "en-US,en-NL;q=0.9,en-GB;q=0.8",
		"authorization":        token,
		"content-length": 		"0",
		"cookie":               Cookies,
		"origin":               "https://discord.com",
		"referer":              "https://discord.com/channels/@me/",
		"sec-fetch-dest":       "empty",
		"sec-fetch-mode":       "cors",
		"sec-fetch-site":       "same-origin",
		"user-agent":           "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
		"x-context-properties": "eyJsb2NhdGlvbiI6IkFkZCBGcmllbmRzIHRvIERNIn0=",
		"x-debug-options":      "bugReporterEnabled",
		"x-discord-locale":     "en-US",
		"x-super-properties":   "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDA2Iiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTUwNzQ4LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==",
	} {
		req.Header.Set(x,o)
	}
	resp, err := Client.Do(req)
	if resp.StatusCode == 204 {
		if method == "PUT" {
			x :="Added  "
			fmt.Println("(\033[32m+\033[39m) "+x+" [USER]: "+user_id+" To [GC]: ", gc_id)
		} else if method == "DELETE" {
			x := "Removed"
			fmt.Println("(\033[32m-\033[39m) "+x+" [USER]: "+user_id+" To [GC]: ", gc_id)
		}
		
	} else if resp.StatusCode == 429 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var ResponseBody structs
		json.Unmarshal(body, &ResponseBody)
		timeout := ResponseBody.Time
		fmt.Println("(\033[33m/\033[39m) Rate Limmit | [TIME]: ", timeout)

	} else {
		fmt.Print("(\033[31mx\033[39m) [ERROR] : ", resp)
	}
}



func change_name() {

}


func spam_gc() {

}




// Build Header Data
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

func build_xconst() {
	//str := `s`
}


//modules
func DecodeBr(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	br := brotli.NewReader(r)
	return ioutil.ReadAll(br)
}

var proxy = "us.proxiware.com:2000"
var p,_ = url.Parse(func() string {
	if !strings.Contains(proxy, "http://") {
		return "http://" + proxy
	}
	return proxy
}())

func read_tokens() ([]string, error) {
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

func read_ids() ([]string, error) {
	file, err := os.Open("gcs.txt")
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

func colors() []string {
	//blue := "\033[34m"
	//cyan := "\033[36m"
	var clr []string
	return clr
}







//main
func main() {
	enclogo := ""
	logo,_ := base64.StdEncoding.DecodeString(enclogo)
	var wg sync.WaitGroup
	fmt.Printf("%q\n", logo)
	scn := bufio.NewScanner(os.Stdin)
	

	//Data
	lines, err := read_tokens()
	if err != nil {
		log.Fatal(err)
	}	
	ids, err := read_ids()
	if err != nil {
		log.Fatal(err)
	}


	//inputs
	fmt.Println("	[1] Create GC\n	[2] Add Member\n	[3] Remove Member\n	[4] Remove/Add Member")
	fmt.Print("\n\n	[]> ")
	scn.Scan()
	choice := scn.Text()

	
	//create Groupchat call
	if choice == "1" {
		fmt.Println("\n	[1] Single Token\n	[2] Multi Token\n")	
		fmt.Print("	[]> ")
		scn.Scan()
		choice := scn.Text()
		if choice == "1" {
		    fmt.Print("[TOKEN]> ")
			scn.Scan()
			token := scn.Text()
			
			for true {
				Create_gc(token)
			}
		}

		if choice == "2" {
			wg.Add(len(lines))
			for i := 0; i < len(lines); i++ {
				for true {
					Create_gc(lines[i])
				}
			}
		}
		

		//add member call
	} else if choice == "2" {
		fmt.Print("	[USER ID]> ")
		scn.Scan()
		user := scn.Text()
		wg.Add(len(lines))
		mth := "PUT"
		for i := 0; i < len(ids); i++ {
			member(lines[i], user, ids[i], mth)
		}

	} else if choice == "3" {
		fmt.Print("	[USER ID]> ")
		scn.Scan()
		user := scn.Text()
		wg.Add(len(lines))
		mth := "DELETE"
		for i := 0; i < len(ids); i++ {
			member(lines[i], user, ids[i], mth)
		}

	} else if choice == "4" {
		fmt.Print("	[USER ID]> ")
		scn.Scan()
		user := scn.Text()
		wg.Add(len(lines))
		for i := 0; i < len(ids); i++ {
			mth := "PUT"
			member(lines[i], user, ids[i], mth)
			mths := "DELETE"
			member(lines[i], user, ids[i], mths)
		}
	} else {
		fmt.Println("\n[\033[31mERROR\033[39m] Wrong Input")
		time.Sleep(1 *time.Second)
        cmd := exec.Command("cmd", "/c", "cls") 
        cmd.Stdout = os.Stdout
        cmd.Run()
		main()
	} 

	
	
	
}
