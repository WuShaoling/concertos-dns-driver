package executor

import (
	"sync"
	"os"
	"bufio"
	"io"
	"log"
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type Executor struct {
}

const DNS_CONFIG_PATH = "/etc/dnsmasq.hosts"

func (exec *Executor) ReadLines() ([]string, error) {

	var lines []string

	fi, err := os.Open(DNS_CONFIG_PATH)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if len(a) != 0 {
			lines = append(lines, string(a))
		}
	}

	return lines, nil
}

func (exec *Executor) DeleteRecord(domain string) error {
	// read all
	lines, err := exec.ReadLines()
	if err != nil {
		return err
	}

	// delete file
	if err = os.Remove(DNS_CONFIG_PATH); nil != err {
		log.Println(err)
		return err
	}

	// create and write to file
	f, err := os.OpenFile(DNS_CONFIG_PATH, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, v := range lines {
		if strings.Contains(v, domain) {
			continue
		}
		f.WriteString(v + "\n")
	}

	// restart dnsmasq
	err = exec.restartDnsmasq()

	return err
}

func (exec *Executor) AddRecord(domainip string) error {
	f, err := os.OpenFile(DNS_CONFIG_PATH, os.O_WRONLY|os.O_APPEND, 0666)
	if nil != err {
		log.Println(err)
	}
	defer f.Close()
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
