package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func dbConnect() *sqlx.DB {

	//// IMPORT CONFIG ////
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	printErr("reading config file", "dbConnect", err)

	//// DB CONNECTION ////
	pathSQL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", viper.GetString("database.user"), viper.GetString("database.pass"), viper.GetString("database.host"), viper.GetInt("database.port"), viper.GetString("database.dbname"))
	db, err := sqlx.Connect("mysql", pathSQL)

	printErr("connect to database", "dbConnect", err)
	return db
}

func getSchools() (list []string, err error) {
	db := dbConnect()

	err = db.Select(&list, "SELECT name FROM schools")
	if err != nil {
		printErr("Unable to get data", "GetSchools", err)
	}
	return list, err
}

func getStudies() (list []string, err error) {
	db := dbConnect()

	err = db.Select(&list, "SELECT name FROM studies")
	if err != nil {
		printErr("Unable to get data", "GetStudies", err)
	}
	return list, err
}

func getMatters() (list []string, err error) {
	db := dbConnect()

	err = db.Select(&list, "SELECT name FROM matters")
	if err != nil {
		printErr("Unable to get data", "GetMatters", err)
	}
	return list, err
}

func RecordUser(subForm Sub) (err error) {
	db := dbConnect()

	// check school
	var (
		school_id  int
		studies_id int
		matter_id  int
		mail       string
	)
	_ = db.QueryRow("SELECT id FROM schools WHERE name = ?", subForm.Campus).Scan(&school_id)

	if school_id == 0 {
		err = fmt.Errorf("this id doesn't exists")
		printErr("Unable to get school id", "RecordUser", err)
		return err
	}

	// check if student or else
	if subForm.AccountType == "student" {
		// check studies
		_ = db.QueryRow("SELECT id FROM studies WHERE name = ?", subForm.Studies).Scan(&studies_id)
		if studies_id == 0 {
			err = fmt.Errorf("this id doesn't exists")
			printErr("Unable to get studies id", "RecordUser", err)
			return err
		}
	} else {
		// check matter
		_ = db.QueryRow("SELECT id FROM matters WHERE name = ?", subForm.Matter).Scan(&matter_id)
		if matter_id == 0 {
			err = fmt.Errorf("this id doesn't exists")
			printErr("Unable to get matter id", "RecordUser", err)
			return err
		}
	}

	// check mail exists
	_ = db.QueryRow("SELECT mail FROM users WHERE mail = ?", subForm.Mail).Scan(&mail)
	fmt.Println(mail)
	if mail != "" {
		err = fmt.Errorf("mail already exists")
		printErr("This mail already exists", "RecordUser", err)
		return err
	}

	_, err = db.Exec("INSERT INTO users (name, surname, mail, repo, campus_id, studies_id, matter_id, pwd) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", subForm.Name, subForm.Surname, subForm.Mail, subForm.Repo, school_id, studies_id, matter_id, subForm.Pwd)

	if err != nil {
		printErr("User insertion failed", "RecordUser", err)
	}
	return err

}
