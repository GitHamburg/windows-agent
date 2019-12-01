package fileop

import (
	"os"
	"path"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// 判断所给路径是否为文件夹

func IsDir(path string) bool {

	s, err := os.Stat(path)

	if err != nil {

		return false

	}

	return s.IsDir()

}

// 判断所给路径是否为文件

func IsFile(path string) bool {

	return !IsDir(path)

}

func WriteBytes(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

func WriteString(filePath string, s string) (int, error) {
	return WriteBytes(filePath, []byte(s))
}
