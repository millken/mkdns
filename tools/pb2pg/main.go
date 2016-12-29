package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"

	// pgsql
	_ "github.com/lib/pq"
)

func main() {
	var err error
	in := flag.String("in", "test.pb", "json file path")
	domain := flag.String("domain", "test.com", "domain name")
	out := flag.String("pg", "postgres://postgres:admin@127.0.0.1/dns?sslmode=disable", "postgres config")
	flag.Parse()
	content, err := ioutil.ReadFile(*in)
	if err != nil {
		log.Fatalln("read json file err:", err)
	}
	db, err := sql.Open("postgres", *out)
	if err != nil {
		log.Fatalf("[%s]: %s", *out, err)
	}
	defer db.Close()

	sStmt := "insert into config.record (domain, value) values ($1, $2)"

	stmt, err := db.Prepare(sStmt)
	if err != nil {
		log.Fatal("prepare err", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(domain, content)
	if err != nil || res == nil {
		log.Fatalf("exec :%s %s %s", domain, content, err)
	}
	log.Println("done")
}
