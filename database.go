package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

const DBInitSQLFileName = "init.sql"

func DBInit(db *sql.DB) error {
	file, err := os.ReadFile(DBInitSQLFileName)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(file))
	return err
}

type Settings struct {
	Token    string
	Admin    int64
	IsPublic int
	Featured int64
}

func (s *Settings) load(db *sql.DB) error {
	err := db.QueryRow(`SELECT token, admin, is_public FROM settings`).
		Scan(&s.Token, &s.Admin, &s.IsPublic)
	s.Featured = 1
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("First time configuration:")
		fmt.Println("Token:")
		_, err = fmt.Scanln(&s.Token)
		if err != nil {
			return err
		}
		fmt.Println("Admin ID:")
		_, err = fmt.Scanln(&s.Admin)
		if err != nil {
			return err
		}
		fmt.Println("Is Public (0 - off, 1 - engage + whitelist, 2 - all):")
		_, err = fmt.Scanln(&s.IsPublic)
		if err != nil {
			return err
		}
		_, err = db.Exec(`INSERT INTO settings (token, admin, is_public) VALUES (?, ?, ?)`,
			s.Token, s.Admin, s.IsPublic)
		return err
	}
	return err
}

func (s *Settings) save(db *sql.DB) error {
	_, err := db.Exec(`UPDATE settings SET token=?, admin=?, is_public=?`,
		s.Token, s.Admin, s.IsPublic)
	return err
}

type Query struct {
	QueryID   string
	QueryText string
	ChatType  string
	UserID    int64
	UserName  string
	UserLang  string
}

type Queries map[string]Query

func (qs *Queries) insert(db *sql.DB, q Query) error {
	_, err := db.Exec(
		`INSERT INTO queries (id, query_text, chat_type, user_id, user_name, user_lang, timestamp) VALUES (?, ?, ?, ?, ?, ?, datetime())`,
		q.QueryID, q.QueryText, q.ChatType, q.UserID, q.UserName, q.UserLang)
	(*qs)[q.QueryID] = q
	return err
}

func (qs *Queries) load(db *sql.DB) error {
	*qs = make(Queries)
	rows, err := db.Query(`SELECT id, query_text, chat_type, user_id, user_name, user_lang FROM queries`)
	if err != nil {
		return err
	}

	for rows.Next() {
		q := Query{}
		err := rows.Scan(&q.QueryID, &q.QueryText, &q.ChatType, &q.UserID, &q.UserName, &q.UserLang)
		if err != nil {
			return err
		}
		(*qs)[q.QueryID] = q
	}
	return rows.Close()
}

type Message struct {
	MessageID string
	QueryID   string
	UserID    int64
	UserName  string
}

type Messages map[string]Message

func (ms *Messages) insert(db *sql.DB, m Message) error {
	_, err := db.Exec(
		`INSERT INTO messages (id, query_id, user_id, user_name, timestamp) VALUES (?, ?, ?, ?, datetime())`,
		m.MessageID, m.QueryID, m.QueryID, m.UserName)
	(*ms)[m.MessageID] = m
	return err
}

func (ms *Messages) load(db *sql.DB) error {
	*ms = make(Messages)
	rows, err := db.Query(`SELECT id, query_id, user_id, user_name FROM messages`)
	if err != nil {
		return err
	}

	for rows.Next() {
		m := Message{}
		err := rows.Scan(&m.MessageID, &m.QueryID, &m.UserID, &m.UserName)
		if err != nil {
			return err
		}
		(*ms)[m.MessageID] = m
	}
	return rows.Close()
}

type Creator struct {
	CreatorID int64
	UserName  string
	Status    int
}

type Creators map[int64]Creator

func (cs *Creators) insert(db *sql.DB, c Creator) error {
	_, err := db.Exec(`INSERT INTO creators (id, user_name, status) VALUES (?, ?, ?)`,
		c.CreatorID, c.UserName, c.Status)
	(*cs)[c.CreatorID] = c
	return err
}

func (cs *Creators) load(db *sql.DB) error {
	*cs = make(Creators)
	rows, err := db.Query(`SELECT id, user_name, status FROM creators`)
	if err != nil {
		return err
	}

	for rows.Next() {
		c := Creator{}
		err := rows.Scan(&c.CreatorID, &c.UserName, &c.Status)
		if err != nil {
			return err
		}
		(*cs)[c.CreatorID] = c
	}
	return rows.Close()
}

type Member struct {
	UserID       int64
	UserName     string
	UserLang     string
	FromID       int64
	MessageID    string
	ChatInstance string
	ContestID    int64
}

type Members map[int64]map[int64]Member //[contest][user]

func (ms *Members) insert(db *sql.DB, m Member) error {
	_, err := db.Exec(
		`INSERT INTO members (user_id, user_name, user_lang, from_id, message_id, chat_instance, contest_id, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, datetime())`,
		m.UserID, m.UserName, m.UserLang, m.FromID, m.MessageID, m.ChatInstance, m.ContestID)
	if (*ms)[m.ContestID] == nil {
		(*ms)[m.ContestID] = make(map[int64]Member)
		(*ms)[m.ContestID][m.UserID] = m
	} else {
		(*ms)[m.ContestID][m.UserID] = m
	}
	return err
}

func (ms *Members) load(db *sql.DB) error {
	*ms = make(Members)
	rows, err := db.Query(`SELECT user_id, user_name, user_lang, from_id, message_id, chat_instance, contest_id FROM members`)
	if err != nil {
		return err
	}

	for rows.Next() {
		m := Member{}
		err := rows.Scan(&m.UserID, &m.UserName, &m.UserLang, &m.FromID, &m.MessageID, &m.ChatInstance, &m.ContestID)
		if err != nil {
			return err
		}
		if (*ms)[m.ContestID] == nil {
			(*ms)[m.ContestID] = make(map[int64]Member)
			(*ms)[m.ContestID][m.UserID] = m
		} else {
			(*ms)[m.ContestID][m.UserID] = m
		}
	}
	return rows.Close()
}

type Post struct {
	ContestID   int64
	Type        string
	Title       string
	Message     string
	Description string
	Image       string
}

type Posts map[int64]map[string]Post //[contest][type]

func (ps *Posts) insert(db *sql.DB, p Post) error {
	_, err := db.Exec(
		`INSERT INTO posts (contest_id, type, title, message, description, image) VALUES (?, ?, ?, ?, ?, ?)`,
		p.ContestID, p.Type, p.Title, p.Message, p.Description, p.Image)
	if (*ps)[p.ContestID] == nil {
		(*ps)[p.ContestID] = make(map[string]Post)
		(*ps)[p.ContestID][p.Type] = p
	} else {
		(*ps)[p.ContestID][p.Type] = p
	}
	return err
}

func (ps *Posts) load(db *sql.DB) error {
	*ps = make(Posts)
	rows, err := db.Query(`SELECT contest_id, type, title, message, description, image FROM posts`)
	if err != nil {
		return err
	}

	for rows.Next() {
		p := Post{}
		err := rows.Scan(&p.ContestID, &p.Type, &p.Title, &p.Message, &p.Description, &p.Image)
		if err != nil {
			return err
		}
		if (*ps)[p.ContestID] == nil {
			(*ps)[p.ContestID] = make(map[string]Post)
			(*ps)[p.ContestID][p.Type] = p
		} else {
			(*ps)[p.ContestID][p.Type] = p
		}
	}
	return rows.Close()
}

type Contest struct {
	ContestID     int64
	CreatorID     int64
	ContestName   string
	ContestStart  string
	ContestEnd    string
	ContestActive int
	Username      string
}

type Contests map[int64]Contest

func (cs *Contests) insert(db *sql.DB, c Contest) error {
	res, err := db.Exec(
		`INSERT INTO contests (creator_id, contest_name, contest_start, contest_end, contest_active, username, timestamp) VALUES (?, ?, ?, ?, ?, ?, datetime())`,
		c.CreatorID, c.ContestName, c.ContestStart, c.ContestEnd, c.ContestActive, c.Username)
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ContestID = id
	(*cs)[c.ContestID] = c
	return err
}

func (cs *Contests) load(db *sql.DB) error {
	*cs = make(Contests)
	rows, err := db.Query(`SELECT id, creator_id, contest_name, contest_start, contest_end, contest_active, username FROM contests`)
	if err != nil {
		return err
	}

	for rows.Next() {
		c := Contest{}
		err := rows.Scan(&c.ContestID, &c.CreatorID, &c.ContestName, &c.ContestStart, &c.ContestEnd, &c.ContestActive, &c.Username)
		if err != nil {
			return err
		}
		(*cs)[c.ContestID] = c
	}
	return rows.Close()
}

type Params struct {
	From    int64 `json:"fr,omitempty"`
	Contest int64 `json:"ct,omitempty"`
}
