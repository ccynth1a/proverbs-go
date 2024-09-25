package main

import (
  "os"
  "fmt"
  "bufio"
  "io/ioutil"
  "math/rand"
  "strings"
)

func filterTextFiles(files *[]os.FileInfo) []os.FileInfo {
  var textFiles []os.FileInfo
  for _, file := range *files {
    if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
      textFiles = append(textFiles, file)
    }
  }
  return textFiles
}

func main() {
  files, _ := ioutil.ReadDir(".")

  // check for all text files 
  textFiles := filterTextFiles(&files)
  if len(textFiles) < 1 {
    fmt.Println("Bad Usage: Directory contains no files")
    return
  }

  // get a random file
  randomIndex := rand.Intn(len(textFiles))
  filePath := textFiles[randomIndex].Name()
  
  file, err := os.Open(filePath)
  if err != nil {
    fmt.Printf("Error Opening: %s", filePath)
    return
  }
  defer file.Close() // this keyword is so cool

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    fmt.Println(scanner.Text())
  } 
}
