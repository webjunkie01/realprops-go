package main

import (
    "encoding/json"
    "fmt"
    _ "html"
    "log"
    "net/http"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
  db_user = ""
  db_pass = ""
  db_name = ""
  
)

var db = InitDB()

func InitDB() *sql.DB {
  dblogin := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", db_user, db_pass, db_name)
  db, err := sql.Open("postgres", dblogin)
  checkErr(err)
  //defer db.Close()
  return db
}


func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request ) {
      fmt.Println("Query DB")
      var rows *sql.Rows
      rows, err := db.Query("Select postalcode as value, postalcode as label from properties where postalcode != '' group by postalcode")
      checkErr(err)
      columns, err := rows.Columns()
      checkErr(err)
      count := len(columns)
      tableData := make([]map[string]interface{}, 0)
      values := make([]interface{}, count)
      valuePtrs := make([]interface{}, count)
  
      for rows.Next(){
        for i := 0; i< count; i++ {
          valuePtrs[i] = &values[i]
        }
        rows.Scan(valuePtrs...)
        entry := make(map[string]interface{})
        for i,col := range columns{
          var v interface{}
          val := values[i]
          b,ok := val.([]byte)
          if ok{
            v = string(b)
          }else{
            v = val
          }
          entry[col] = v
        }
        tableData = append(tableData, entry)
      }
      js,err := json.Marshal(tableData)
      checkErr(err)
      w.Header().Set("Content-Type", "application/json")
      w.Write(js)
      //fmt.Fprintf(w, string(tableData))
    })

    log.Fatal(http.ListenAndServe(":8080", nil))

}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
