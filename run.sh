#!/bin/bash

red_background="\e[41m"
white_text="\e[97m"
reset="\e[0m"

# Print disclaimer
echo -e "${red_background}${white_text}警告:${reset} 该脚本可能无法正常运行，因为它缺乏适当的检查。继续运行吗？"
# 检测用户输入是否以N开头，如果是则退出
read -p "输入Y开始执行，或任意字符退出：" userInput
if [[ $userInput != Y ]]; then
    echo "已退出。"
    exit
fi

# 检测文件是否存在
if [ -f "a.exe" ]; then
    # 获取文件a.exe的最后修改日期
    lastModified=$(stat -c %Y "a.exe")

    # 计算距离现在的时间间隔
    timeDifference=$(( $(date +%s) - lastModified ))

    # 转换时间间隔为天、小时、分钟和秒
    days=$(( timeDifference / 86400 ))
    timeDifference=$(( timeDifference % 86400 ))
    hours=$(( timeDifference / 3600 ))
    timeDifference=$(( timeDifference % 3600 ))
    minutes=$(( timeDifference / 60 ))
    seconds=$(( timeDifference % 60 ))

    echo "构建产物距离现在已经过去 $days 天 $hours 小时 $minutes 分钟 $seconds 秒。"
else
    echo "构建产物不存在。请先运行go build"
    exit
fi


# 检测用户输入是否以N开头，如果是则退出
read -p "按回车键开始执行，或输入以N开头的字符退出：" userInput

if [[ $userInput == N* ]]; then
    echo "已退出。"
    exit
fi

# 检测有没有gnome-terminal先

# Command to check existence
command_to_check="gnome-terminal"

# Check if the command exists
if ! type "$command_to_check" &> /dev/null; then
    echo "Error: $command_to_check does not exist."
    exit 1
fi

# If the command exists, continue with your script
echo "$command_to_check exists, continuing with the script..."

# 执行若干节点启动逻辑（这部分需要根据实际情况修改，bash无法直接启动Windows可执行文件）
gnome-terminal -e "./blockchainEmulator.exe -n 1"
gnome-terminal -e "./blockchainEmulator.exe -n 2"
gnome-terminal -e "./blockchainEmulator.exe -c -f"
gnome-terminal -e "./blockchainEmulator.exe -n 0"