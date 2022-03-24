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
	Subjects   []string
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
	Subject     string `db:"subject"`
	Bareme      string `db:"bareme"`
	Level       string `db:"level"`
	Creator     string `db:"creator"`
	Created     string `db:"created"`
}

type exoSearch struct {
	Name    string
	Date    string
	Level   string
	Subject string
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
	MoyS  int64  `db:"moy_student"`
	MoyB  int64  `db:"moy_boot"`
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

// export csv data
type Export struct {
	Infos  User
	Exos   []Exercises
	Grades []Grades
}

type User struct {
	Id          int64
	Name        string
	Surname     string
	Mail        string
	Repo        string
	Type        string
	CampusID    int64
	CampusName  string
	StudiesID   int64
	StudiesName string
	MatterID    int64
	MatterName  string
	Pic         string
	Level       string
	Created     string `db:"created"`
	Modified    string `db:"modified"`
}

type Exercises struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Path        string `db:"git_path"`
	Due         string `db:"due_at"`
	Description string `db:"description"`
	Creator     string `db:"creator"`
	LevelID     int64  `db:"level_id"`
	LevelName   string `db:"level_name"`
	MatterID    int64  `db:"matter_id"`
	MatterName  string `db:"matter_name"`
	SubjectID   int64  `db:"subject_id"`
	SubjectName string `db:"subject_name"`
	Bareme      int64  `db:"bareme"`
	Created     string `db:"created"`
	Modified    string `db:"modified"`
	Rendus      []Grades
}

// notes rendues par le prof ou notes reçues par l'élève
type Grades struct {
	Id          int64  `db:"id"`
	ExerciceID  int64  `db:"exercice_id"`
	StudentID   int64  `db:"student_id"`
	StudentName string `db:"student_name"`
	Score       int64  `db:"score"`
	Created     string `db:"created"`
	ExoDetails  Exercises
}
