package fs_util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}

	return true
}

func Sha256(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetFileList(dir string) ([]string, error) {
	var fileList []string

	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// 遍历目录项
	for _, entry := range entries {
		// 排除文件夹，只保留文件
		if !entry.IsDir() {
			fileList = append(fileList, entry.Name())
		}
	}

	return fileList, nil
}
