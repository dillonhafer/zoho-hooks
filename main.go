package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

type Template struct {
	From    string
	To      string
	Subject string
	Body    string
}

func sendEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Command Must be a POST\n")
			return
		}

		emailUser := &EmailUser{
			os.Getenv("GO_GMAIL_USERNAME"),
			os.Getenv("GO_GMAIL_PASSWORD"),
			"smtp.gmail.com",
			587,
		}

		auth := smtp.PlainAuth("",
			emailUser.Username,
			emailUser.Password,
			emailUser.EmailServer,
		)

		r.ParseForm()
		reported_by := fmt.Sprintf("Reported By: %s\n", r.FormValue("reported_by"))
		case_reason := fmt.Sprintf("Reason: %s\n", r.FormValue("case_reason"))
		feedback := fmt.Sprintf("Feedback: %s\n", r.FormValue("feedback"))
		case_number := fmt.Sprintf("Case #%s\n", r.FormValue("case_number"))

		var err error
		var doc bytes.Buffer

		const emailTemplate = `From: &#123;&#123;.From&#125;&#125;
		To: &#123;&#123;.To&#125;&#125;
		Subject: &#123;&#123;.Subject&#125;&#125;

		&#123;&#123;.Body&#125;&#125;

		Sincerely,

		&#123;&#123;.From&#125;&#125;
		`
		context := &Template{
			"dh@dillonhafer.com",
			"dh@dillonhafer.com",
			"New ZoHo Dwelling Case",
			reported_by + case_reason + feedback + case_number,
		}

		t := template.New("emailTemplate")
		t, err = t.Parse(emailTemplate)
		if err != nil {
			log.Print("error trying to parse mail template")
		}

		err = t.Execute(&doc, context)
		if err != nil {
			log.Print("error trying to execute mail template")
		}

		err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
			auth,
			emailUser.Username,
			[]string{"dh@dillonhafer.com"},
			doc.Bytes(),
		)

		if err != nil {
			log.Print("ERROR: attempting to send a mail ", err)
		}
		//err := mc.Run(cmd)
		//fmt.Prinjktln(r.FormValue("email"))
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	fmt.Fprintf(w, "Error: %v", err)
		//	return
		//}
	}
}

func main() {
	http.HandleFunc("/zohocase", sendEmail())
	fmt.Println("Server running and listening on port 3900")
	fmt.Println("Ctrl-C to shutdown server")
	err := http.ListenAndServe(":3900", nil)
	fmt.Fprintln(os.Stderr, err)
}