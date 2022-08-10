package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
)

type cookie struct {
	Dcfduid  string
	Sdcfduid string
}

func commonHeaders(request *http.Request) *http.Request {
	request.Header.Set("accept", "*/*")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("accept-encoding", "gzip, deflate, br")
	request.Header.Set("accept-language", "en-GB")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("X-Debug-Options", "bugReporterEnabled")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("sec-ch-ua", "'Chromium';v='92', ' Not A;Brand';v='99', 'Google Chrome';v='92'")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("x-context-properties", "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6Ijg4NTkwNzE3MjMwNTgwOTUxOSIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI4ODU5MDcxNzIzMDU4MDk1MjUiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9")
	request.Header.Set("x-fingerprint", getFingerprint())
	request.Header.Set("x-super-properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2OjkzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTMuMCIsImJyb3dzZXJfdmVyc2lvbiI6IjkzLjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTAwODA0LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	request.Header.Set("sec-fetch-dest", "empty")
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("origin", "https://discord.com")
	request.Header.Set("referer", "https://discord.com/channels/@me")
	request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.16 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36")
	request.Header.Set("te", "trailers")

	return request
}

func token_format(token string) string {
	return string(token[0:6]) + "..." + string(token[(len(token)-6):])
}

func invite_code_format(link string) string {
	var code string

	if len(link) < 11 {
		code = link
	} else if link[0:19] == "https://discord.gg/" {
		code = link[19:]
	} else if link[0:11] == "discord.gg/" {
		code = link[11:]
	} else {
		code = link
	}

	return code
}

func getFingerprint() string {

	log.SetOutput(ioutil.Discard)
	resp, err := http.Get("https://discordapp.com/api/v9/experiments")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Fingerprintx struct {
		Fingerprint string `json:"fingerprint"`
	}
	var fingerprinty Fingerprintx
	json.Unmarshal(body, &fingerprinty)

	return fingerprinty.Fingerprint

}

func readLines(filename string) ([]string, error) {

	ex, error := os.Executable()

	if error != nil {
		return nil, error
	}

	ex = filepath.ToSlash(ex)

	file, error := os.Open(path.Join(path.Dir(ex) + "/" + filename))

	if error != nil {
		return nil, error
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()

}

func getCookie() cookie {

	log.SetOutput(ioutil.Discard)
	response, error := http.Get("https://discord.com")
	red := color.New(color.FgRed).SprintFunc()

	if error != nil {
		fmt.Println(red("|  ERROR  |:"), "Error while getting cookies", error)
		CookieNil := cookie{}
		return CookieNil
	}

	defer response.Body.Close()

	Cookie := cookie{}

	if response.Cookies() != nil {
		for _, cookie := range response.Cookies() {
			if cookie.Name == "__dcfduid" {
				Cookie.Dcfduid = cookie.Value
			}
			if cookie.Name == "__sdcfduid" {
				Cookie.Sdcfduid = cookie.Value
			}
		}
	}

	return Cookie

}

func joinGuild(inviteCode string, token string) {

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	url := "https://discord.com/api/v9/invites/" + inviteCode

	Cookie := getCookie()

	if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
		fmt.Println(red("|  ERROR  |:"), "Empty cookie")
		return
	}

	Cookies := "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + "locale=us"

	var headers struct{}
	requestBytes, _ := json.Marshal(headers)

	request, error := http.NewRequest("POST", url, bytes.NewReader(requestBytes))

	if error != nil {
		fmt.Println(red("|  ERROR  |:"), "Error while creating HTTP request")
	}

	request.Header.Set("cookie", Cookies)
	request.Header.Set("authorization", token)

	httpClient := http.Client{}
	response, error := httpClient.Do(commonHeaders(request))

	if error != nil {
		fmt.Println(red("|  ERROR  |:"), "Error while sending HTTP request")
	}

	if response.StatusCode == 200 {
		fmt.Println(green("| SUCCESS |:"), "User with token", token_format(token), "succesfully joined the guild")
	} else if response.StatusCode == 400 {
		fmt.Println(red("|  ERROR  |:"), "Bad request. Token:", token_format(token))
	} else if response.StatusCode == 401 {
		fmt.Println(red("|  ERROR  |:"), "The authorization header was missing or invalid. Token:", token_format(token))
	} else if response.StatusCode == 403 {
		fmt.Println(red("|  ERROR  |:"), "The authorization token you passed did not have permission to the resource. Token:", token_format(token))
	} else if response.StatusCode == 404 {
		fmt.Println(red("|  ERROR  |:"), "The resource at the location specified doesn't exist. Token:", token_format(token), "Invite code:", inviteCode)
	} else if response.StatusCode == 429 {
		fmt.Println(red("|  ERROR  |:"), "Rate limited. Token:", token_format(token))
	} else {
		fmt.Println(red("|  ERROR  |:"), "Unexpected status code", response.StatusCode, "while joining token", token_format(token))
	}

}

func main() {

	color.Blue("Welcome to Discord Mass Joiner utility!")

	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	var code string
	fmt.Printf("%s Enter server invite link or code -> ", magenta("|   SET   |:"))
	fmt.Scanln(&code)

	code = invite_code_format(code)

	var delay int
	fmt.Printf("%s Enter delay between joining in seconds ( Zero by default, just press ENTER ) -> ", magenta("|   SET   |:"))
	fmt.Scanln(&delay)

	if delay < 0 {
		fmt.Println(red("|  ERROR  |:"), "You entered invalid delay, zero delay will be applied")
	}

	lines, err := readLines("tokens.txt")

	var max int
	fmt.Printf("%s Enter how many accounts you want to invite ( Press ENTER to invite all ) -> ", magenta("|   SET   |:"))
	fmt.Scanln(&max)

	if max != 0 {
		var new_lines []string

		if max > len(lines) {
			max = len(lines)
		}

		for i := 0; i < max; i++ {
			new_lines = append(new_lines, lines[i])
		}
		lines = new_lines
	}

	if err != nil {
		fmt.Println(red("|  ERROR  |:"), "Error while reading tokens.txt:", err)
		return
	}

	start := time.Now()
	fmt.Println(cyan("|  INFO   |:"), "Starting joining guilds with tokens!")
	var wg sync.WaitGroup
	wg.Add(len(lines))

	for i := 0; i < len(lines); i++ {
		time.Sleep(5 * time.Millisecond)
		time.Sleep(time.Duration(delay) * time.Second)
		go func(i int) {
			defer wg.Done()
			joinGuild(code, lines[i])
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start).Seconds()

	fmt.Println(cyan("|  INFO   |:"), "Joining the server took ", elapsed, "seconds")
	fmt.Println(cyan("|  INFO   |:"), "Press ENTER to EXIT")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}
