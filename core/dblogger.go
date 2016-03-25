package core

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	create_table_log string = "create table if not exists msg (id integer not null primary key, utime integer, username text, msgtext text)"
	log_db_path      string = "./log.db"
)

type Dblogger struct {
	db *sql.DB
}

func NewDblogger() *Dblogger {

	dblogger := new(Dblogger)

	db, err := sql.Open("sqlite3", log_db_path)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := create_table_log
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	dblogger.db = db
	return dblogger
}

func (dblogg *Dblogger) SaveMsg(msg Message) {

	tx, err := dblogg.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(fmt.Sprintf("insert into msg(utime, username, msgtext) values(%v, '%v', '%v')", msg.Time.Unix(), msg.Username, msg.Text))
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()
	tx.Commit()
}

func (dblogg *Dblogger) FindAllMsg() []Message {

	var msgs []Message

	rows, err := dblogg.db.Query("select id, utime, username, msgtext from msg")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var utime int64
		var username string
		var msgtext string
		rows.Scan(&id, &utime, &username, &msgtext)
		msgs = append(msgs, Message{Username: username, Text: msgtext, Time: time.Unix(utime, 0)})
	}

	return msgs
}
