package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	rand2 "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"time"
)

const (
	B int64 = 1 << (iota * 10)
	KB
	MB
	GB
	TB
)

func MainDefer() {
	if e := recover(); e != nil {
		pwd, _ := os.Getwd()
		fd, e2 := os.Create(pwd + "/" + time.Now().Format("20060102030405") + ".stack")
		if e2 == nil && fd != nil {
			e3 := e.(error)
			if e3 != nil {
				fd.WriteString("recover:" + e3.(error).Error() + "\n")
			}
			fmt.Println(pwd + "/" + time.Now().Format("20060102030405") + ".stack")
			fd.Write(debug.Stack())
			fd.Close()
		} else {
			log.Println("error:", e, e2)
			log.Println("error:", string(debug.Stack()))
		}
	}
}

func Md5(data []byte) []byte {
	x := md5.Sum(data)
	return x[:16]
}

func MD5(data []byte) string {
	r := md5.Sum(data)
	return hex.EncodeToString(r[:])
}

func HmacMD5(message []byte, key []byte) string {
	mac := hmac.New(md5.New, []byte(key))
	mac.Write((message))
	return hex.EncodeToString(mac.Sum(nil))
}

func Sha1(data []byte) string {
	tmp := sha1.Sum(data)
	return hex.EncodeToString(tmp[:])
}

func Sha1Hmac(data, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return fmt.Sprintf("%x", mac.Sum(nil))
}
func Sha256Hmac(data, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func Base64_Encode(data []byte) []byte {
	var buf = new(bytes.Buffer)
	w := base64.NewEncoder(base64.StdEncoding, buf)
	w.Write(data)
	return buf.Bytes()
}

func GetRanStr(n int) string {
	var set = "0123456789abcdefghijklmnopqrstuvwxyz"

	var buf []byte
	for i := 0; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		ran := rand.Intn(len(set))
		buf = append(buf, set[ran])
	}
	return string(buf)
}

func RandomBytes(n int) []byte {
	var buffer bytes.Buffer
	dictionary := "0123456789abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(dictionary))
		buffer.WriteByte(dictionary[index])
	}
	return buffer.Bytes()
}

func RandIntn(n int) int {
	rand.Seed(time.Now().UnixMicro())
	return rand.Intn(n)
}

func Sha256Str(src []byte) string {
	h := sha256.New()
	h.Write(src)
	return hex.EncodeToString(h.Sum(nil))
}

func GetRandomToken(length int) string {
	r := make([]byte, length)
	io.ReadFull(rand2.Reader, r)
	return hex.EncodeToString(r)
}

var BusyBox bool

// RunCommand run command without shell
func RunCommand2(name string, args ...string) ([]byte, error) {
	var cmd *exec.Cmd
	if BusyBox {
		newArgs := []string{name}
		newArgs = append(newArgs, args...)
		cmd = exec.Command("busybox", newArgs...)
	} else {
		cmd = exec.Command(name, args...)
	}
	log.Println("debug:", "RunCommand:", name, args)
	return cmd.CombinedOutput()
}

func RunCommand(cmd string) ([]byte, error) {
	c := exec.Command("sh", "-c", cmd)
	log.Println("debug:", "RunCommand:", cmd)
	return c.CombinedOutput()
}

// ExecCommand exec command failed
func ExecCommand(name string, args ...string) (stdout []byte, stderr []byte, e error) {
	log.Println("debug:", "run command:", name, args)
	var outbuf, errbuf bytes.Buffer
	var cmd *exec.Cmd
	if BusyBox {
		newArgs := []string{name}
		newArgs = append(newArgs, args...)
		cmd = exec.Command("busybox", newArgs...)
	} else {
		cmd = exec.Command(name, args...)
	}
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.Bytes()
	stderr = errbuf.Bytes()
	if err != nil {
		var exitError *exec.ExitError
		var pathError *os.PathError
		if errors.As(err, &exitError) {
			e = exitError
		} else if errors.As(err, &pathError) {
			e = pathError
		} else {
			e = err
		}
	}
	return
}

func IsExistFile(f string) bool {
	info, e := os.Stat(f)
	if os.IsNotExist(e) {
		return false
	}
	if os.IsPermission(e) {
		return false
	}
	if info == nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}
func IsExist(f string) bool {
	_, e := os.Stat(f)
	if e != nil {
		if os.IsNotExist(e) {
			return false
		}
		if os.IsPermission(e) {
			return false
		}
	}
	return true
}

// IsExistDir check f is exist dir
func IsExistDir(f string) bool {
	info, e := os.Stat(f)
	if e != nil {
		if os.IsNotExist(e) {
			return false
		}
		if os.IsPermission(e) {
			return false
		}
	}
	if info == nil {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true

}

func Page(count int, page, pagesize int) (pages, start, end int) {
	if count <= 0 {
		pages = 0
		start = 0
		end = 0
		return
	}

	if page < 0 {
		page = 0
	}

	if pagesize < 0 {
		pagesize = 20
	}

	pages = int(math.Ceil(float64(count) / float64(pagesize)))
	if page > pages {
		page = pages
	}
	start = page * pagesize
	end = start + pagesize

	if end > count {
		end = count
	}
	return
}

func LikeEscape(s string) string {
	count := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '%' || c == '_' {
			count++
		}
	}

	if count == 0 {
		return s
	}

	required := len(s) + count
	t := make([]byte, required)
	j := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '%' || c == '_' {
			t[j] = '\\'
			j++
		}
		t[j] = c
		j++
	}
	return string(t)
}

func PathSize(path ...string) (sum int64) {
	for i := range path {
		filepath.Walk(path[i], func(path string, info fs.FileInfo, err error) error {
			sum += info.Size()
			return err
		})
	}

	return
}
