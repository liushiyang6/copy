package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

/**
执行此程序后 首次复制的文件夹目录为目标目录 后续复制的路径所有文件将自动复制到目标目录
*/
func main() {
	// 初始化剪切板
	var initShearPlate string
	// 目标路径
	var targetShearPlate string
	// 是否初始化目标路径
	var isInitTargetDirectory = false
	initShearPlate, _ = clipboard.ReadAll()
	fmt.Println("初始化剪切板:", initShearPlate)

	for i := 0; i > -1; i++ {
		localShearPlate, _ := clipboard.ReadAll()
		if initShearPlate != localShearPlate {
			initShearPlate = localShearPlate
			fmt.Println("检测到剪切板变化:", initShearPlate)
			if IsExist(initShearPlate) && IsDir(initShearPlate) {
				if !isInitTargetDirectory {
					targetShearPlate = initShearPlate
					fmt.Println("新复制文本为文件夹路径,设置此路径为目标文件夹成功!", targetShearPlate)
					isInitTargetDirectory = true
				} else {
					fmt.Println("开始将此路径所有文件拷贝至目标目录:", targetShearPlate)
					_ = Dir(initShearPlate, targetShearPlate)
					fmt.Println("拷贝成功!")
				}
			}
		}
		time.Sleep(time.Second / 2)
	}
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

// File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// Dir copies a whole directory recursively
func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
