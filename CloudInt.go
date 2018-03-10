
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"crypto/tls"
	"github.com/fatih/color"
)

var faliure = color.New(color.FgRed).PrintlnFunc()
var success = color.New(color.Bold, color.FgGreen).PrintlnFunc()
var info = color.New(color.FgGreen).PrintlnFunc()
//var debug = color.New(color.FgCyan).PrintlnFunc()

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
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
func checkAlias(alias string) {
	// IAM
	// Checks if an alias / Account ID exists
	url := ("https://" + alias + ".signin.aws.amazon.com")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		success("[+] Exists: " + alias)
	} else if resp.StatusCode == 404 {
		//info(alias + " does not exist")
	} else {
		faliure("error")
	}

}
func checkSpace(alias string) {
	// Checks if a Digital Ocean space exists
	regions := []string{"nyc3", "ams3", "sgp1"}
	tr := &http.Transport{
	        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	    }
	client := &http.Client{Transport: tr}
	for _, region := range regions {
		url := ("https://" + alias + "." + region + ".digitaloceanspaces.com")
		resp, err := client.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode == 400 || resp.StatusCode == 403 || resp.StatusCode == 200 {
			info("[*] Exists [" + region + "]: " + alias)
			if resp.StatusCode == 200 {
				success("\t[+] " + alias + " anonymous access allowed !!!")
			}else{
				//info("\t[-] " + alias + " anonymous access forbidden")
			}
		} else if resp.StatusCode == 404 {
			//info(alias + " does not exist")
		} else {
			faliure("error")
		}
	}


}
func checkBucket(alias string) {
	// S3
	// Checks if a S3 Bucket exists
	url := ("http://" + alias + ".s3.amazonaws.com")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 400 || resp.StatusCode == 403 || resp.StatusCode == 200 {
		info("[*] Exists: " + alias)
		if resp.StatusCode == 200 {
			success("\t[+] " + alias + " anonymous access allowed !!!")
		}else{
			//info("\t[-] " + alias + " anonymous access forbidden")
		}
	} else if resp.StatusCode == 404 {
		//faliure(alias + " does not exist")
	} else {
		faliure("error")
	}
}
func generateWordlist(alias string) []string {
	wordlist := make([]string, 1)
	wordlist[0] = alias
	// adding years
	year, _, _ := time.Now().Date()
	wordlist = append(wordlist, alias+strconv.Itoa(year))
	for i := 0; i < 10; i++ {
		wordlist = append(wordlist, alias+strconv.Itoa(year+i))
		wordlist = append(wordlist, alias+"_"+strconv.Itoa(year+i))
		wordlist = append(wordlist, alias+"__"+strconv.Itoa(year+i))
		wordlist = append(wordlist, alias+"."+strconv.Itoa(year+i))
		wordlist = append(wordlist, alias+"-"+strconv.Itoa(year+i))
		wordlist = append(wordlist, alias+"--"+strconv.Itoa(year+i))
	}

	// adding numbers
	for i := 0; i < 10; i++ {
		wordlist = append(wordlist, alias+strconv.Itoa(i))
		wordlist = append(wordlist, alias+"_"+strconv.Itoa(i))
		wordlist = append(wordlist, alias+"__"+strconv.Itoa(i))
		wordlist = append(wordlist, alias+"."+strconv.Itoa(i))
		wordlist = append(wordlist, alias+"-"+strconv.Itoa(i))
		wordlist = append(wordlist, alias+"--"+strconv.Itoa(i))
	}

	// adding interesting word
	words := [...]string{"appliance", "test", "testing", "development", "acceptance", "production", "bucket", "buckets", "s3", "space", "spaces", "local", "trunk", "integration", "qa", "stage", "staging", "pre-prod", "pre-production", "live", "education", "backup", "deployment", "uat", "dev", "server", "client", "secret"}
	for i := range words {
		wordlist = append(wordlist, alias+"_"+words[i])
		wordlist = append(wordlist, alias+"__"+words[i])
		wordlist = append(wordlist, alias+"."+words[i])
		wordlist = append(wordlist, alias+"-"+words[i])
		wordlist = append(wordlist, alias+"--"+words[i])
	}
	return wordlist

}
func main() {
	var wg sync.WaitGroup
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nCloudInt v0.1 - Cloud Container Enumerator via HTTP response code")
	fmt.Println("-----------------------------------------------------------------")
	for {
		info("\nEnumeration available:")
		fmt.Println("[1] AWS Alias / Account ID enumeration")
		fmt.Println("[2] AWS S3 Bucket enumeration")
		fmt.Println("[3] DigitalOcean Spaces enumeration")
		fmt.Println("\n[99] Exit")
		fmt.Print("\n\nSelect an option ")
		module, _ := reader.ReadString('\n')
		module = strings.Replace(module, "\n", "", -1)
		switch module{
			case "1": // AWS Alias / Account ID enumeration
			info("\nAWS Alias / Account ID enumeration:")
			fmt.Println("[1] Single Alias / Account ID")
			fmt.Println("[2] Wordlist of Aliases / Account IDs")
			fmt.Println("[3] Generate a wordlist based on a word\n")
			fmt.Println("\n[99] Exit")
			fmt.Print("\n\nSelect an option ")
			iamOption, _ := reader.ReadString('\n')
			iamOption = strings.Replace(iamOption, "\n", "", -1)
			switch iamOption{
				case "1":
					fmt.Print("\nInsert a word: ")
					alias, _ := reader.ReadString('\n')
					alias = strings.Replace(alias, "\n", "", -1)
					checkAlias(alias)
				case "2":
					fmt.Print("\nInsert a wordlist (full path): ")
					wordlistFile, _ := reader.ReadString('\n')
					wordlistFile = strings.Replace(wordlistFile, "\n", "", -1)
					lines, err := readLines(wordlistFile)
					if err != nil {
						faliure("File error")
					}
					for line := range lines {
						wg.Add(1)
						go func(alias string) {
							checkAlias(alias)
							wg.Done()
						}(lines[line])
					}
					wg.Wait()
				case "3":
					fmt.Print("\nInsert a word to create the wordlist: ")
					alias, _ := reader.ReadString('\n')
					alias = strings.Replace(alias, "\n", "", -1)
					wordlist := generateWordlist(alias)
					for i := range wordlist {
						wg.Add(1)
						go func(alias string) {
							if !strings.Contains(alias, ".") {
								checkAlias(alias)
							}
							wg.Done()
						}(wordlist[i])
					}
					wg.Wait()
				case "99":
					info("\nExiting...")
					os.Exit(0)
				default:
					faliure("Wrong option")
				}
		case "2": // AWS S3 Bucket enumeration
			info("\nAWS S3 Bucket enumeration:")
			fmt.Println("[1] Single Bucket")
			fmt.Println("[2] Wordlist of Buckets")
			fmt.Println("[3] Generate a wordlist based on a word")
			fmt.Println("\n[99] Exit")
			fmt.Print("\n\nSelect an option ")
			s3Option, _ := reader.ReadString('\n')
			s3Option = strings.Replace(s3Option, "\n", "", -1)
			switch s3Option{
				case "1":
					fmt.Print("Insert a word: ")
					alias, _ := reader.ReadString('\n')
					alias = strings.Replace(alias, "\n", "", -1)
					checkBucket(alias)
				case "2":
					fmt.Print("Insert a wordlist (full path): ")
					wordlistFile, _ := reader.ReadString('\n')
					wordlistFile = strings.Replace(wordlistFile, "\n", "", -1)
					lines, err := readLines(wordlistFile)
					if err != nil {
						faliure("File error")
					}
					for line := range lines {
						wg.Add(1)
						go func(alias string) {
							checkBucket(alias)
							wg.Done()
							}(lines[line])
						}
						wg.Wait()
				case "3":
					fmt.Print("Insert a word to create the wordlist: ")
					alias, _ := reader.ReadString('\n')
					alias = strings.Replace(alias, "\n", "", -1)
					wordlist := generateWordlist(alias)
					for i := range wordlist {
						wg.Add(1)
						go func(alias string) {
							checkBucket(alias)
							wg.Done()
						}(wordlist[i])
					}
					wg.Wait()
				case "99":
					info("\nExiting...")
					os.Exit(0)
				default:
					faliure("Wrong option")
			}
		case "3": // Digital Ocean
			info("\nDigitalOcean Spaces enumeration:")
			fmt.Println("[1] Single Spaces")
			fmt.Println("[2] Wordlist of Spaces")
			fmt.Println("[3] Generate a wordlist based on a word\n")
			fmt.Println("\n[99] Exit")
			fmt.Print("\n\nSelect an option ")
			dOption, _ := reader.ReadString('\n')
			dOption = strings.Replace(dOption, "\n", "", -1)
			switch dOption{
				case "1":
					fmt.Print("Insert a word: ")
					alias, _ := reader.ReadString('\n')
					alias = strings.Replace(alias, "\n", "", -1)
					checkSpace(alias)
				case "2":
					fmt.Print("Insert a wordlist (full path): ")
					wordlistFile, _ := reader.ReadString('\n')
					wordlistFile = strings.Replace(wordlistFile, "\n", "", -1)
					lines, err := readLines(wordlistFile)
					if err != nil {
						faliure("File error")
					}
					for line := range lines {
						wg.Add(1)
						go func(alias string) {
							checkSpace(alias)
							wg.Done()
						}(lines[line])
					}
					wg.Wait()
				case "3":
					fmt.Print("Insert a word to create the wordlist: ")
					alias, _ := reader.ReadString('\n')
					alias = strings.Replace(alias, "\n", "", -1)
					wordlist := generateWordlist(alias)
					for i := range wordlist {
						wg.Add(1)
						go func(alias string) {
							checkSpace(alias)
							wg.Done()
						}(wordlist[i])
					}
					wg.Wait()
				case "99":
					info("\nExiting...")
					os.Exit(0)
				default:
					faliure("Wrong option")
				}
		case "99":
			info("\nExiting...")
			os.Exit(0)
		default:
			faliure("Wrong option")
		}

	}
}
