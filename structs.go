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

type studentRank struct {
	Id      int64  `db:"id"`
	Name    string `db:"name"`
	Surname string `db:"surname"`
	Studies string `db:"studies"`
	Level   string `db:"level"`
	Score   int64  `db:"score"`
	Rank    int64  `db:"rank"`
}

// for student/id and histo
type Student struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Surname     string `db:"surname"`
	Studies     string `db:"studies"`
	Rank        string `db:"rank"`
	Score       int64  `db:"score"`
	ScoreByLang []ScoreLang
	DaysDetails []Days
}

type ScoreLang struct {
	Lang  string `db:"lang"`
	Score int64  `db:"score_by_lang"`
	Moy   int64  `db:"moy_score"`
}

type Days struct {
	Score int64  `db:"score_by_day"`
	Date  string `db:"date"`
	Exos  []Rendus
}

type Rendus struct {
	Name  string `db:"exo_name"`
	Repo  string `db:"repo"`
	Lang  string `db:"exo_lang"`
	Score int64  `db:"exo_score"`
	Total int64  `db:"exo_total"`
}

type NewGrade struct {
	ExerciceID int `json:"exercice_id"`
	StudentID  int `json:"student_id"`
	Score      int `json:"score"`
}

type AllResults struct {
	ScoreByLang []LangScore
	DaysDetails []Days
}

type LangScore struct {
	Lang        string `db:"lang"`
	Moy         int64  `db:"moy_score"`
	TotalPoints int64
}

type OverviewSearch struct {
	Level   string
	Studies string
}
