package main

import (
	"fmt"
	"github.com/firdasafridi/gocrypt"
	"github.com/joho/godotenv"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	err := godotenv.Load(filepath.Join("C:\\Users\\Simon\\GolandProjects\\progress-team-viewer-v2", ".env"))

	if err != nil {
		fmt.Println(err)
	}

	aeskey := os.Getenv("PRIVATE_DECRYPT_KEY")

	if _, err := os.Stat(fmt.Sprintf("%s%s%s", os.Getenv("EXE_FILE_DIR"), `\`, os.Getenv("EXE_FILE_NAME"))); err != nil {
		var path = os.Getenv("REG_FILE_PATH")

		var file, err = os.Create(path)
		if err != nil {
			log.Println("ERR create path", err)
			return
		}
		_, err = file.WriteString(
			fmt.Sprintf("%s \n %s\n %s \n %s \n %s\n %s \n %s \n",
				"Windows Registry Editor Version 5.00 ",
				`[HKEY_CLASSES_ROOT\webant]`,
				`"URL Protocol"="\"\""`,
				`[HKEY_CLASSES_ROOT\webant\shell]`,
				`[HKEY_CLASSES_ROOT\webant\shell\open]`,
				`[HKEY_CLASSES_ROOT\webant\shell\open\command]`,
				`@="C:\\Program Files (x86)\\webant\\webantTeamViewer %1"`))
		if err != nil {
			log.Println("ERR write reg file", err)
			return
		}

		// save changes
		err = file.Sync()
		if err != nil {
			log.Println("ERR sync file", err)
			return
		}
		err = file.Close()
		if err != nil {
			log.Println("ERR close file", err)
			return
		}

		defer startReg()
		err = os.MkdirAll(os.Getenv("EXE_FILE_DIR"), os.ModePerm)
		if err != nil {
			log.Println("ERR", err)
			return
		}
		cmdCopy := exec.Command("cmd", "/C", fmt.Sprintf("copy %s %s", os.Getenv("EXE_FILE_NAME"), os.Getenv("EXE_FILE_PATH")))
		err = cmdCopy.Run()
		if err != nil {
			log.Println("ERR exec command copy file", err)
			return
		}
		return
	}
	arg := os.Args[1]

	if arg != "" {
		arr := strings.Split(arg, "/")

		desOpt, err := gocrypt.NewDESOpt(aeskey)
		if err != nil {
			log.Println("ERR init decrypt", err)
			return
		}

		decryptPass, err := desOpt.Decrypt([]byte(arr[3]))
		if err != nil {
			log.Println("ERR decrypt pass", err)
			return
		}

		cmd := exec.Command(os.Getenv("TEAMVIEWER_PATH"), "-i", arr[2], "-P", decryptPass)
		err = cmd.Start()

		if err != nil {
			log.Println("ERR exec command start tv", err)
			return
		}
	}

	return
}

func startReg() {
	cmd := exec.Command("cmd", "/C", fmt.Sprintf("start %s"), os.Getenv("REG_FILE_PATH"))
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		log.Println("ERR add reg", err)
		return
	}

	for true {
		k, err := registry.OpenKey(registry.CLASSES_ROOT, `webant`, registry.QUERY_VALUE)
		if err == nil {

			break
		}
		err = k.Close()

		if err != nil {
			log.Println("ERR with reg", err)
			return
		}
	}
}
