package main

import "github.com/gorilla/sessions"

type Sub struct {
	// General
	AccountType string
	Name        string
	Surname     string
	Mail        string
	Repo        string

	// School
	Campus  string
	Studies string
	Matter  string

	// Secu
	Pwd        string
	PwdConfirm string

	// others
	Pic string
}

type SubFormInfos struct {
	Schools    []string
	Formations []string
	Matters    []string
	Levels     []string
	Languages  []string
}

type SessionInfos struct {
	Session    sessions.Session
	Token      string
	Atype      string
	Name       string
	Surname    string
	Campus_id  string
	Matter_id  string
	Studies_id string
}

type UserInfos struct {
	Id         int64  `db:"id"`
	Name       string `db:"name"`
	Surname    string `db:"surname"`
	Mail       string `db:"mail"`
	Type       string `db:"type"`
	Repository string `db:"repo"`
	Campus     string `db:"school"`
	CampusID   int64  `db:"campus_id"`
	Studies    string `db:"study"`
	StudiesID  int64  `db:"studies_id"`
	Matter     string `db:"matter"`
	MatterID   int64  `db:"matter_id"`
	Pic        string `db:"pic"`
}

type Exos struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	GitPath     string `db:"git_path"`
	Due         string `db:"due_at"`
	Description string `db:"description"`
	Matter      string `db:"matter"`
	Score       int64  `db:"score"`
	Language    string `db:"language"`
	Bareme      string `db:"bareme"`
	Level       string `db:"level"`
	Creator     string `db:"creator"`
	Created     string `db:"created"`
}

type exoSearch struct {
	Name     string
	Date     string
	Level    string
	Language string
}
