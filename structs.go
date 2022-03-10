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
