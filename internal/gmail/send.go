package gmail

import (
	"fmt"
	"log"
	"net/smtp"
	"upbit-api/config"
)

// Send .
// subject : 메일 제목
// body : 내용
func Send(subject, body string) {
	// 발신자 이메일 정보
	from := "myungsworld@gmail.com"
	password := config.GmailAppPassword

	// 수신자 이메일 주소
	to := "myungsworld@gmail.com"

	// SMTP 서버 설정
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"

	// 이메일 보내기
	err := sendEmail(smtpServer, smtpPort, from, password, to, subject, body)
	if err != nil {
		log.Fatal(err)
	}
}

func sendEmail(smtpServer, smtpPort, from, password, to, subject, body string) error {
	// 이메일의 헤더와 본문 생성
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// SMTP 인증 설정
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// SMTP 연결 설정
	serverAddr := fmt.Sprintf("%s:%s", smtpServer, smtpPort)

	// 메일 전송
	err := smtp.SendMail(serverAddr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("이메일 전송 실패: %v", err)
	}
	log.Print(subject)
	log.Print(body)
	return nil
}
