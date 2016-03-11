// The MIT License (MIT)
//
// Copyright (c) 2016 aerth
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

package main

import (
	//"strconv"
	"bytes"
	"fmt"
	"github.com/astaxie/beego/session"
	"github.com/logpacker/PayPal-Go-SDK"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type User struct {
	ID    int
	Name  []byte
	Email []byte
	Date  []byte
}

var PayPalC = os.Getenv("PayPalC")
var PayPalK = os.Getenv("PayPalK")

func main() {
	if PayPalC == "" || PayPalK == "" {

		log.Fatalln("Set PayPalC and PayPalK environmental variables.")
	}

	http.HandleFunc("/join", joinhandler)
	http.HandleFunc("/", homehandler)
	http.HandleFunc("/success", successhandler)
	http.HandleFunc("/confirm", confirmhandler)
	http.HandleFunc("/login", loginhandler)
	http.HandleFunc("/cancel", cancelhandler)
	http.HandleFunc("/fail", failhandler)
	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
func joinhandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	logdb := log.New(&buf, "logdb: ", 1)
	f, err := os.OpenFile("./sub.db", os.O_RDWR|os.O_CREATE|
		os.O_APPEND, 0600)
	if err != nil {
		log.Fatal("error opening file: %v", err)
		os.Exit(1)
	}
	logdb.SetOutput(f)

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		log.Println("Session Error:")
		log.Println(err)
	}
	defer sess.SessionRelease(w)
	log.Println(sess)
	//email := sess.Get("email")
	//if email != "" {fmt.Println(email)}
	if r.Method == "GET" {

		t, err := template.New("Index").ParseFiles("login.gtpl")
		if err != nil {
			log.Println("Template Error:")
			log.Println(err)
		}
		log.Println("Template Parsed. Login form presented.")
		t.ExecuteTemplate(w, "Index", nil)
	} else {
		r.ParseForm()
		logdb.Println("email:", r.Form["email"])
		sess.Set("email", r.Form["email"])
		logdb.Println("password:", r.Form["password"])

		c, err := paypalsdk.NewClient(PayPalC, PayPalK, paypalsdk.APIBaseSandBox)
		if err == nil {
			log.Println("ClientID: " + c.ClientID)
			log.Println("APIBase: " + c.APIBase)
		} else {
			log.Println("ERROR: " + err.Error())
		}

		log.Println("Getting new AccessToken")
		token, err := c.GetAccessToken()
		if err == nil {
			log.Println("AccessToken: " + token.Token)

		} else {
			fmt.Println("ERROR: " + err.Error())
		}

		amount := paypalsdk.Amount{
			Total:    "10.01",
			Currency: "USD",
		}
		redirectURI := "http://127.0.0.1:8080/success"
		cancelURI := "http://127.0.0.1:8080/cancel"
		now := time.Now()
		nowtime := now.Format("Mon, 2 Jan 2006 15:04:05 -0700")
		description := "Membership (one year) " + nowtime
		log.Println(description)
		paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)
		/*
		   payment, err := c.GetPayment(paymentResult.ID)
		   payments, err := c.GetPayments()
		   if err == nil {
		     fmt.Println("DEBUG: PaymentsCount=" + strconv.Itoa(len(payments)))
		   } else {
		     fmt.Println("ERROR: " + err.Error())
		   }

		   fmt.Println("Payment ID:")
		   fmt.Println(payment.ID)
		   fmt.Println("Payment:")
		   fmt.Println(payment)
		   fmt.Println("Payments:")
		   fmt.Println(payments)
		*/

		//http.Redirect(w, r, paymentResult.Links[1].Href, 302)
		fmt.Fprintf(w, "<!DOCTYPE html><html><a href=\""+paymentResult.Links[1].Href+"\">Click here to use PayPal ($ 10)</a></html>")

		log.Println("Redirecting " + r.RemoteAddr + paymentResult.ID + paymentResult.Links[1].Href)
	}
}

func successhandler(w http.ResponseWriter, r *http.Request) {
	m, _ := url.ParseQuery(r.URL.RawQuery)
	log.Println("Returned from PayPal.")
	paymentid := m["paymentId"][0]
	payerid := m["PayerID"][0]
	log.Println("Payment ID: " + paymentid)
	log.Println("Payer ID: " + payerid)
	c, err := paypalsdk.NewClient(PayPalC, PayPalK, paypalsdk.APIBaseSandBox)
	if err == nil {
		log.Println("ClientID: " + c.ClientID)
		log.Println("APIBase: " + c.APIBase)
	} else {
		log.Println("ERROR: " + err.Error())
	}

	log.Println("Getting new AccessToken")
	token, err := c.GetAccessToken()
	if err == nil {
		log.Println("AccessToken: " + token.Token)

	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	payment, err := c.GetPayment(paymentid)
	if err != nil {
		log.Println("Paypal Error")
		log.Println(err.Error())
		http.Redirect(w, r, "/fail?id="+paymentid, 302)
	} else {
		log.Println("Pre-Confirm:")
		log.Println(payment.Intent, payment.Payer.PayerInfo.FirstName, payment.Payer.PayerInfo.LastName)
		log.Println(payment.Payer.PayerInfo.Email, payment.Payer.PayerInfo.Phone, payment.Payer.PayerInfo.PayerID)
		transaction := payment.Transactions
		fmt.Fprintf(w, "<!DOCTYPE html><html><a href=\"/confirm?id="+payment.ID+"&u="+payment.Payer.PayerInfo.PayerID+"\">Click here to Confirm Payment of $%s</a></html>", transaction[0].Amount.Total)
	}

}
func cancelhandler(w http.ResponseWriter, r *http.Request) {

	m, _ := url.ParseQuery(r.URL.RawQuery)
	fmt.Println(m)

}
func failhandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Failed.")
	c, err := paypalsdk.NewClient(PayPalC, PayPalK, paypalsdk.APIBaseSandBox)
	if err != nil {
		log.Println(err)
	}
	m, _ := url.ParseQuery(r.URL.RawQuery)
	paymentid := m["id"][0]
	log.Println(paymentid)
	payment, err := c.GetPayment(paymentid)
	if payment != nil {
		log.Println(payment)
	}
	if err != nil {
		log.Println(err.Error())
	}
	if payment != nil {
		log.Println(payment.ID)
	}

	http.Redirect(w, r, "/", 302)

}
func homehandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<!DOCTYPE html><html><a href=\"/join\">Join now -- $10 bucks!</a></html>")
	log.Println("Home Visitor:")
	log.Println(r)
}
func loginhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<!DOCTYPE html><html><a href=\"/login\">Login!</a></html>")
}
func confirmhandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Payment Confirmed by User. Sending to Paypal.")
	c, err := paypalsdk.NewClient(PayPalC, PayPalK, paypalsdk.APIBaseSandBox)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	m, _ := url.ParseQuery(r.URL.RawQuery)
	paymentid := m["id"][0]
	payerid := m["u"][0]
	//payment, err := c.GetPayment(paymentid)
	log.Println("Payment ID: " + paymentid)
	log.Println("Payer ID: " + payerid)
	c.GetAccessToken()
	response, err := (c.ExecuteApprovedPayment(paymentid, payerid))
	log.Println(response)
	if err != nil {
		log.Println("Error:")
		log.Println(err.Error())

	}
	if response.State != "" {
		log.Println("Approval State: " + response.State)

	}
	//c.ExecuteApprovedPayment("nil", "nil")
	//if err != nil {fmt.Println(err)}
	//fmt.Println(executeResponse)

	log.Println("Redirecting to /login")
	http.Redirect(w, r, "/login", 302)
}
func init() {
	c, err := paypalsdk.NewClient("clientID", "secretID", paypalsdk.APIBaseSandBox)
	if err != nil {
		log.Fatalln(err)
	}
	c.SetLogFile("debug.log") // Set log file if necessary
	globalSessions, _ = session.NewManager("file", `{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"./tmp"}`)
	go globalSessions.GC()

	go func() {
		for {
			payments, err := c.GetPayments()
			time.Sleep(60 * time.Second)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("PULSE")
			fmt.Println(payments)
		}
	}()

}

var globalSessions *session.Manager
