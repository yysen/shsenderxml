package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/signal"

	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/linlexing/datelogger"
)

var (
	dlog    *datelogger.DateLogger
	workDir string
)

func uploadAll() {
	files, err := ioutil.ReadDir(filepath.Join(workDir, "out"))
	if err != nil {
		dlog.Error(err)
		return
	}
	for _, one := range files {
		filename := filepath.Join(workDir, "out", one.Name())
		dlog.Println("file:", filename, "start upload...")
		if err = doUpload(vconfig.URL, filename,
			filepath.Join(workDir, vconfig.FinishOut), vconfig.UserName,
			vconfig.Password); err != nil {
			dlog.Error(err)
			break
		}
		dlog.Println("file:", filename, "uploaded")
	}

}
func main() {
	build := flag.Bool("buildnow", false, "build immediate")

	flag.Parse()

	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
	//获得当前文件的路径
	str, err := os.Executable()
	if err != nil {
		panic(err)
	}
	//获取路径的中没有文件的部分，即文件所在的文件夹
	workDir = filepath.Dir(str)
	dlog = datelogger.NewDateLog(filepath.Join(workDir, "log"))
	lend := make(chan bool)
	if err := readConfig(filepath.Join(workDir, "config.yaml")); err != nil {
		panic(err)
	}
	//凌晨两点执行
	if err := jobs.AddFunc("0 0 2 * * *", taskRun); err != nil {
		panic(err)
	}
	jobs.Start()
	for _, one := range jobs.Entries() {
		fmt.Printf("job next:%s\n", one.Next)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	//uploadAll()
	dlog.Println("start service...")
	if *build {
		//马上执行一次
		go taskRun()
	}
	go func() {
		for range c {
			dlog.Info("received ctrl+c,wait back job finished...")
			jobRun.Lock()
			dlog.Info("success shutdown")
			lend <- true
			break
		}
	}()
	//blockforever
	<-lend
}
func 认证(strURL, userName, password string) (string, error) {
	u, err := url.Parse(strURL)
	if err != nil {
		log.Panic(err)
	}
	u.Path = path.Join(u.Path, "login.html")
	q := u.Query()
	q.Add("username", userName)
	q.Add("password", password)
	q.Add("json", "true")
	u.RawQuery = q.Encode()

	rs, err := http.Post(u.String(), "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		return "", fmt.Errorf("url:%s,err:%s", u.String(), err.Error())
	}
	if rs.ContentLength == 0 {
		return "", fmt.Errorf("respon nil")
	}
	defer rs.Body.Close()
	rev := map[string]interface{}{}
	bys, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(bys, &rev); err != nil {

		return "", fmt.Errorf("res:%s,err:%s", string(bys), err)
	}
	if code := rev["message"].(map[string]interface{})["code"].(string); code != "000000" {
		return "", fmt.Errorf(code)
	}
	return rev["record"].(map[string]interface{})["authorization"].(string), nil
}
func hash_file_md5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

func Upload(url, file, austr string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("uploadFile", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", austr)

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	bys, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	rev := map[string]interface{}{}
	if err = json.Unmarshal(bys, &rev); err != nil {
		log.Panic(err)
	}
	if code := rev["message"].(map[string]interface{})["code"].(string); code != "000000" {
		err = fmt.Errorf("code:%s,message:%s", code, rev["message"].(map[string]interface{})["info"].(string))
	}
	return
}

func 上传(strURL, strSrcFile, strFinishOut, austr string) error {
	u, err := url.Parse(strURL)
	if err != nil {
		log.Panic(err)
	}
	md5Str, err := hash_file_md5(strSrcFile)
	if err != nil {
		log.Panic(err)
	}
	fname := filepath.Base(strSrcFile)
	fname = fname[:len(fname)-len(filepath.Ext(fname))]
	u.Path = path.Join(u.Path, "upload/fileupload/", fname, md5Str+".json")
	if err = Upload(u.String(), strSrcFile, austr); err != nil {
		return err
	}
	return nil
}
func doUpload(strURL, strSrcFile, strFinishOut, userName, password string) error {
	//认证
	austr, err := 认证(strURL, userName, password)
	if err != nil {
		return err
	}
	if err = 上传(strURL, strSrcFile, strFinishOut, austr); err != nil {
		return err
	}
	return os.Rename(strSrcFile, filepath.Join(strFinishOut, filepath.Base(strSrcFile)))
}
