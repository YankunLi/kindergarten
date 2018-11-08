package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sort"
	"strconv"
)

type result struct {
	Title   string `json:"title"`
	Project string `json:"project"`
	User_id string `json:"user_id"`
	Answers map[string]answer
}

type answer struct {
	Type   string `json:"type"`
	Result string `json:"result"`
}

func decode(data []byte) string {

	var ret result
	var tempStr string
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return ""
	}

	tempStr = ret.Title + "\t" + ret.Project + "\t" + ret.User_id

	var kv = make(map[int]string)
	for k, v := range ret.Answers {
		num, _ := strconv.Atoi(k)
		kv[num] = v.Result
	}

	var keys []int

	for k := range kv {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for _, k := range keys {
		str := kv[k]
		if str == "" {
			str = "*X*"
		}
		tempStr = tempStr + "\t" + str
	}

	return tempStr

}

func main() {
	//parse parmater
	var projectP *string = flag.String("project", "", "project name")
	flag.Parse()
	if *projectP == "" {
		fmt.Printf("not found project %s", *projectP)
		panic("not found project")
	}

	projectName := *projectP
	outFilePath := "out_" + projectName
	//connect db
	db, err := sql.Open("mysql", "user:password@tcp(IP:3306)/questionnaire?charset=utf8")
	if err != nil {
		fmt.Println("open db fail")
		panic("fail")
	}
	defer db.Close()
	stmt, _ := db.Prepare(`select answers from answers where project_name = ?`)

	rows, err := stmt.Query(projectName)
	if err != nil {
		fmt.Printf("execute sql fail")
		panic("execute sql fail")
	}
	defer rows.Close()
	var outdata string

	fd, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer fd.Close()
	for rows.Next() {
		rows.Scan(&outdata)
		fd.WriteString(decode([]byte(outdata)) + "\n")
	}

}
