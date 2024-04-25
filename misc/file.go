package misc

import (
	"fmt"
	"os"
)

func CreateUniqueFolder(folderName string) (string, error) {
	// Check if the folder already exists
	_, err := os.Stat(folderName) // 先判断一下目录存不存在
	if os.IsNotExist(err) {
		// Folder doesn't exist, create it
		err := os.Mkdir(folderName, 0755) // 不存在的话那很好。直接创建。否则继续。
		if err != nil {
			return "", err
		}
		return folderName, nil
	}

	// Folder already exists, find a unique name
	for i := 1; ; i++ {
		newFolderName := fmt.Sprintf("%s(%d)", folderName, i) // 看看folder(i)存不存在。不存在的话就好。否则继续。
		_, err := os.Stat(newFolderName)
		if os.IsNotExist(err) {
			// Unique folder name found, create it
			err := os.Mkdir(newFolderName, 0755)
			if err != nil {
				return "", err
			}
			return newFolderName, nil
		}
	}
}
