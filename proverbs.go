package main

import (
  "os"
  "fmt"
  "io/ioutil"
  "math/rand"
  "strings"
  "flag"
  "net/http"
  "io"
)

const (
  lineLimit int = 40
)

// currently redundant, kept it in just in case
func filterTextFiles(files *[]os.FileInfo) []os.FileInfo {
  var textFiles []os.FileInfo
  for _, file := range *files {
    if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
      textFiles = append(textFiles, file)
    }
  }
  return textFiles
}

func fetchProverbs() []byte {
  response, err := http.Get("https://ccynth1a.github.io/proverbs.html")
  // error handling
  if err != nil {
    // first error, so now lets loop and see if we can try again. if not, assume that the user does not have an internet connection
    for i := 0;; i++{
      if i > 10 {
        fmt.Println("Error: Could Not Resolve Hostname")
        os.Exit(1)
      }
      response, err = http.Get("https://ccynth1a.github.io/proverbs.html")
      if err == nil {
        break
      }
    }
  }
  defer response.Body.Close()
  if response.StatusCode != http.StatusOK {
    fmt.Println("Error: Webpage Did Not Return HTTP OK")
    os.Exit(1)
  }
  body, err := io.ReadAll(response.Body) // read HTTP response into variable
  return body
}

func fFetchProverbs(filename string) []byte {
  file, err := ioutil.ReadFile(filename)
  if err != nil {
    fmt.Println("Error: File Not Found") 
    os.Exit(1)
  }
  return file
}

// if a proverb ends up longer than 20 characters, insert a newline in
func insertNewLines(proverb string) string {
  var formattedString strings.Builder // OH MY GOD THIS IS SO USEFUL AFTER COMING FROM C THANK YOU LORD
  var count int // count the amount of characters since last newline
  for _, c := range proverb {
    formattedString.WriteRune(c) // write the char to the string
    if count >= lineLimit && c == ' '{ // put a newline char in if it exceeds limit, and importantly the last character was a space
      formattedString.WriteString("\n") 
      count = 0
    }
    count++
  }
  return formattedString.String()
}

/*
EXAMPLE FORMATTING FOR A PROVERB ON THE WEBPAGE/FILE
Proverb: A jade stone is useless before it is processed; a man is good-for-nothing until he is educated.
*/
func extractProverbs(body string) []string { 
  var proverbs []string // will hold the list of strings parsed
  lines := strings.Split(body, "\n") // split the content by newlines
  // iterate through every line 
  for _, line := range lines {
    // first, check for the start of a new proverb. if we hit this if statement, we know we've successfully parsed a completed one
    if strings.Contains(line, "Proverb:") {
      proverb := strings.TrimSpace(strings.TrimPrefix(line, "Proverb:"))
      proverb = insertNewLines(proverb)
      proverbs = append(proverbs, proverb)
    }
  }
  return proverbs
}

func main() {
  localFileName := flag.String("l", "", "Use Local Text Files")
  flag.Parse() // parse flags until you hit a non-flag argument

  var body []byte
  if *localFileName != "" {
    body = fFetchProverbs(*localFileName)
  } else {
    body = fetchProverbs()
  }
  proverbs := extractProverbs(string(body))
  randomIndex := rand.Intn(len(proverbs))
  fmt.Println("----------------------------------------")
  fmt.Println(proverbs[randomIndex])
  fmt.Println("----------------------------------------")
}
