package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var VERSION string = "1.0.0"

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func main() {

	// Intro
	clear()
	fmt.Printf("Welcome to Resyfer's Password Manager v%v\n", VERSION)

	// Opening file or Creating it
	var fileExists bool = true
	_, errStat := os.Stat("password.txt")
	if errStat != nil {
		fileExists = false
		file, errCreate := os.Create("password.txt")

		if errCreate != nil {
			log.Fatal("ERROR: Could not create Password File")
		}
		file.Close()
	}

	// Secret Code if it exists
	fileByteContent, _ := os.ReadFile("password.txt")
	fileStringContent := strings.Split(string(fileByteContent), "\n")

	// Secret
	var secret string = ""
	var attemptLimit int = 5
	var attempts int = 0

	// Getting Secret Code
	if fileExists {
		survey.AskOne(&survey.Password{
			Message: "Please enter your secret code > ",
			Help: "Please enter your secret code that only you know to access the passwords and save/edit them",
		}, &secret, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					
					clientSecret, _ := val.(string)
					decoded, _ := Decrypt(fileStringContent[0], clientSecret)

					if attempts <= attemptLimit && decoded == clientSecret {
						return nil
					} else if attempts > attemptLimit {
						return nil
					} else {
						attempts++
						return fmt.Errorf("wrong secret code...bad luck (╥_╥)")
					}

				}
			}()))
	} else {
		survey.AskOne(&survey.Password{
			Message: "Please enter a secret code > ",
			Help: "Please enter a secret code that only you know to access the passwords and save/edit them",
		}, &secret, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					
					clientSecret, _ := val.(string)

					if len(clientSecret) < 8 && len(clientSecret) > 32 {
						return fmt.Errorf("length of secret code needs to be between 8-32 (っ˘̩╭╮˘̩)っ")
					}
					return nil
				}
			}()))
		clear()
		fmt.Println("Please remember this secret code ⊂(￣▽￣)⊃")
	}

	if attempts > attemptLimit {
		fmt.Println("Too many tries mate, we know you ain't the Chosen One ヽ( `д´*)ノ")
		os.Exit(0)
	}


	// //Clear Login
	// loginClear := exec.Command("clear")
	// loginClear.Stdout = os.Stdout
	// loginClear.Run()

}