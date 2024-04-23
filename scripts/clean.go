package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// 提示用户输入 Y/n 决定是否运行脚本
	fmt.Print("是否运行清理脚本？(输入 Y 确认，输入其他任意字符取消): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm == "Y" || confirm == "y" {
		// 创建一个列表 dirToRemove，包含 record、log、result 三个文件夹。
		dirToRemove := []string{"record", "log", "result"}

		// 遍历列表
		for _, dir := range dirToRemove {
			dirPath := filepath.Join("文件夹路径", dir)

			// 检测路径是否存在
			if _, err := os.Stat(dirPath); err == nil {
				// 如果存在，则删除文件夹
				err := os.RemoveAll(dirPath)
				if err == nil {
					fmt.Printf("文件夹 %s 被删除。\n", dir)
				} else {
					fmt.Printf("无法删除文件夹 %s：%v\n", dir, err)
				}
			} else if os.IsNotExist(err) {
				fmt.Printf("文件夹 %s 不存在，不用删除。\n", dir)
			} else {
				fmt.Printf("无法访问文件夹 %s：%v\n", dir, err)
			}
		}
	} else {
		fmt.Println("取消运行清理脚本。")
	}
}
