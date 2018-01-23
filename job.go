package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/linlexing/datelogger"
	"github.com/robfig/cron"
)

const (
	//每次写入最大数量
	batchNum = 500
	//间隔日期
	dayNum = 5
)

var (
	jobs      = cron.New()
	running   = false
	jobRun    = &sync.Mutex{}
	FileMatch = regexp.MustCompile("^R433")
	logOut    = datelogger.NewDateLog("log")
	//系统当前时间
	now = time.Now()
	//存储打开的文件夹
	useFileDirs = new(userFiles)

	//最开始的日期
	up string = "20160102"
)

type Files struct {
	List []os.FileInfo
}
type userFiles struct {
	fileDirs []string
}

func (f *Files) Swap(i, j int) {
	f.List[i], f.List[j] = f.List[j], f.List[i]
}
func (f *Files) Less(i, j int) bool {
	//比较i对应的name是否比j对应的小：是为true
	return -1 == strings.Compare(f.List[i].Name()[12:], f.List[j].Name()[12:])
}
func (f *Files) Len() int {
	return len(f.List)
}
func taskRun() {
	if running {
		return
	}
	running = true
	defer func() {
		running = false
	}()
	dlog.Println("before lock")
	jobRun.Lock()
	defer jobRun.Unlock()
	dlog.Println("start job")
	up = vconfig.Uptime
	err := useFileDirs.openFileDir(vconfig.Filedir, vconfig.Startmonth, vconfig.Stopmonth)
	if err != nil {
		dlog.Error(err)
		return
	}
	err = buildDataFile()
	if err != nil {
		dlog.Error(err)
		return
	}
	//然后开始上传
	//uploadAll()

	dlog.Println("job finished")
}

//创建zip文件，并将templa文件夹中的文件复制到zip文件中
func createNewZipFile() (*os.File, *zip.Writer, error) {
	//添加在workDir后添加out路径
	outPath := filepath.Join(workDir, "out")
	//使用指定的权限和名称创建一个目录，包括任何必要的上级目录，并返回nil，否则返回错误
	if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
		return nil, nil, err
	}
	//确定文件名过程为：
	//out目录中没有同名文件
	//upload目录中也没有同名文件
	var fileName string
	for i := 1; ; i++ {
		fileName = fmt.Sprintf("gsdata_%s_%s_%06d.zip", time.Now().Format("20060102"),
			vconfig.AreaCode, i)
		var not1, not2 bool
		//Stat返回一个描述name指定的文件对象的FileInfo。
		//如果指定的文件对象是一个符号链接，返回的FileInfo描述该符号链接指向的文件的信息，
		//本函数会尝试跳转该链接
		//IsNotExist返回一个布尔值说明该错误是否表示一个文件或目录不存在
		if _, err := os.Stat(filepath.Join(workDir, vconfig.FinishOut, fileName)); os.IsNotExist(err) {
			not1 = true
		} else if err != nil {
			return nil, nil, err
		}
		if _, err := os.Stat(filepath.Join(workDir, "out", fileName)); os.IsNotExist(err) {
			not2 = true
		} else if err != nil {
			return nil, nil, err
		}
		if not1 && not2 {
			break
		}
	}
	//创建文件名为fileName的zip
	file, err := os.Create(filepath.Join(outPath, fileName))
	if err != nil {
		return nil, nil, err
	}
	//得到一个将zip文件写入file的*Writer
	zipw := zip.NewWriter(file)
	//先复制模板文件
	//返回template指定的目录的目录信息的有序列
	files, err := ioutil.ReadDir(filepath.Join(workDir, "template"))
	if err != nil {
		return nil, nil, err
	}
	//将template文件夹里的文件复制到Zip中
	for _, f := range files {
		//使用给出的文件名添加一个文件进zip文件。本方法返回的w是一个io.Writer接口（用于写入新添加文件的内容）
		w, err := zipw.Create(f.Name())
		if err != nil {
			return nil, nil, err
		}
		//ReadFile 从filename指定的文件中读取数据并返回文件的内容
		bys, err := ioutil.ReadFile(filepath.Join(workDir, "template", f.Name()))
		if err != nil {
			return nil, nil, err
		}
		//通过w向文件中写入bys
		if _, err = w.Write(bys); err != nil {
			return nil, nil, err
		}
	}
	//bufio.NewWriter创建一个具有默认大小缓冲、写入w的*Writer。
	return file, zipw, err
}

//向文件中写入一行
func writeLine(w *bufio.Writer, strs []string) error {
	for _, str := range strs {
		if _, err := w.WriteString(str); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return nil
}

//向文件中写入数据
func buildDataFile() error {
	file, zipw, err := createNewZipFile()
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	defer zipw.Close()

	//创建dat文件，并得到一个向该文件写的流
	datFile, err := zipw.Create(fmt.Sprint("ent_info", ".dat"))
	if err != nil {
		log.Println(err)
		return err
	}
	//得到一个向dat文件写的缓冲流
	datw := bufio.NewWriter(datFile)
	defer datw.Flush()
	//共写多少行
	icount := 0
	//写入文件
	returnTime, err := readFile(vconfig.Uptime, func(i int, datas [][]string) error {
		if len(datas) > 0 {
			icount += i
			dlog.Println("rownum:", i, "write", i, "rows,total:", icount)
		}
		for _, line := range datas {
			if err := writeLine(datw, line); err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		logOut.Println(err)
		return err
	}
	//如果返回时间等于当前时间，退出程序
	if returnTime < now.Format("20060102") {
		buildDataFile()
	}
	return nil
}

//获得需打开的文件夹路径:具体到月份文件夹,读一段时间的文件夹，读开始不读结束的月份，如无开始或无结束则读最新月的
func (user *userFiles) openFileDir(fileDir, startMonth, stopMonth string) error {
	if len(stopMonth) == 0 {
		stopMonth = "999999"
	}
	if len(startMonth) == 0 {
		startMonth = "999999"
	}
	file, err := os.Open(fileDir)
	if err != nil {
		logOut.Println(err)
		return err
	}
	defer file.Close()
	//获得file文件里的所有文件对象
	info, err := file.Readdir(0)
	if err != nil {
		logOut.Println(err)
		return err
	}
	dirs := []string{}
	for _, v := range info {
		if v.Name() >= startMonth && v.Name() < stopMonth {
			dirs = append(dirs, filepath.Join(fileDir, v.Name()))
		} else if v.Name() < startMonth && v.Name() <= stopMonth && v.Name() == now.Format("200601") {
			dirs = append(dirs, filepath.Join(fileDir, v.Name()))
		}
	}
	user.fileDirs = dirs
	return nil
}

//从文件夹里读取xml文件，
func readFile(upTime string, cd func(num int, datas [][]string) error) (string, error) {
	returnTime := now.Format("20060102")
	fileDirs := useFileDirs.fileDirs
	//结束日期
	stopTime := returnTime
	for _, dir := range fileDirs {
		//读取文件夹
		file, err := os.Open(dir)
		if err != nil {
			logOut.Println(err)
			return returnTime, err
		}
		defer file.Close()
		info, err := file.Readdir(0)
		if err != nil {
			logOut.Println(err)
			return returnTime, err
		}
		//判断是否有上传日期限制
		if len(upTime) < 8 {
			up = now.Format("20060102")
		}
		stopTime, err = turnTime(up, dayNum)
		if err != nil {
			logOut.Println(err)
			return returnTime, err
		}
		upTime, err = turnTime(stopTime, -dayNum)
		if err != nil {
			logOut.Println(err)
			return returnTime, err
		}
		//存放所需文件
		outList := []os.FileInfo{}
		for _, v := range info {
			//判断文件名是否符合要求
			if FileMatch.MatchString(v.Name()) && v.Name()[12:20] >= upTime && v.Name()[12:20] < stopTime && v.Size() > 0 {
				outList = append(outList, v)
				up = stopTime
			}
		}
		flist := &Files{outList}
		//排序
		sort.Sort(flist)
		//已写文件数量
		icount := 0
		//已读出的文件内容
		datas := [][]string{}
		//判断是否有文件
		if flist.Len() > 0 {
			//读取每个文件
			for _, oneFile := range flist.List {
				icount++
				strs, err := readOneFile(filepath.Join(dir, oneFile.Name()))
				if err != nil {
					logOut.Println(err)
					return returnTime, err
				}
				datas = append(datas, strs...)
				//判断是否读入batchNum个文件
				if icount%batchNum == 0 {
					//将文件写入
					if err := cd(icount, datas); err != nil {
						logOut.Println(err)
						return returnTime, err
					}
					icount = 0
					datas = nil
				}

			}
			if icount > 0 {
				if err := cd(icount, datas); err != nil {
					logOut.Println(err)
					return returnTime, err
				}
			}
			return stopTime, nil
		}
	}
	return stopTime, nil
}

//给指定的字符串时间beforDay，加上或减去一定的天数，返回正确的字符串时间afterDay，num大于零为加
func turnTime(beforDay string, num int) (afterDay string, err error) {
	oldTime, err := time.Parse("20060102", beforDay)
	if err != nil {
		fmt.Print(err)
		return
	}
	dayNum := num * 24
	//解析一个时间段,得到可供时间类型相加减的Duration
	d, err := time.ParseDuration(fmt.Sprintf("%vh", dayNum))
	if err != nil {
		return
	}
	newTime := oldTime.Add(d)
	afterDay = newTime.Format("20060102")
	return
}
