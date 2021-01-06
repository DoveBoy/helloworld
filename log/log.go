package log

import (
	"github.com/gookit/color"
	"io"
	"log"
	"os"
	"time"
)

//重写log标准库，需要多少个方法就加多少个

var file = "./logs/jd_seckill_" + time.Now().Format("20060102") + ".log"

//将日志同时输出到控制台和文件
func Println(v ...interface{}) {
	logFile, logErr := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if logErr != nil {
		panic(logErr)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	//log.SetPrefix("[jd_seckill]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Println(v...)
}

//将日志同时输出到控制台和文件
func ColorPrintln(color2 color.Color, v ...interface{}) {
	logFile, logErr := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if logErr != nil {
		panic(logErr)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	//log.SetPrefix("[jd_seckill]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Println(v...)
	color2.Light().Println(v...)
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Success(v ...interface{}) {
	ColorPrintln(color.Green, v...)
}

func Info(v ...interface{}) {
	ColorPrintln(color.LightCyan, v...)
}

func Warning(v ...interface{}) {
	ColorPrintln(color.Yellow, v...)
}

func Error(v ...interface{}) {
	ColorPrintln(color.FgLightRed, v...)
	os.Exit(0)
}
