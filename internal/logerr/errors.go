package logerr

import (
	"io"
	"os"
)

// -------------------------------------------------------------------------------------------
func (l *Logerr) Exist() bool {
	_, err := os.Stat(l.FilePath)
	return !os.IsNotExist(err)
}

// -------------------------------------------------------------------------------------------
func (l *Logerr) ErrLog(s string) {
	var f *os.File
	var err error

	if l.Exist() {
		f, err = os.OpenFile(l.FilePath, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			println("[ ERR ]", "open log file :", err.Error())
			return
		}
	} else {
		f, err = os.OpenFile(l.FilePath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			println("[ ERR ]", "create log file :", err.Error())
			return
		}
	}
	_, err = io.WriteString(f, s+"\n")
	if err != nil {
		println("[ ERR ]", "write to log file :", err.Error())
		return
	}
}
