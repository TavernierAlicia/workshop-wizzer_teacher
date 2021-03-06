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

func getStudiesByMatter(matter int64) (list []string, err error) {
	db := dbConnect()

	err = db.Select(&list, "SELECT name FROM studies WHERE matter_id = ?", matter)
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

func getSubjects(matterId int64) (list []string, err error) {
	db := dbConnect()

	err = db.Select(&list, "SELECT name FROM subjects WHERE matter_id = ?", matterId)
	if err != nil {
		printErr("get data", "getSubjects", err)
	}
	return list, err
}

func getMail(mail string) (result string) {
	db := dbConnect()

	err = db.QueryRow("SELECT mail FROM users WHERE mail = ?", mail).Scan(&result)

	return result
}

func RecordUser(subForm Sub) (err error) {
	db := dbConnect()

	subForm.Pwd = encodePWD(subForm.Pwd)

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
	mail = getMail(subForm.Mail)
	if mail != "" {
		err = fmt.Errorf("mail already exists")
		printErr("This mail already exists", "RecordUser", err)
		return err
	}

	id, err := db.Exec("INSERT INTO users (name, surname, mail, repo, type, campus_id, studies_id, matter_id, pwd) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", subForm.Name, subForm.Surname, subForm.Mail, subForm.Repo, subForm.AccountType, school_id, studies_id, matter_id, subForm.Pwd)

	if err != nil {
		printErr("User insertion failed", "RecordUser", err)
		return err
	}

	if subForm.AccountType == "prof" {
		botToken := tokenGenerator()
		u_id, err := id.LastInsertId()
		if err != nil {
			printErr("get lastrowid", "RecordUser", err)
			return err
		}

		// add botToken
		_, err = db.Exec("INSERT INTO botToken (user_id, token) VALUES (?, ?)", u_id, botToken)

		if err != nil {
			printErr("insert bot token", "RecordUser", err)
			return err
		}
	}
	return err
}

func updatePWD(pwd string, id int64) (err error) {
	db := dbConnect()
	pwd = encodePWD(pwd)
	_, err = db.Exec("UPDATE users SET pwd = ? WHERE id = ?", pwd, id)

	printErr("update pwd", "updatePWD", err)

	return err
}

func updateBotToken(id int64) (botToken string, err error) {
	db := dbConnect()
	token := tokenGenerator()
	_, err = db.Exec("UPDATE botToken SET token = ? WHERE user_id = ?", token, id)

	if err != nil {
		printErr("update token", "updateBotToken", err)
	}

	return token, err
}

func getBotToken(id int64) (botToken string, err error) {
	db := dbConnect()
	err = db.QueryRow("SELECT token FROM botToken WHERE user_id = ?", id).Scan(&botToken)

	if err != nil {
		printErr("get bot token", "getBotToken", err)
	}
	return botToken, err
}

func checkBotToken(token string) (id int64, err error) {
	db := dbConnect()
	err = db.QueryRow("SELECT user_id FROM botToken WHERE token = ?", token).Scan(&id)

	if err != nil {
		printErr("get user from token", "checkToken", err)
		return 0, err
	}

	return id, err
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
		err = db.Select(&exos, `
			SELECT 
				exercices.level_id AS level, 
				exercices.id, 
				exercices.name AS name, 
				CONCAT(?, exercices.git_path) AS git_path, 
				DATE(exercices.due_at) AS due_at, 
				exercices.description, 
				matters.name AS matter, 
				0 AS score, 
				IFNULL(subjects.name, '') AS subject, 
				exercices.bareme, 
				CONCAT(users.name, ' ', users.surname) AS creator, 
				exercices.created 
			FROM exercices 
				JOIN matters ON exercices.matter_id = matters.id 
				JOIN subjects ON subjects.id = exercices.subject_id 
				JOIN users on users.id = exercices.user_id 
			WHERE 
				CAST(exercices.due_at AS DATE) = CAST(NOW() AS DATE) AND users.campus_id = ? 
				AND exercices.matter_id = ? 
				AND level_id = ?`,
			repo, campus_id, matter_id, level_id)
	} else {

		err = db.Select(&exos, `
		SELECT levels.name AS level, exercices.id, exercices.name AS name, exercices.git_path, CASE WHEN CAST(exercices.due_at AS DATE) = CAST(NOW() AS DATE) THEN "Aujourd'hui" ELSE DATE(exercices.due_at) END AS due_at, exercices.description, matters.name AS matter, 0 AS score, IFNULL(subjects.name, "") AS subject, exercices.bareme, CONCAT(users.name, ' ', users.surname) AS creator, exercices.created FROM exercices LEFT JOIN matters ON exercices.matter_id = matters.id LEFT JOIN subjects ON subjects.id = exercices.subject_id LEFT JOIN users on users.id = exercices.user_id LEFT JOIN levels ON exercices.level_id = levels.id WHERE users.campus_id = ? AND exercices.matter_id = ? AND YEAR(exercices.due_at) = YEAR(NOW()) 
		AND IF(? != "", DATE(exercices.due_at) = DATE(?), 1) 
		AND IF(? != "", exercices.name LIKE ?, 1)
		AND IF(? != "", levels.name = ?, 1)
		AND IF(? != "", subjects.name = ?, 1) ORDER BY exercices.due_at ASC`, campus_id, matter_id, params.Date, params.Date, params.Name, "%"+params.Name+"%", params.Level, params.Level, params.Subject, params.Subject)
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

func updateParams(id int64, pic string, repo string, campus string, studies string, matter string) (err error) {
	db := dbConnect()

	deletePic := ""
	err = db.QueryRow("SELECT pic FROM users WHERE id = ?", id).Scan(&deletePic)

	if err != nil {
		printErr("get old pic", "updateParams", err)
	}

	if pic != deletePic && deletePic != viper.GetString("default.default_pic") {
		deleteOldPic(deletePic)
	}

	if repo != "" {
		_, err = db.Exec("UPDATE users SET pic = ?, repo = ?, campus_id = (SELECT id FROM schools WHERE name = ?), studies_id = (SELECT id FROM studies WHERE name = ?) WHERE id = ?", pic, repo, campus, studies, id)
	} else {
		_, err = db.Exec("UPDATE users SET pic = ?, campus_id = (SELECT id FROM schools WHERE name = ?), matter_id = (SELECT id FROM matters WHERE name = ?) WHERE id = ?", pic, campus, matter, id)
	}

	if err != nil {
		printErr("updare params", "updateParams", err)
	}
	return err
}

func postExo(exoName string, gitPath string, date string, desc string, level string, matter int64, subject string, bareme string, user_id int64) (err error) {

	db := dbConnect()

	_, err = db.Exec(`
	INSERT INTO exercices (
		name, 
		git_path, 
		due_at, 
		description, 
		level_id, 
		matter_id, 
		subject_id, 
		bareme, 
		user_id
		) VALUES (
			?, ?, ?, ?, 
			(SELECT id FROM levels WHERE name = ?), 
			?,
			(SELECT id FROM subjects WHERE name = ?),
			 ?, ?
		)`, exoName, gitPath, date, desc, level, matter, subject, bareme, user_id)

	if err != nil {
		printErr("insert exercice", "postExo", err)
	}
	return err

}

func putExo(exoName string, gitPath string, date string, desc string, level string, matter int64, subject string, bareme string, user_id int64, exo_id string) (err error) {

	db := dbConnect()

	_, err = db.Exec(`
	UPDATE exercices SET 
		name = ?, 
		git_path = ?, 
		due_at = ?, 
		description = ?, 
		level_id = (SELECT id FROM levels WHERE name = ?), 
		matter_id = ?, 
		subject_id = (SELECT id FROM subjects WHERE name = ?), 
		bareme = ? 
		WHERE user_id = ? AND id = ?`,
		exoName, gitPath, date, desc, level, matter, subject, bareme, user_id, exo_id)

	if err != nil {
		printErr("update exercice", "putExo", err)
	}
	return err

}

func deleteExo(user_id int64, exo_id string) (err error) {

	db := dbConnect()

	_, err = db.Exec("DELETE FROM exercices WHERE id = ? AND user_id = ?", exo_id, user_id)

	if err != nil {
		printErr("delete exercice", "deleteExo", err)
	}
	return err

}

func getExoDetails(exo_id string, user_id int64) (exo Exos, err error) {
	db := dbConnect()

	err = db.Get(&exo, `
	SELECT levels.name AS level, 
		exercices.id, 
		exercices.name AS name, 
		exercices.git_path AS git_path, 
		DATE(exercices.due_at) AS due_at, 
		exercices.description, 
		matters.name AS matter, 
		0 AS score, 
		subjects.name AS subject, 
		exercices.bareme, 
		CONCAT(users.name, ' ', users.surname) AS creator, 
		exercices.created 
	FROM exercices 
		LEFT JOIN matters ON exercices.matter_id = matters.id 
		LEFT JOIN subjects ON subjects.id = exercices.subject_id 
		LEFT JOIN users ON users.id = exercices.user_id 
		LEFT JOIN levels ON levels.id = exercices.level_id
	WHERE exercices.user_id = ? AND exercices.id = ?`,
		user_id, exo_id)

	if err != nil {
		printErr("get single exercice", "getExoDetails", err)
	}
	return exo, err
}

func getStudents(studies_id int64, campus_id int64, matter_id int64) (students []*studentRank, err error) {
	db := dbConnect()

	err = db.Select(&students, `
	SELECT users.id, 
		users.name, 
		users.surname, 
		studies.name AS studies, 
		levels.name AS level,
		IFNULL(SUM(score), 0) AS score
	FROM users 
		JOIN studies ON studies.id = users.studies_id
		JOIN matters ON matters.id = studies.matter_id
		JOIN levels ON studies.level_id = levels.id
		LEFT JOIN rendus ON rendus.student_id = users.id
	WHERE users.campus_id = ?
		AND (matters.id = ? OR (SELECT studies.matter_id FROM studies WHERE studies.id = ?) = matters.id)
		AND users.type = "student"
	GROUP BY users.id
	ORDER BY score DESC
	`,
		campus_id, matter_id, studies_id)

	if err != nil {
		printErr("get student rank", "getStudents", err)
	}

	return students, err
}

func getStudentScoring(student_id string) (studentScoring Student, err error) {

	db := dbConnect()

	// need some data before filling structs
	campus_id := 0
	matter_id := 0
	studies_id := 0
	level_id := 0

	err = db.QueryRow("SELECT campus_id, studies_id, users.matter_id, studies.level_id AS level_id FROM users JOIN studies ON studies.id = users.studies_id WHERE users.id = ?", student_id).Scan(&campus_id, &studies_id, &matter_id, &level_id)
	if err != nil {
		printErr("cannot get init data", "getStudentScoring", err)
		return studentScoring, err
	}

	// get student sample data
	err = db.Get(&studentScoring, `
		
		SELECT a.* FROM (
			SELECT 
				users.id AS id,
				users.name AS name, 
				@rank:=IFNULL(@rank,0)+1 AS rank,
				users.surname AS surname,
				studies.name AS studies,
				IFNULL(SUM(rendus.score), 0) AS score
				FROM users 
				JOIN studies ON studies.id = users.studies_id
				JOIN matters ON matters.id = studies.matter_id
				JOIN levels ON studies.level_id = levels.id
				LEFT JOIN rendus ON rendus.student_id = users.id
			WHERE users.campus_id = ?
				AND (matters.id = ? OR (SELECT studies.matter_id FROM studies WHERE studies.id = ?) = matters.id)
				AND users.type = "student"
			GROUP BY users.id
			ORDER BY score DESC
		) AS a WHERE a.id = ?
	`, campus_id, matter_id, studies_id, student_id)

	if err != nil {
		printErr("get user sample data", "getStudentScoring", err)
		return studentScoring, err
	}

	// get score by lang
	err = db.Select(&studentScoring.ScoreByLang, `
	SELECT 
		subjects.name AS lang,
		IFNULL(SUM(IF(users.id = ?, rendus.score, 0)), 0) AS score_by_lang,
		IFNULL(Cast(AVG(rendus.score) AS Int), 0) AS moy_boot,
		IFNULL(Cast(AVG(IF(users.id = ?, rendus.score, NULL)) AS Int), 0) AS moy_student
	FROM rendus
		JOIN exercices ON exercices.id = rendus.exercice_id
		JOIN subjects ON exercices.subject_id = subjects.id
		JOIN users ON rendus.student_id = users.id
		JOIN studies ON studies.id = users.studies_id
		JOIN matters ON matters.id = studies.matter_id
		JOIN levels ON studies.level_id = levels.id
		WHERE users.campus_id = ? AND (matters.id = ? OR (SELECT studies.matter_id FROM studies WHERE studies.id = ?) = matters.id) AND levels.id = ?
		GROUP BY subjects.id;
	`, student_id, student_id, campus_id, matter_id, studies_id, level_id)

	if err != nil {
		printErr("get student scoring by lang", "getStudentScoring", err)
		return studentScoring, err
	}

	// get sample days data
	err = db.Select(&studentScoring.DaysDetails, `
		SELECT 
			DATE(created) AS date,
			IFNULL(SUM(score), 0) AS score_by_day
		FROM rendus 
		WHERE student_id = ?
		GROUP BY date
	`, student_id)

	if err != nil {
		printErr("get user sample days data", "getStudentScoring", err)
		return studentScoring, err
	}

	// get exos by day
	for i, day := range studentScoring.DaysDetails {
		err = db.Select(&studentScoring.DaysDetails[i].Exos, `
		SELECT 
			exercices.name AS exo_name,
			exercices.git_path AS repo,
			subjects.name AS exo_lang,
			rendus.score AS exo_score,
			exercices.bareme AS exo_total
		FROM rendus
			JOIN exercices ON exercices.id = rendus.exercice_id
			JOIN subjects ON subjects.id = exercices.subject_id
		WHERE 
			rendus.student_id = ?
			AND DATE(rendus.created) = ?
	`, student_id, day.Date)

		if err != nil {
			printErr("get day exos data", "getStudentScoring", err)
			return studentScoring, err
		}
	}

	return studentScoring, err
}

func insertRendu(exercice NewGrade) (err error) {
	db := dbConnect()
	_, err = db.Exec("INSERT INTO rendus (exercice_id, student_id, score) VALUES (?, ?, ?)", exercice.ExerciceID, exercice.StudentID, exercice.Score)

	if err != nil {
		printErr("insert new rendu", "insertRendu", err)
	}
	return err
}

func getAllStudentScoring(campus_id int64, matter_id int64, params OverviewSearch) (studentScoring AllResults, err error) {

	db := dbConnect()

	// get score by lang
	err = db.Select(&studentScoring.ScoreByLang, `
	SELECT
		subjects.name AS lang,
		IFNULL(Cast(AVG(IF(studies.id IS NULL, NULL, rendus.score)) AS Int),0) AS moy_score
	FROM exercices
		JOIN subjects ON exercices.subject_id = subjects.id
		JOIN users AS u1 ON exercices.user_id = u1.id
		JOIN levels ON levels.id = exercices.level_id
		LEFT JOIN rendus ON exercices.id = rendus.exercice_id
		LEFT JOIN users AS u2 ON u2.id = rendus.student_id
		LEFT JOIN studies ON studies.id = u2.studies_id AND IF(? != "", studies.name LIKE ?, 1)
	WHERE exercices.matter_id = ? AND u1.campus_id = ? 
	AND IF(? != "", levels.name LIKE ?, 1)
	GROUP BY subjects.id;
	`, params.Studies, params.Studies, matter_id, campus_id, params.Level, params.Level)

	if err != nil {
		printErr("get score by lang", "getAllStudentScoring", err)
		return studentScoring, err
	}

	for i, lang := range studentScoring.ScoreByLang {
		err = db.QueryRow(`
			SELECT
				IFNULL(Cast(AVG(exercices.bareme) AS Int), 0)
			FROM exercices
			JOIN subjects ON exercices.subject_id = subjects.id
			JOIN users ON exercices.user_id = users.id
			JOIN levels ON levels.id = exercices.level_id
			WHERE exercices.matter_id = ? AND users.campus_id = ? AND subjects.name = ?
			AND IF(? != "", levels.name LIKE ?, 1)
			GROUP BY subjects.id;
		`, matter_id, campus_id, lang.Lang, params.Level, params.Level).Scan(&studentScoring.ScoreByLang[i].TotalPoints)
		if err != nil {
			printErr("get total points by lang", "getAllStudentScoring", err)
			return studentScoring, err
		}
	}

	// get sample days data
	err = db.Select(&studentScoring.DaysDetails, `
		SELECT
			DATE(due_at) AS date,
			0 AS score_by_day
		FROM exercices
		JOIN users ON exercices.user_id = users.id
		JOIN levels ON levels.id = exercices.level_id
		WHERE exercices.matter_id = ? AND users.campus_id = ? AND due_at < NOW()
		AND IF(? != "", levels.name LIKE ?, 1)
		GROUP BY date
	`, matter_id, campus_id, params.Level, params.Level)

	if err != nil {
		printErr("get exos date", "getAllStudentScoring", err)
		return studentScoring, err
	}

	// get exos by day
	for i, day := range studentScoring.DaysDetails {
		err = db.Select(&studentScoring.DaysDetails[i].Exos, `
		SELECT
			exercices.name AS exo_name,
			exercices.git_path AS repo,
			subjects.name AS exo_lang,
			IFNULL(Cast(AVG(IF(studies.id IS NULL,0, rendus.score)) AS Int),0) AS exo_score,
			exercices.bareme AS exo_total
		FROM exercices
		JOIN subjects ON exercices.subject_id = subjects.id
		JOIN users AS u1 ON exercices.user_id = u1.id
		JOIN levels ON levels.id = exercices.level_id
		LEFT JOIN rendus ON exercices.id = rendus.exercice_id
		LEFT JOIN users AS u2 ON u2.id = rendus.student_id
		LEFT JOIN studies ON studies.id = u2.studies_id AND IF(? != "", studies.name LIKE ?, 1)
		WHERE exercices.matter_id = ? AND u1.campus_id = ?
			AND DATE(exercices.due_at) = ?
			AND IF(? != "", levels.name LIKE ?, 1)
		GROUP BY exercices.id
	`, params.Studies, params.Studies, matter_id, campus_id, day.Date, params.Level, params.Level)

		if err != nil {
			printErr("get exos details", "getAllStudentScoring", err)
			return studentScoring, err
		}
	}

	return studentScoring, err
}

func deleteAllUserData(id int64) (err error) {
	db := dbConnect()

	// delete pic too
	pic := ""
	err = db.QueryRow("SELECT pic FROM users WHERE id = ?", id).Scan(&pic)
	if err != nil {
		printErr("get pic", "deleteAllUserData", err)
		return err
	}

	if pic != viper.GetString("default.default_pic") {
		deleteOldPic(pic)
	}

	// Delete things from less to more important
	_, err = db.Exec("DELETE FROM botToken WHERE user_id = ?", id)
	if err != nil {
		printErr("erase bot token", "deleteAllUserData", err)
		return err
	}

	_, err = db.Exec("DELETE FROM rendus WHERE student_id = ?", id)
	if err != nil {
		printErr("erase rendus table", "deleteAllUserData", err)
		return err
	}

	_, err = db.Exec("DELETE FROM exercices WHERE user_id = ?", id)
	if err != nil {
		printErr("erase exercices table", "deleteAllUserData", err)
		return err
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		printErr("erase users table", "deleteAllUserData", err)
		return err
	}

	return err
}

func getUserData(id int64, data Export) (Export, error) {

	db := dbConnect()

	// get created/modified

	err := db.QueryRow(`
		SELECT created, modified
		FROM users
		WHERE id = ?
	`, id).Scan(&data.Infos.Created, &data.Infos.Modified)

	if err != nil {
		printErr("get created and modified for profile", "getUserData", err)
		return data, err
	}

	if data.Infos.Type == "student" {

		// now get for student
		err = db.Select(&data.Grades, `
			SELECT 
				rendus.id, 
				rendus.exercice_id,
				rendus.student_id,
				CONCAT(users.name, " ", users.surname) AS student_name,
				rendus.score,
				rendus.created
			FROM rendus
				JOIN users ON rendus.student_id = users.id
			WHERE rendus.student_id = ?
		`, id)

		if err != nil {
			printErr("get rendus data", "getUserData", err)
		}

		for i, _ := range data.Grades {
			err = db.Get(&data.Grades[i].ExoDetails, `
				SELECT 
					exercices.id,
					exercices.name,
					CONCAT(users.name, " ", users.surname) AS creator,
					exercices.git_path,
					exercices.due_at,
					exercices.description,
					exercices.level_id,
					levels.name AS level_name,
					exercices.matter_id, 
					matters.name AS matter_name,
					exercices.subject_id,
					subject.name AS subject_name, 
					exercices.bareme,
					exercices.created, 
					exercices.modified
				FROM exercices
					JOIN matters ON exercices.matter_id = matters.id
					JOIN subjects ON exercices.ssubjectd = subject.id
					JOIN levels ON exercices.level_id = levels.id
					JOIN users ON exercices.user_id = users.id
				WHERE exercices.id = ?
				`, data.Grades[i].ExerciceID)

			if err != nil {
				printErr("get exercices from rendus", "getUserData", err)
			}
		}
	} else if data.Infos.Type == "prof" {

		// get all exercises
		err = db.Select(&data.Exos, `
			SELECT 
				exercices.id,
				exercices.name,
				CONCAT(users.name, " ", users.surname) AS creator,
				exercices.git_path,
				exercices.due_at,
				exercices.description,
				exercices.level_id,
				levels.name AS level_name,
				exercices.matter_id, 
				matters.name AS matter_name,
				exercices.subject_id,
				subjects.name AS subject_name, 
				exercices.bareme,
				exercices.created, 
				exercices.modified
			FROM exercices
				JOIN matters ON exercices.matter_id = matters.id
				JOIN subjects ON exercices.subject_id = subjects.id
				JOIN levels ON exercices.level_id = levels.id
				JOIN users ON exercices.user_id = users.id
			WHERE exercices.user_id = ?
		`, id)

		if err != nil {
			printErr("get exercices from prof", "getUserData", err)
		}

		// get corrections by exercice
		for i, _ := range data.Exos {
			err = db.Select(&data.Exos[i].Rendus, `
				SELECT 
					rendus.id,
					rendus.student_id,
					CONCAT(users.name, " ", users.surname) AS student_name,
					rendus.score,
					rendus.created
				FROM rendus 
					JOIN users ON rendus.student_id = users.id
				WHERE rendus.exercice_id = ?
			`, data.Exos[i].Id)

			if err != nil {
				printErr("get rendus from exercices", "getUserData", err)
			}
		}
	}

	return data, err
}
