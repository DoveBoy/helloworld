package common

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	goQrcode "github.com/skip2/go-qrcode"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"image"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
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
	s := []string{
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
	str := ""
	for i := 1; i <= length; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		str += s[r.Intn(len(s)-1)]
	}
	return str
}

func Substr(s string, start, end int) string {
	strRune := []rune(s)
	if start == -1 {
		return string(strRune[:end])
	}
	if end == -1 {
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
	if runtime.GOOS == "windows" { //windows
		cmd := exec.Command("cmd", "/c", "rundll32.exe", "C:\\Windows\\System32\\shimgvw.dll,ImageView_FullscreenA", qrPath)
		_ = cmd.Start()
		//扫码后二维码自动删除，自动关闭照片查看器
		go func() {
			for {
				time.Sleep(time.Duration(1) * time.Second)
				if !Exists(qrPath) {
					_ = exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(cmd.Process.Pid)).Run()
					break
				}
			}
		}()
	} else if runtime.GOOS == "darwin" { //Macos
		QrWithBackground(qrPath)
		cmd := exec.Command("open", qrPath)
		_ = cmd.Start()
		//扫码后二维码自动删除，自动关闭照片查看器
		go func() {
			for {
				time.Sleep(time.Duration(1) * time.Second)
				if !Exists(qrPath) {
					_ = exec.Command("pkill","-f","Preview").Run()
					break
				}
			}
		}()
	} else {
		//linux或者其他系统
		file, _ := os.Open(qrPath)
		img, _, _ := image.Decode(file)
		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
		qrReader := qrcode.NewQRCodeReader()
		res, _ := qrReader.Decode(bmp, nil)
		//输出控制台
		qr, _ := goQrcode.New(res.String(), goQrcode.High)
		fmt.Println(qr.ToSmallString(false))
	}
}

func QrWithBackground(path string) {
	bg, _ := Asset("bg.png")
	background, _ := imaging.Decode(bytes.NewReader(bg))
	qr, _ := imaging.Open(path)
	dst := imaging.Paste(background, qr, image.Pt(1555, 500))
	imaging.Save(dst, path)
}

//指定位数随机数
func RandomNumber(len int) string {
	var container string
	var str = "0123456789"
	b := bytes.NewBufferString(str)
	length := b.Len()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		container += string(str[rand.Intn(length)])
	}
	return container
}
