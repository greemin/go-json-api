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
  Id string `json:"id`
  UserName string `json:"username"`
  Text string `json:"text"`
  Date time.Time `json:"date"`
  VoteScore int `json:"votescore"`
}

type Response struct {
  Success bool `json:"success"`
  Count int `json:count"`
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
  fmt.Println(r.Method, r.URL.Path, r.RemoteAddr)
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
      newComment.Id = newComment.Date.Format(time.RFC3339) + "-" + newComment.UserName

      allComments := loadCommentsFromFile("comments.json")
      allComments = append(allComments, *newComment)

      success := writeCommentsToFile("comments.json", allComments)
      if success == true {
        resultJSON, _ := json.Marshal(newComment)
        fmt.Println("New Comment: " + newComment.Text, "(" + newComment.Id + ")")
        w.Write(resultJSON)
      } 
    }
    case "DELETE": {
      response := &Response{ Success: false, Count: 0 }
      query := r.URL.Query()
      id := query.Get("id")
      fmt.Println("/comment DELETE QUERY", query.Get("id"))

      allComments := loadCommentsFromFile("comments.json")
      count := 0
      for i := 0; i < len(allComments);i++ {
        if allComments[i].Id == id {
          count++

          if len(allComments) > 1 {
            tempVar := allComments[len(allComments)-1]
            allComments[i] = tempVar
            allComments[len(allComments)-1] = allComments[i]
            allComments = allComments[(len(allComments)-1):]
          } else {
            allComments = []Comment{}
          }

          success := writeCommentsToFile("comments.json", allComments)
          response.Success = success
          response.Count = count
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
