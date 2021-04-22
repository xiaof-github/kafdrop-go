package agent

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

func main() {

	for true {
		isRunning := isExeRunning(
			"abc.EXE",
			"abc.EXE",
		)
		println("word是否在运行:%t", isRunning)
		time.Sleep(time.Second * 2)
	}

}

// strKey 检索的关键词 如WINWORD.EXE、WINWORD
//strExeName exe的完整名字
func isExeRunning(strKey string, strExeName string) bool {
	buf := bytes.Buffer{}
	cmd := exec.Command("wmic", "process", "get", "name,executablepath")
	cmd.Stdout = &buf
	cmd.Run()

	cmd2 := exec.Command("findstr", strKey)
	cmd2.Stdin = &buf
	data, err := cmd2.CombinedOutput()
	if err != nil && err.Error() != "exit status 1" {
		return false
	}

	strData := string(data)
	if strings.Contains(strData, strExeName) {
		return true
	} else {
		return false
	}
}
