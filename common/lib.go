package common

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/CodyGuo/win"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	goQrcode "github.com/skip2/go-qrcode"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func NewRandStr(length int) string {
	s:=[]string{
		"a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l",
		"m", "n", "o", "p", "q", "r",
		"s", "t", "u", "v", "w", "x",
		"y", "z", "A", "B", "C", "D",
		"E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P",
		"Q", "R", "S", "T", "U", "V",
		"W", "X", "Y", "Z",
	}
	str:=""
	for i:=1;i<=length;i++  {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		str+=s[r.Intn(len(s)-1)]
	}
	return str
}

func Substr(s string,start,end int) string {
	strRune:=[]rune(s)
	if start==-1 {
		return string(strRune[:end])
	}
	if end==-1 {
		return string(strRune[start:])
	}
	return string(strRune[start:end])
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func OpenImage(qrPath string) {
	if "windows" == runtime.GOOS { // Windows系统
		cmd := "cmd /c rundll32.exe C:\\Windows\\System32\\shimgvw.dll,ImageView_Fullscreen " + qrPath
		if err := ExecRun(cmd); err != nil {
			log.Println(cmd)
			log.Fatal(err)
		}
	} else { // 非Windows系统(Linux等)输出控制台
		//解码二维码
		file, _ := os.Open(qrPath)
		img, _, _ := image.Decode(file)
		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
		qrReader := qrcode.NewQRCodeReader()
		res, _ := qrReader.Decode(bmp, nil)

		//输出控制台
		qr, err := goQrcode.New(res.String(), goQrcode.High)
		if err != nil {
			log.Println("二维码获取成功，请打开图片用京东APP扫描")
		}
		fmt.Println(qr.ToSmallString(false))
	}
}

func Hour2Unix(hour string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormatStr, time.Now().Format(DateFormatStr) + " " + hour, time.Local)
}

func ExecRun(cmd string) error {
	lpCmdLine := win.StringToBytePtr(cmd)
	ret := win.WinExec(lpCmdLine, win.SW_HIDE)
	if ret <= 31 {
		return errors.New(winExecError[ret])

	}
	return nil
}
