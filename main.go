package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dlclark/regexp2" //Internal Go Regexp doesn't support backtracking in exchange for constant times

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// Global Constants
const VERSION string = "1.0.0"
const FILE_NAME string = "password.txt"

// Global Variables
var ValidPassword *regexp2.Regexp
var data map[string]string = map[string]string{}

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func options() {

	choice := ""
	survey.AskOne(&survey.Select{
		Message: "Please select an option (Enter Q to quit)",
		Options: []string{
			"Get Password",
			"Suggest Password",
			"Add Password",
			"Change Name",
			"Change Password",
			"Quit",
		},
	}, &choice)

	switch choice {
	case "Get Password":
		name := ""
		
		survey.AskOne(&survey.Select{
			Message: "Name of Site",
			Options: func () []string {
				names := []string{}
				names = append(names, "QUIT")

				for name := range data {
					if name != "secret" {
						names = append(names, name)
					}
				}

				return names
			}(),

		}, &name)

		if name == "QUIT" {
			break
		}
		fmt.Println(data[name])

	case "Suggest Password":
		suggest()

	case "Add Password":
		name := ""
		newPassword := ""
	
		survey.AskOne(&survey.Input{
			Message: "Name for New Site (Atleast 4 characters long)",
		}, &name, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					newName := val.(string)

					
					if newName == "Q" {
						return nil
					}
					
					if len(newName) < 4 {
						return fmt.Errorf("name should be atleast 4 characters long ლ(ಠ_ಠ ლ)")
					}

					_, ok := data[newName]

					if !ok {
						return nil
					} else {
						return fmt.Errorf("name already exists, please select a new one ლ(ಠ_ಠ ლ)")
					}
				}
			}()))

		if name == "Q" {
			break
		}

		survey.AskOne(&survey.Password{
			Message: "Enter Password for site (Enter Q to quit)",
		}, &newPassword, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					pswd := val.(string)

					if pswd == "Q" {
						return nil
					}

					if len(pswd) < 8 {
						return fmt.Errorf("a good password needs to be aleast 8 characters long |･д･)ﾉ")
					}

					// Didn't add password strength checker to allow password freedom
					return nil
				}
			}()))

		if newPassword == "Q" {
			break
		}

		add(name, newPassword)

	case "Change Name":

		old := ""
		new := ""

		// Old Name question
		survey.AskOne(&survey.Select{
			Message: "Name of Site",
			Options: func () []string {
				names := []string{}
				names = append(names, "QUIT")

				for name := range data {
					if name != "secret" {
						names = append(names, name)
					}
				}

				return names
			}(),

		}, &old)

		if old == "QUIT" {
			break
		}

		// New Name question
		survey.AskOne(&survey.Input{
			Message: "New Name for Site (Atleast 4 characters long)",
		}, &new, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					newName := val.(string)

					
					if newName == "Q" {
						return nil
					}
					
					if len(newName) < 4 {
						return fmt.Errorf("name should be atleast 4 characters long ლ(ಠ_ಠ ლ)")
					}

					_, ok := data[newName]

					if !ok {
						return nil
					} else {
						return fmt.Errorf("name already exists, please select a new one ლ(ಠ_ಠ ლ)")
					}
				}
			}()))

		if new == "Q" {
			break
		}
		
		changeName(old, new)

	case "Change Password":
		newPassword := ""
		siteName := ""

		// Old Name question
		survey.AskOne(&survey.Select{
			Message: "Name of Site",
			Options: func () []string {
				names := []string{}
				names = append(names, "QUIT")

				for name := range data {
					if name != "secret" {
						names = append(names, name)
					}
				}

				return names
			}(),

		}, &siteName)

		if siteName == "QUIT" {
			break
		}

		//New Password
		survey.AskOne(&survey.Password{
			Message: "Enter Password for site (Enter Q to quit)",
		}, &newPassword, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					pswd := val.(string)

					if pswd == "Q" {
						return nil
					}

					if len(pswd) < 8 {
						return fmt.Errorf("a good password needs to be aleast 8 characters long |･д･)ﾉ")
					}

					// Didn't add password strength checker to allow password freedom
					return nil
				}
			}()))

		if newPassword == "Q" {
			break
		}

		changePassword(siteName, newPassword)

	case "Quit":
		quit()
	}
	fmt.Printf("\n\n")

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

	color.Green(password)
}

func add(name, password string) {

	data[name] = password

}

func changeName(oldName, newName string) {
	data[newName] = data[oldName]
	delete(data, oldName)
}

func changePassword(name, newPassword string) {
	data[name] = newPassword
}

func quit() {
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
	fmt.Printf("\n")
	color.Magenta("Disclaimer : The file %v is where your passwords are stored. It is secure enough not to be decryptable without the secret code, even by the developers.\n\n", FILE_NAME)
	color.Magenta("Every data is stored only on your PC in your %v file. Don't forget your secret code, and don't lose the file (placing another file with the same name can work as import/export of passwords) else it is irrecoverable. For recovering guarantee, always keep a copy of %v in Cloud Storage.", FILE_NAME, FILE_NAME)
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

	//TODO: Stop program on keyboard interrupts
	if fileExists {
		survey.AskOne(&survey.Password{
			Message: "Please enter your secret code ",
			Help: "Please enter your secret code that only you know to access the passwords and save/edit them",
		}, &secret, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					
					clientSecret, _ := val.(string)
					decoded, _ := Decrypt(fileStringContent[0], clientSecret)

					if attempts <= attemptLimit && decoded == clientSecret {
						data["secret"] = decoded
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
			Message: "Please enter a secret code",
			Help: "Please enter a secret code that only you know to access the passwords and save/edit them",
		}, &secret, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					
					clientSecret, _ := val.(string)

					if len(clientSecret) < 8 && len(clientSecret) > 32 {
						return fmt.Errorf("length of secret code needs to be between 8-32 (っ˘̩╭╮˘̩)っ")
					}

					data["secret"] = clientSecret
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

	for {
		options()
	}

}