package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// 需要在项目根目录运行该脚本

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func log(format string, a ...any) {
	_, err := os.Stdout.Write([]byte(fmt.Sprintf(format, a...)))
	check(err)
}

func main() {
	var err error

	err = os.RemoveAll("./test/log")
	check(err)
	err = os.MkdirAll("./test/log", 0o755)
	check(err)

	// 生成代码
	{
		logPath := "./test/log/generate.log"
		logFile, err := os.Create(logPath)
		check(err)
		defer func() {
			err := logFile.Close()
			check(err)
		}()

		log(">>> 生成代码\n")
		log("日志：%s\n", logPath)

		cmd := exec.Command("go", "run", "cmd/n-cl/main.go", "./test")
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		err = cmd.Run()
		if err != nil {
			log("异常：%s\n", err)
		}
	}

	// 初始化测试目录
	{
		logPath := "./test/log/init.log"
		logFile, err := os.Create(logPath)
		check(err)
		defer func() {
			err := logFile.Close()
			check(err)
		}()

		log(">>> 初始化项目\n")
		log("日志：%s\n", logPath)

		{
			log("golang 项目初始化\n")
			cmd := exec.Command("go", "mod", "tidy")
			cmd.Dir = "./test/golang"
			cmd.Stdout = logFile
			cmd.Stderr = logFile
			err = cmd.Run()
			if err != nil {
				log("异常：%s\n", err)
			}
		}
		{
			log("python 项目初始化\n")
			cmd := exec.Command("bash", "-c", "source .venv/bin/activate && pip install -r requirements.txt")
			cmd.Dir = "./test/python"
			cmd.Stdout = logFile
			cmd.Stderr = logFile
			err = cmd.Run()
			if err != nil {
				log("异常：%s\n", err)
			}
		}
		{
			log("typescript 项目初始化\n")
			cmd := exec.Command("pnpm", "i")
			cmd.Dir = "./test/typescript"
			cmd.Stdout = logFile
			cmd.Stderr = logFile
			err = cmd.Run()
			if err != nil {
				log("异常：%s\n", err)
			}
		}
	}

	// 测试
	i := 0

	for serverName, serverCmd := range map[string]func(io.Writer) func(){
		"golang": golangServer,
		"python": pythonServer,
	} {
		for clientName, clientCmd := range map[string]func(io.Writer){
			"golang":     golangClient,
			"typescript": typescriptClient,
		} {
			i += 1

			log(">>> 测试：%d\n", i)

			sLogPath := fmt.Sprintf("./test/log/test%d-server-%s.log", i, serverName)
			sLogFile, err := os.Create(sLogPath)
			check(err)
			defer func() {
				err := sLogFile.Close()
				check(err)
			}()

			cLogPath := fmt.Sprintf("./test/log/test%d-client-%s.log", i, clientName)
			cLogFile, err := os.Create(cLogPath)
			check(err)
			defer func() {
				err := cLogFile.Close()
				check(err)
			}()

			// 启动服务端
			time.Sleep(time.Second * 5)
			close := serverCmd(sLogFile)

			// 客户端
			time.Sleep(time.Second * 15)
			clientCmd(cLogFile)

			// 停止服务端
			close()
		}
	}

}

func golangServer(logFile io.Writer) (cancel func()) {
	c := make(chan struct{})

	cmd := exec.Command("go", "build", "-o", "./bin/serer", "./cmd/server/main.go")
	cmd.Dir = "./test/golang"
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err := cmd.Run()
	check(err)

	cmd = exec.Command("./bin/serer")
	cmd.Dir = "./test/golang"
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Start()
	check(err)

	go func() {
		<-c
		err = cmd.Process.Kill()
		check(err)
	}()

	return func() {
		close(c)
	}
}

func pythonServer(logFile io.Writer) (cancel func()) {
	c := make(chan struct{})

	cmd := exec.Command(".venv/bin/python", "server.py")
	cmd.Dir = "./test/python"
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err := cmd.Start()
	check(err)

	go func() {
		<-c
		err = cmd.Process.Kill()
		check(err)
	}()

	return func() {
		close(c)
	}
}

func golangClient(logFile io.Writer) {
	cmd := exec.Command("go", "run", "cmd/client/main.go")
	cmd.Dir = "./test/golang"
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err := cmd.Run()
	check(err)
}

func typescriptClient(logFile io.Writer) {
	cmd := exec.Command("pnpm", "run", "client")
	cmd.Dir = "./test/typescript"
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err := cmd.Run()
	check(err)
}
