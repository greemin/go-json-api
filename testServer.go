package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "html"
  "io/ioutil"
  "log"
  "net/http"
  "time"
)

func check (e error) {
  if e != nil {
    panic(e)
  }
}

type Comment struct {
  ID string `json:"id`
  UserName string `json:"username"`
  Text string `json:"text"`
  Date time.Time `json:"date"`
  VoteScore int `json:"votescore"`
}

type Response struct {
  Success bool `json:"success"`
}

func loadCommentsFromFile(path string) []Comment {
  file, fileError := ioutil.ReadFile(path)
  check(fileError)

  data := []Comment{}
  jsonError := json.Unmarshal([]byte(file), &data)
  check(jsonError)

  return data
}

func writeCommentsToFile(path string, comments []Comment) bool {
  json, jsonError := json.Marshal(comments)
  check(jsonError)

  fileError := ioutil.WriteFile("comments.json", json, 0644)
  check(fileError)

  return true
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET": {
      file, fileError := ioutil.ReadFile("comments.json")
      check(fileError)
      w.Write(file)
    }
    case "POST": {
      d := json.NewDecoder(r.Body)
      newComment := &Comment{}

      decodeError := d.Decode(newComment)
      check(decodeError)
      newComment.Date = time.Now()
      newComment.ID = newComment.Date.String() + newComment.UserName

      allComments := loadCommentsFromFile("comments.json")
      allComments = append(allComments, *newComment)

      success := writeCommentsToFile("comments.json", allComments)
      if success == true {
        resultJSON, _ := json.Marshal(newComment)
        w.Write(resultJSON)
      } 
    }
    case "DELETE": {
      response := &Response{ Success: false }
      query := r.URL.Query()
      id := query.Get("id")
      fmt.Println("/comment DELETE QUERY", query.Get("id"))

      allComments := loadCommentsFromFile("comments.json")
      for i := 0; i < len(allComments);i++ {
        if allComments[i].ID == id {
          tempVar := allComments[len(allComments)-1]
          allComments[len(allComments)-1] = allComments[i]
          allComments[i] = tempVar
          allComments = allComments[(len(allComments)-1):]

          success := writeCommentsToFile("comments.json", allComments)
          response.Success = success
          break
        }
      }

      responseJSON, _ := json.Marshal(response)
      w.Write(responseJSON)
    }
  }
}

func main() {
  portPtr := flag.String("port", "8080", "port number\":\"")
  flag.Parse()

  fmt.Println("Listening on port", *portPtr)

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	  fmt.Println("Hello,", html.EscapeString(r.URL.Path))
  })

  http.HandleFunc("/comment", commentHandler);

  log.Fatal(http.ListenAndServe(":" + *portPtr, nil))
}
