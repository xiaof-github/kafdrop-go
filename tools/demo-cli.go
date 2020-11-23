package maicli

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
    var msgFile string

    app := & cli.App{
      Flags : []cli.Flag {
        &cli.StringFlag{
          Name: "file",
          Usage: "send msg from `FILE`",
          Destination: &msgFile,
        },
      },
      Action: func(c *cli.Context) error {        
        fileName := c.String("file")
        fmt.Println("read file:", fileName)
        filePath := "./"+fileName
        file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
        // File file = os.
        if err != nil {
            fmt.Println("打开文件：", err)
            panic("文件打开失败")
        }
        //读原来文件的内容，并且显示在终端
        reader := bufio.NewReader(file)
        for {
            str, err := reader.ReadString('\n')
            if err == io.EOF {
                break
            }
            fmt.Print(str)
        }

        defer file.Close()
        return nil
      },
    }
  
    err := app.Run(os.Args)
    if err != nil {
      log.Fatal(err)
    }	

  }