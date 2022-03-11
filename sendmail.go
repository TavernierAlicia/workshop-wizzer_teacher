package main

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendResetMail(mail string, link string) (err error) {

	//configure sending mailbox
	from := viper.GetString("sendmail.service_mail")
	pass := viper.GetString("sendmail.service_pwd")

	subject := "Wizzer Teacher - Réinitialisation de mot de passe"
	//set message
	message := ` 
	<p> Bonjour! Vous avez demandé une réinitialisation de mot de passe \n Pour réinitialiser votre mot de passe, veuillez cliquer sur ce lien: <a href="` + link + `"> ` + link + ` </a> </p> 
	<p> Si vous n'avez pas demandé de réinitialisation de mot de passe, nous vous conseillons de vous reconnecter avec vos identifiants dans les plus brefs délais afin d'écarter toute menace. </p>
	`
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", mail)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewPlainDialer("smtp.gmail.com", 587, from, pass)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return err

}
