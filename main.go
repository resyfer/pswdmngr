package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/dlclark/regexp2" //Internal Go Regexp doesn't support backtracking in exchange for constant times

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
)

// Global Constants
const VERSION string = "1.0.0"
const FILE_NAME string = "password.txt"

// Global Variables	
var secret string = ""
var ValidPassword *regexp2.Regexp
var data map[string]string = map[string]string{}

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func save() {

	var saveString string = ""

	encrSecret, _ := Encrypt(secret, secret)
	saveString += encrSecret

	for elem, val := range data {

		if elem == "secret" {
			continue
		}

		encrElem, _ := Encrypt(elem, secret)
		encrVal, _ := Encrypt(val, secret)
		saveString += "\n" + encrElem + " " + encrVal

	}
	
	os.WriteFile(FILE_NAME, []byte(saveString), 0644)

}

func retrieve(fileBytes [][]byte) {

	for i:=1; i < len(fileBytes); i++ {
	
		site := bytes.Split(fileBytes[i], []byte(" "))
		decrElem, _ := Decrypt(string(site[0]), secret)
		decrVal, _ := Decrypt(string(site[1]), secret)

		data[decrElem] = decrVal
	}

}

func options() {

	choice := ""

	err := survey.AskOne(&survey.Select{
		Message: "Please select an option",
		Options: []string{
			"Quit",
			"Get Password",
			"Suggest Password",
			"Add Password",
			"Change Name",
			"Change Password",
			"Delete Password",
		},
	}, &choice)

	if err != nil {
		if err == terminal.InterruptErr {
			log.Fatal("interrupted")
		}
	}

	switch choice {
	case "Get Password":
		name := ""
		
		err := survey.AskOne(&survey.Select{
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

		if err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}

		if name == "QUIT" {
			break
		}
		fmt.Println(data[name])

	case "Suggest Password":
		suggest()

	case "Add Password":
		name := ""
		newPassword := ""
		confirmPassword := " " //Different from newPassword
	
		err1 := survey.AskOne(&survey.Input{
			Message: "Name for New Site (Atleast 4 characters long, Q to quit)",
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

		if err1 != nil {
			if err1 == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}

		if name == "Q" {
			break
		}

		for newPassword != confirmPassword {

			if newPassword != "" {
				color.Red("Passwords Don't Match. Please Try Again.")
			}
			
			// New Password
			err2 := survey.AskOne(&survey.Password{
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
	
			if err2 != nil {
				if err2 == terminal.InterruptErr {
					log.Fatal("interrupted")
				}
			}
	
			if newPassword == "Q" {
				break
			}
	
			// Confirm Password
			err3 := survey.AskOne(&survey.Password{
				Message: "Confirm Password for site (Enter Q to quit)",
			}, &confirmPassword)
	
			if err3 != nil {
				if err3 == terminal.InterruptErr {
					log.Fatal("interrupted")
				}
			}
	
			if confirmPassword == "Q" {
				break
			}
		}
	
		if newPassword == "Q" || confirmPassword == "Q" {
			break
		}

		add(name, newPassword)

	case "Change Name":

		old := ""
		new := ""
		confirm := false

		// Old Name question
		err1 := survey.AskOne(&survey.Select{
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
		
		if err1 != nil {
			if err1 == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}

		if old == "QUIT" {
			break
		}

		// New Name question
		err2 := survey.AskOne(&survey.Input{
			Message: "New Name for Site (Atleast 4 characters long, Q for quit)",
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
			
		if err2 != nil {
			if err2 == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}

		if new == "Q" {
			break
		}

		survey.AskOne(&survey.Confirm{
			Message: "Are you sure you want to delete it?",
		}, &confirm)

		if !confirm {
			break
		}
		
		changeName(old, new)

	case "Change Password":
		siteName := ""
		newPassword := ""
		confirmPassword := " " //Different from newPassword

		// Old Name question
		err1 := survey.AskOne(&survey.Select{
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
		
		if err1 != nil {
			if err1 == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}

		if siteName == "QUIT" {
			break
		}

		for newPassword != confirmPassword {

			//New Password
			err2 := survey.AskOne(&survey.Password{
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
				
			if err2 != nil {
				if err2 == terminal.InterruptErr {
					log.Fatal("interrupted")
				}
			}
	
			if newPassword == "Q" {
				break
			}
			
			//Confirm Password
			err3 := survey.AskOne(&survey.Password{
				Message: "Enter Password for site (Enter Q to quit)",
			}, &confirmPassword, survey.WithValidator(
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
				
			if err3 != nil {
				if err3 == terminal.InterruptErr {
					log.Fatal("interrupted")
				}
			}
	
			if confirmPassword == "Q" {
				break
			}
		}
	
		if newPassword == "Q" || confirmPassword == "Q" {
			break
		}

		changePassword(siteName, newPassword)

	case "Delete Password":
		site := ""
		confirm := false

		// Old Name question
		err := survey.AskOne(&survey.Select{
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

		}, &site)
		
		if err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}

		if site == "QUIT" {
			break
		}

		survey.AskOne(&survey.Confirm{
			Message: "Are you sure you want to delete it?",
		}, &confirm)

		if !confirm {
			break
		}

		delete(data, site)

	case "Quit":
		save()
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

	for length < 12 {
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
	var fileEmpty bool = true
	_, errStat := os.Stat("password.txt")
	if errStat != nil {
		fileEmpty = false
		file, errCreate := os.Create("password.txt")

		if errCreate != nil {
			log.Fatal("ERROR: Could not create Password File")
		}
		file.Close()
	}

	// Secret Code if it exists
	fileByteContent, _ := os.ReadFile(FILE_NAME)
	if len(fileByteContent) == 0 {
		fileEmpty = false
	}

	fileByteArray := bytes.Split(fileByteContent, []byte("\n"))
	
	// Getting Secret Code and if present, verifying it
	var attemptLimit int = 5
	var attempts int = 0

	if fileEmpty {


		err := survey.AskOne(&survey.Password{
			Message: "Please enter your secret code ",
			Help: "Please enter your secret code that only you know to access the passwords and save/edit them",
		}, &secret, survey.WithValidator(
			func () survey.Validator {
				return func (val interface{}) error {
					
					clientSecret, _ := val.(string)
					decoded, _ := Decrypt(string(fileByteArray[0]), clientSecret)

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

			if err != nil {
				if err == terminal.InterruptErr {
					log.Fatal("interrupted")
				}
			}

			retrieve(fileByteArray)


	} else {

		confirmSecret := " " //Different from Secret

		for secret != confirmSecret {

			if secret != "" {
				color.Red("Passwords don't match. Please try again.")
			}

			err1 := survey.AskOne(&survey.Password{
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
	
			if err1 != nil {
				if err1 == terminal.InterruptErr {
					log.Fatal("interrupted")
				}
			}

			err2 := survey.AskOne(&survey.Password{
				Message: "Please confirm your secret code",
				Help: "Please confirm your secret code that only you know to access the passwords and save/edit them",
			}, &confirmSecret)
	
			if err2 != nil {
				if err2 == terminal.InterruptErr {
					log.Fatal("interrupted")
				}
			}
		}
		
		color.Green("\nPlease remember this secret code ⊂(￣▽￣)⊃\n")
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