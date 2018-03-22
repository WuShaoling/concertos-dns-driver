package executor

import (
	"sync"
	"os"
	"bufio"
	"io"
	"log"
	"strings"
)

type Executor struct {
}

const DNS_CONFIG_PATH = "/etc/hosts"

func (exec *Executor) DeleteRecord(domain string) error {
	// read from file
	var lines []string
	f, err := os.OpenFile(DNS_CONFIG_PATH, os.O_RDONLY, 0666)
	if err != nil {
		log.Println(err)
		f.Close()
		return err
	}
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.Contains(line, domain) {
			continue
		}
		lines = append(lines, line)
	}
	f.Close()

	// delete file
	if err = os.Remove(DNS_CONFIG_PATH); nil != err {
		log.Println(err)
		return err
	}

	// create and write to file
	f, err = os.OpenFile(DNS_CONFIG_PATH, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		f.Close()
		return err
	}
	for _, v := range lines {
		f.WriteString(v)
	}
	f.Close()

	return err
}

func (exec *Executor) AddRecord(domainip string) error {
	f, err := os.OpenFile(DNS_CONFIG_PATH, os.O_WRONLY|os.O_APPEND, 0666)
	defer f.Close()
	if nil != err {
		log.Println(err)
		panic(err)
	}
	f.WriteString("\n" + domainip)
	return err
}

var once2 sync.Once
var Exectucor *Executor

func GetExecutor() *Executor {
	once2.Do(func() {
		Exectucor = &Executor{
		}
	})
	return Exectucor
}
