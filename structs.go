package main

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
}

type UserInfos struct {
	Id         int64  `db:"id"`
	Name       string `db:"name"`
	Surname    string `db:"surname"`
	Mail       string `db:"mail"`
	Type       string `db:"type"`
	Repository string `db:"repo"`
	Campus     string `db:"school"`
	Studies    string `db:"study"`
	Matter     string `db:"matter"`
	Pic        string `db:"pic"`
}
