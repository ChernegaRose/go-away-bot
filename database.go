package main

import (
	"database/sql"
)

type XMessage struct {
	From int64  `json:"from,omitempty"`
	Post string `json:"post,omitempty"`
}

type Settings struct {
	admin   int64
	token   string
	channel int64
}

type Member struct {
	id   int64
	from int64
	date string
	post string
}

type Members map[int64]Member

type InlineQuery struct {
	id       string
	sender   int64
	name     string
	chatType string
}

type InlineQueries map[string]InlineQuery

func (s *Settings) load(db *sql.DB) error {
	return db.QueryRow(`SELECT admin, token, channel FROM settings`).Scan(&s.admin, &s.token, &s.channel)
}

func (s *Settings) save(db *sql.DB) error {
	_, err := db.Exec(`UPDATE settings SET admin=?, token=?, channel=?`, s.admin, s.token, s.channel)
	return err
}

func (m *Members) insert(db *sql.DB, mb Member) error {
	_, err := db.Exec(`INSERT INTO member (id, invited, ts, post) VALUES (?, ?, ?, ?)`, mb.id, mb.from, mb.date, mb.post)
	(*m)[mb.id] = mb
	return err
}

func (m *Members) remove(db *sql.DB, mb Member) error {
	_, err := db.Exec(`DELETE FROM member WHERE id=?`, mb.id)
	delete(*m, mb.id)
	return err
}

func (m *Members) removeID(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM member WHERE id=?`, id)
	delete(*m, id)
	return err
}

func (m *Members) load(db *sql.DB) error {
	rows, err := db.Query(`SELECT id, invited, ts, post FROM member`)
	if err != nil {
		return err
	}

	for rows.Next() {
		p := Member{}
		err := rows.Scan(&p.id, &p.from, &p.date, &p.post)
		if err != nil {
			return err
		}
		(*m)[p.id] = p
	}
	return rows.Close()
}

func (iq *InlineQueries) insert(db *sql.DB, q InlineQuery) error {
	_, err := db.Exec(`INSERT INTO inline (id, sender, name, type, ts) VALUES (?, ?, ?, ?, datetime())`,
		q.id, q.sender, q.name, q.chatType)
	(*iq)[q.id] = q
	return err
}
func (iq *InlineQueries) load(db *sql.DB) error {
	rows, err := db.Query(`SELECT id, sender, name, type FROM inline`)
	if err != nil {
		return err
	}

	for rows.Next() {
		p := InlineQuery{}
		err := rows.Scan(&p.id, &p.sender, &p.name, &p.chatType)
		if err != nil {
			return err
		}
		(*iq)[p.id] = p
	}
	return rows.Close()
}
