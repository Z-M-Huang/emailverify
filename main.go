package main

import (
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//VerifyEmail verify if email is a valid email
func VerifyEmail(emailAddress string) bool {
	//expression checking
	if !emailRe.Match([]byte(emailAddress)) {
		return false
	}

	_, host := splitEmail(emailAddress)
	//check MX
	mx, err := net.LookupMX(host)
	if err != nil {
		return false
	}

	//dail check
	client, err := dialCheck(fmt.Sprintf("%s:%d", mx[0].Host, 25))
	if err != nil {
		return false
	}
	defer client.Close()

	//HELO check
	err = client.Hello("emailcheck.zhcode.com")
	if err != nil {
		return false
	}

	//fake email check
	err = client.Mail("noreply@zh-code.com")
	if err != nil {
		return false
	}

	//check rcpt
	err = client.Rcpt(emailAddress)
	if err != nil {
		return false
	}
	return true
}

func splitEmail(emailAddress string) (account, host string) {
	index := strings.LastIndexByte(emailAddress, '@')
	account = emailAddress[:index]
	host = emailAddress[index+1:]
	return
}

func dialCheck(address string) (*smtp.Client, error) {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}

	t := time.AfterFunc(timeout, func() { conn.Close() })
	defer t.Stop()

	host, _, _ := net.SplitHostPort(address)
	return smtp.NewClient(conn, host)
}
