package executor

import (
	"sync"
	"os"
	"bufio"
	"io"
	"log"
	"strings"
	"bytes"
	"errors"
	"os/exec"
)

type Executor struct {
}

const DNS_CONFIG_PATH = "/etc/dnsmasq.hosts"
//const DNS_CONFIG_PATH = "a.txt"

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

	err = exec.restartDnsmasq()

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

	err = exec.restartDnsmasq()
	return err
}

func (exe *Executor) restartDnsmasq() error {
	cmd := exec.Command("/bin/bash", "-c", " service dnsmasq restart")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if nil != err {
		log.Println(err)
		return errors.New(string(out.Bytes()))
	}
	log.Println("RestartDnsmasq : ", string(out.Bytes()))
	return nil
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
