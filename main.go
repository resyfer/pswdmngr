package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dlclark/regexp2" //Internal Go Regex doesn't support backtracking in exchange for constant times

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

var VERSION string = "1.0.0"
var FILE_NAME string = "password.txt"

var ValidPassword *regexp2.Regexp

var data [][]string = [][]string{}

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func options() {

	choice := ""
	survey.AskOne(&survey.Select{
		Message: "Please select an option",
		Options: []string{
			"Get Password",
			"Suggest Password",
			"Save Password",
			"Quit",
		},
	}, &choice)

	switch choice {
	case "Suggest Password":
		suggest()
	case "Quit":
		quit()
	}

}

func passwordStrengthChecker(password string) bool {
	match, _ := ValidPassword.MatchString(password)
	return match
}

func randomPassword() string {

	length := 0

	for length < 8 {
		length = rand.Intn(32)
	}

	password := ""
	for i:=0; i<length; i++ {

		ascii := 0
		for ascii < 33 {
			ascii = rand.Intn(126)
		}

		password += string(rune(ascii))
	}

	return password
}

func suggest() {
	password := ""

	for !passwordStrengthChecker(password) {
		password = randomPassword()
	}

	fmt.Println(password)
}

func quit() {
	fmt.Println("Hello")
	os.Exit(0)
}

func main() {
	// Compiling RegEx for Password Checker
	ValidPassword = regexp2.MustCompile(`^(?=.*[0-9])(?=.*[!"#$%&'()*+,-./:;<=>?@[\]^_{|}~])[a-zA-Z0-9!"#$%&'()*+,-./:;<=>?@[\]^_{|}~]{8,100}$`, 0)

	rand.Seed(time.Now().UnixNano())

	// Intro
	clear()
	color.Yellow("Welcome to Resyfer's Password Manager v%v\n", VERSION)
	color.Red("NOTE: Only use the Quit option to exit the program, or changes won't be saved")
	fmt.Printf("\n\n")

	// Checking file or Creating it
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
	fileByteContent, _ := os.ReadFile(FILE_NAME)
	fileStringContent := strings.Split(string(fileByteContent), "\n")
	
	// Secret
	var secret string = ""
	
	// Getting Secret Code and if present, verifying it
	var attemptLimit int = 5
	var attempts int = 0
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
						data = append(data, strings.Split(decoded, ""))
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

					data = append(data, strings.Split(clientSecret, ""))
					return nil
				}
			}()))
		fmt.Println("Please remember this secret code ⊂(￣▽￣)⊃")
	}
	// Stop for too many attempts
	if attempts > attemptLimit {
		fmt.Println("Too many tries mate, we know you ain't the Chosen One ヽ( `д´*)ノ")
		os.Exit(0)
	}

	options()

}