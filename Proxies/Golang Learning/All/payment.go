package main

// import (
// 	"errors"
// 	"fmt"
// )

// type PaymentGateway interface {
// 	Pay(amount float64) error
// }

// type CreditCard struct {
// 	CardNumber string
// }

// type PayPal struct {
// 	Email string
// }
// type phonepe struct {
// 	UpiID string
// }

// func (up phonepe) Pay(amount float64) error {
// 	if amount <= 0 {
// 		return errors.New("amont must be greater than zero")
// 	}
// 	fmt.Printf("Paying %.2f using Phonpe with UPI ID : %s\n", amount, up.UpiID)
// 	return nil
// }
// func (cc CreditCard) Pay(amount float64) error {
// 	if amount <= 0 {
// 		return errors.New("amount must be greater than zero")
// 	}
// 	fmt.Printf("Paying %.2f using Credit Card %s\n", amount, cc.CardNumber)
// 	return nil
// }

// func (pp PayPal) Pay(amount float64) error {
// 	if amount <= 0 {
// 		return errors.New("amount must be greater than zero")
// 	}
// 	fmt.Printf("Paying %.2f using Paypal %s\n", amount, pp.Email)
// 	return nil
// }

// func ProcessPayment(gateway PaymentGateway, amount float64) {
// 	err := gateway.Pay(amount) //Creditcard.pay method or paypal.pay method is executed
// 	if err != nil {
// 		fmt.Println("Payment failed:", err)
// 	} else {
// 		fmt.Println("Payment completed successfully")
// 	}
// }

// func main() {
// 	cc := CreditCard{CardNumber: "1234567890"}
// 	pp := PayPal{Email: "channu@gmail.com"}
// 	upi := phonepe{UpiID: "1234567890@ybl"}
// 	ProcessPayment(cc, 150.0)
// 	ProcessPayment(pp, 200.0)
// 	ProcessPayment(upi, 250.0)
// }
