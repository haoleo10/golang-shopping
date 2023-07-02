package csv_helper

import (
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
)

//从给定的FileHeader读信息，返回二维数组
func ReadCsv(fileHeader *multipart.FileHeader) ([][]string, error) {
	//打开文件
	f, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	//关闭文件
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			log.Print(err)
		}
	}(f)
//读取文件，new一个reader
	reader := csv.NewReader(f)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("err")
		return nil, err
	}
	var result [][]string
//遍历第一二列
	for _, line := range lines[1:] {

		data := []string{line[0], line[1]}
		result = append(result, data)
	}

	return result, nil
}
