package utility

import (
	"log"
	"os"
	"fmt"
	"time"
	"os/exec"
	"syscall"
	"io"
)


// creates logfiles or writes to them. if the current logfile crawler.log becomes too big, logfiles are split.
// because logging is absolutey crucial this function is allowed to panic on errors. we dont want a zombie crawler that doesnt protocoll
// its actions.
func LogToFile(what string) {

	//check is ./logs/ subdirectory exists and make it if not
	if _, err := os.Stat(CRAWLER_LOG_FILE_PATH); os.IsNotExist(err) {
		cmd := exec.Command("mkdir", "-p", CRAWLER_LOG_FILE_PATH)
		_ = syscall.Umask(0077) // Set umask for this process
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		//cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
		cmd.Run()
	}

	// check if crawler.log exists. if it exists, check its size.
	// if size if too big, copy it into logs with name crawlerYYYY-MM-DD.log and make a new one.
	// if it doesnt exist, make it an write to it
	// ATTENTION: the unlikely case in which there are 2 logfile splits in the same second is currently unhandled and causes problems
	if _, err := os.Stat(CRAWLER_LOG_FILE_PATH + CRAWLER_LOG_FILE_NAME + CRAWLER_LOG_FILE_MIMETYPE); os.IsNotExist(err) {
		//crawler.log doesnt exist. we make a new one and write to it.
		f, err := os.OpenFile(CRAWLER_LOG_FILE_PATH + CRAWLER_LOG_FILE_NAME + CRAWLER_LOG_FILE_MIMETYPE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}

		// don't forget to close it
		defer f.Close()

		// assign it to the standard logger
		log.SetOutput(f)
		log.Output(1, what)
	} else {
		//crawler log does exists
		f,_ := os.Stat(CRAWLER_LOG_FILE_PATH + CRAWLER_LOG_FILE_NAME + CRAWLER_LOG_FILE_MIMETYPE)
		if f.Size() > 500000/*00*/ { // 500 KB
			//crawler log does exist but is too big
			// first we copy the file crawler.log and rename the copy with current timestamp
			from, err := os.Open(CRAWLER_LOG_FILE_PATH + CRAWLER_LOG_FILE_NAME + CRAWLER_LOG_FILE_MIMETYPE)
			if err != nil {
				log.Fatal(err)
			}
			defer from.Close()

			current_time := time.Now().Local()
			dt := fmt.Sprintf(current_time.Format("2006-01-02 15:04:05"))

			to, err := os.OpenFile(CRAWLER_LOG_FILE_PATH + CRAWLER_LOG_FILE_NAME + dt + CRAWLER_LOG_FILE_MIMETYPE, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer to.Close()

			_, err = io.Copy(to, from)
			if err != nil {
				log.Fatal(err)
			}

			//then we remove crawler.log
			cmd := exec.Command("rm", "-r", CRAWLER_LOG_FILE_PATH + CRAWLER_LOG_FILE_NAME + CRAWLER_LOG_FILE_MIMETYPE)
			_ = syscall.Umask(0077) // Set umask for this process
			cmd.SysProcAttr = &syscall.SysProcAttr{}
			//cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
			cmd.Run()

			//now we create a new crawler log and write to it.
			f, err := os.OpenFile(CRAWLER_LOG_FILE_PATH + CRAWLER_LOG_FILE_NAME + CRAWLER_LOG_FILE_MIMETYPE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				panic(err)
			}

			// don't forget to close it
			defer f.Close()

			// assign it to the standard logger
			log.SetOutput(f)
			log.Output(1, what)

		} else {
			//crawler log does exist but is fine. just open it and write to it
			f, err := os.OpenFile(CRAWLER_LOG_FILE_PATH + CRAWLER_LOG_FILE_NAME + CRAWLER_LOG_FILE_MIMETYPE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				panic(err)
			}

			// don't forget to close it
			defer f.Close()

			// assign it to the standard logger
			log.SetOutput(f)
			log.Output(1, what)
		}
	}
}

func CheckErr(err error) {
	if err != nil {
		LogToFile(err.Error())
	}
}
