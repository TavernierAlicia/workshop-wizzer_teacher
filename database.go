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

func getLevels() (list []string, err error) {
	db := dbConnect()

	err = db.Select(&list, "SELECT name FROM levels")
	if err != nil {
		printErr("get data", "getLevels", err)
	}
	return list, err
}

func getLanguages() (list []string, err error) {
	db := dbConnect()

	err = db.Select(&list, "SELECT name FROM languages")
	if err != nil {
		printErr("get data", "getLanguages", err)
	}
	return list, err
}

func getMail(mail string) (result string) {
	db := dbConnect()

	_ = db.QueryRow("SELECT mail FROM users WHERE mail = ?", mail).Scan(&result)
	return result
}

func RecordUser(subForm Sub) (err error) {
	db := dbConnect()

	subForm.Pwd, _ = encodePWD(subForm.Pwd)

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
	mail = getMail(mail)
	if mail != "" {
		err = fmt.Errorf("mail already exists")
		printErr("This mail already exists", "RecordUser", err)
		return err
	}

	_, err = db.Exec("INSERT INTO users (name, surname, mail, repo, type, campus_id, studies_id, matter_id, pwd) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", subForm.Name, subForm.Surname, subForm.Mail, subForm.Repo, subForm.AccountType, school_id, studies_id, matter_id, subForm.Pwd)

	if err != nil {
		printErr("User insertion failed", "RecordUser", err)
	}
	return err

}

func updatePWD(pwd string, id int64) (err error) {
	db := dbConnect()
	fmt.Println(pwd)
	pwd, _ = encodePWD(pwd)
	fmt.Println(pwd)
	_, err = db.Exec("UPDATE users SET pwd = ? WHERE id = ?", pwd, id)

	fmt.Println(err)
	return err
}

func insertToken(token string, mail string) (err error) {
	db := dbConnect()

	_, err = db.Exec("UPDATE users SET token = ? WHERE mail = ?", token, mail)

	return err
}

func deleteToken(identifier string) (err error) {
	db := dbConnect()

	_, err = db.Exec("UPDATE users SET token = '' WHERE mail = ? OR token = ?", identifier, identifier)
	return err
}

func checkToken(token string) (id int64, err error) {
	db := dbConnect()
	err = db.QueryRow("SELECT id FROM users WHERE token = ?", token).Scan(&id)

	if err != nil {
		printErr("get user from token", "checkToken", err)
		return 0, err
	}

	return id, err
}

func getConnected(mail string, pwd string) (token string, err error) {
	db := dbConnect()

	var id int64
	err = db.QueryRow("SELECT id FROM users WHERE mail = ? AND pwd = ?", mail, pwd).Scan(&id)

	if err != nil || id == 0 {
		printErr("Cannot get this user", "getConnected", err)
		return "", err
	}

	token = tokenGenerator()
	_, err = db.Exec("UPDATE users SET token = ? WHERE id = ?", token, id)

	if err != nil {
		printErr("Cannot update token", "getConnected", err)
		return "", err
	}
	return token, err
}

func getUserInfos(token string) (infos UserInfos, err error) {

	db := dbConnect()

	err = db.Get(&infos, "SELECT users.id, users.name, surname, mail, type, IFNULL(repo, '') AS repo, schools.name AS school, users.campus_id, IFNULL(matters.name, '') AS matter, users.matter_id, IFNULL(studies.name, '') AS study, users.studies_id, IFNULL(pic, '') AS pic FROM users LEFT JOIN schools ON users.campus_id = schools.id LEFT JOIN matters ON users.matter_id = matters.id LEFT JOIN studies ON users.studies_id = studies.id WHERE users.token = ?", token)

	return infos, err
}

func getExercices(id int64, aType string, studies_id string, campus_id string, matter_id string, params exoSearch) (exos []*Exos, err error) {

	db := dbConnect()

	// stmt, err := db.Prepare()

	if aType == "student" {
		// get level_id first
		var level_id int64
		var repo string
		err = db.QueryRow("SELECT level_id, matter_id FROM studies WHERE id = ?", studies_id).Scan(&level_id, &matter_id)
		if err != nil {
			printErr("get level_id & matter_id", "getExercices", err)
		}

		// get repo
		err = db.QueryRow("SELECT repo FROM users WHERE id = ?", id).Scan(&repo)
		if err != nil {
			printErr("get repo", "getExercices", err)
		}

		// now get exercices
		err = db.Select(&exos, "SELECT exercices.level_id, exercices.id, exercices.name AS name, CONCAT(?, exercices.git_path) AS git_path, exercices.due_at, exercices.description, matters.name AS matter, 0 AS score, languages.name AS language, exercices.bareme, CONCAT(users.name, ' ', users.surname) AS creator, exercices.created FROM exercices LEFT JOIN matters ON exercices.matter_id = matters.id LEFT JOIN languages ON languages.id = exercices.language_id LEFT JOIN users on users.id = exercices.user_id WHERE CAST(exercices.due_at AS DATE) = CAST(NOW() AS DATE) AND users.campus_id = ? AND exercices.matter_id = ? AND level_id = ?", repo, campus_id, matter_id, level_id)
	} else {

		err = db.Select(&exos, `
		SELECT levels.name AS level, exercices.id, exercices.name AS name, exercices.git_path, CASE WHEN CAST(exercices.due_at AS DATE) = CAST(NOW() AS DATE) THEN "Aujourd'hui" ELSE exercices.due_at END AS due_at, exercices.description, matters.name AS matter, 0 AS score, languages.name AS language, exercices.bareme, CONCAT(users.name, ' ', users.surname) AS creator, exercices.created FROM exercices LEFT JOIN matters ON exercices.matter_id = matters.id LEFT JOIN languages ON languages.id = exercices.language_id LEFT JOIN users on users.id = exercices.user_id LEFT JOIN levels ON exercices.level_id = levels.id WHERE users.campus_id = ? AND exercices.matter_id = ? AND YEAR(exercices.due_at) = YEAR(NOW()) 
		AND IF(? != "", DATE(exercices.due_at) = DATE(?), 1) 
		AND IF(? != "", exercices.name LIKE ?, 1)
		AND IF(? != "", levels.name = ?, 1)
		AND IF(? != "", languages.name = ?, 1) ORDER BY exercices.due_at ASC`, campus_id, matter_id, params.Date, params.Date, params.Name, "%"+params.Name+"%", params.Level, params.Level, params.Language, params.Language)
	}

	if err != nil {
		printErr("get exercices", "getExercices", err)
		return nil, err
	}
	return exos, err

}

func getScore(id int64) (score int64) {
	db := dbConnect()

	err = db.QueryRow("SELECT IFNULL(SUM(score), 0) FROM rendus WHERE student_id = ?", id).Scan(&score)
	if err != nil {
		printErr("get score", "getScore", err)
		return 0
	}

	return score
}
