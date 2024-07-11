# Linux 常用命令详解

## 文件和目录操作

### ls
```bash
ls [-l -a -h] [linux文件路径]
```
显示文件和目录信息
- `-l`：使用长格式显示
- `-a`：显示所有文件，包括隐藏文件
- `-h`：以人类可读的格式显示文件大小

### cd
```bash
cd [linux目录路径]
```
切换到指定目录

### pwd
```bash
pwd
```
显示当前工作目录的路径

### ~
表示用户主目录

### /
表示根目录

### mkdir
```bash
mkdir [-p] [目录路径]
```
创建目录
- `-p`：递归创建目录

### touch
```bash
touch [文件路径]
```
创建一个空文件

### cat
```bash
cat [文件路径/文件夹路径]
```
显示文件或文件夹内容

### more
```bash
more [文件路径]
```
逐页显示文件内容，按 `q` 退出

### cp
```bash
cp [-r] [源文件路径] [目标文件路径]
```
复制文件或目录
- `-r`：递归复制目录

### mv
```bash
mv [源文件路径] [目标文件路径]
```
移动或重命名文件或目录

### rm
```bash
rm [-r -f] [文件路径]
```
删除文件或目录
- `-r`：递归删除目录
- `-f`：强制删除

#### 危险命令
```bash
rm -rf /*
```
删除根目录，**危险操作**

### which
```bash
which [命令名称]
```
显示命令的绝对路径

### find
```bash
find 起始目录 -name 匹配文件名 -type 匹配文件类型 -size 匹配文件大小(+|- 区分大于小于 [k|m|g]) -perm 匹配权限 -user 匹配所有者 -group 匹配组 -mtime 匹配修改时间 -atime 匹配访问时间 -ctime 匹配创建时间
```
在目录中查找文件

### grep
```bash
grep -n [匹配模式] [文件路径]
```
匹配文件内容
- `-n`：显示匹配行号

### wc
```bash
wc [-lwc] [文件路径]
```
统计文件的行数、单词数和字节数
- `-l`：行数
- `-w`：单词数
- `-c`：字节数

### head
```bash
head [-n] [文件路径]
```
显示文件的前 `n` 行

### tail
```bash
tail [-n -f] [文件路径]
```
显示文件的后 `n` 行
- `-f`：实时显示文件新增内容

## 命令组合和历史

### 管道命令
```bash
| 
```
前一个命令的输出作为后一个命令的输入

### 逻辑与
```bash
&&
```
前一个命令成功后才执行后一个命令

### 逻辑或
```bash
||
```
前一个命令失败后才执行后一个命令

### history
```bash
history
```
显示历史命令

### echo
```bash
echo [内容`命令`]
```
输出内容或执行命令并输出结果

### 重定向
```bash
>> 追加内容到文件末尾
> 覆盖文件内容
```

## 文件编辑

### vi/vim
```bash
vi/vim [文件路径]
```
编辑文件

#### 进入不同模式
- `i`：进入输入模式
- `esc`：退出输入模式
- `:`：进入命令模式

#### 命令模式操作
- `:q!`：强制退出不保存
- `:wq`：保存并退出
- `i`：插入
- `a`：插入
- `o`：新行插入
- `I`：插入行首
- `A`：插入行尾
- `O`：新行插入
- `/`：搜索
- `n`：下一个
- `N`：上一个

#### 底线命令模式操作
- `:w`：保存
- `:set nu`：显示行号
- `:set nonu`：取消显示行号

## 用户和权限管理

### su
```bash
su - 用户名
```
切换用户

### exit
```bash
exit
```
退出当前用户

### sudo
```bash
sudo [命令]
```
以管理员权限执行命令

### visudo
```bash
visudo
```
编辑 `/etc/sudoers` 文件添加用户权限

### 用户和用户组管理
```bash
useradd [-g -d] 用户名
userdel -r 用户名
groupadd 用户组名
groupdel 用户组名
usermod -aG 用户组 用户名
id [用户名]
getent passwd
```
- `useradd`：添加用户，`-g` 指定用户组，`-d` 指定用户主目录
- `userdel`：删除用户，`-r` 递归删除用户目录
- `groupadd`：添加用户组
- `groupdel`：删除用户组
- `usermod`：添加用户到用户组
- `id`：查看用户信息
- `getent passwd`：获取系统所有用户信息

### 权限管理
```bash
chmod [-R] 权限 (u=rwx,g=rwx,o=rx)\777 目标文件或目录
chown [-R] 用户名:组名 目标文件或目录
```
- `chmod`：改变文件或目录的权限，`-R` 递归改变目录权限
- `chown`：改变文件或目录的拥有者和组，`-R` 递归改变目录所有者和组

### 权限表示
```text
drwxrwxrwx
```
第一个字母表示文件类型，后面三个字母表示拥有者权限，接下来的三个字母表示组权限，最后三个字母表示其他用户权限。

## 系统和服务管理

### yum
```bash
yum -y install 软件名称
yum -y remove 软件名称
yum -y search 软件名称
```
- 安装软件
- 卸载软件
- 搜索软件

### systemctl
```bash
systemctl start/stop/restart/reload 服务名称
systemctl enable/disable 服务名称
systemctl status 服务名称
```
- 启动/停止/重启服务
- 开机启动/关闭开机启动
- 查看服务状态

### ln
```bash
ln -s 源文件 目标文件
```
创建软链接

### date
```bash
date +%Y-%m-%d %H:%M:%S
date -d "+1 year"
date -d "-1 day"
date -d "2022-12-31"
```
显示当前日期和时间
- 显示明年日期
- 显示昨天日期
- 显示指定日期

### 设置时区
```bash
rm -f /etc/localtime
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```

### hostname
```bash
hostname
hostnamectl set-hostname 主机名
```
- 显示主机名
- 设置主机名

### ping
```bash
ping [-c num] 主机名
```
网络连通性测试，`-c` 发送的包数量，默认4个

### ifconfig
```bash
ifconfig
```
显示网络接口信息

### wget
```bash
wget [-b] url
```
下载文件，`-b` 后台下载

### curl
```bash
curl -O 文件名 url
```
下载文件，`-O` 下载到本地文件名

### nmap
```bash
nmap 127.0.0.1
```
扫描本机开放端口

### netstat
```bash
netstat -anp | grep ssh
```
扫描本机开放的ssh端口

### ps
```bash
ps [-e -f]
ps -ef | grep ssh
```
显示进程信息
- `-e`：显示所有进程
- `-f`：显示详细信息

### kill
```bash
kill [-9] PID
```
杀死进程，`-9` 强制杀死进程

### top
```bash
top
```
实时显示系统信息

## 环境变量管理

### echo
```bash
echo $PATH
echo ${PATH}ABC
echo $HOME
echo $PWD
```
- 显示环境变量 `PATH`
- 显示环境变量 `PATH` 后追加 `ABC`
- 显示用户主目录
- 显示当前目录路径

### export
```bash
export 变量名=值
export PATH=$PATH:/usr/local/bin
```
临时设置环境变量

### 永久设置

环境变量
```bash
vim /etc/profile
vim /etc/bashrc
source /etc/profile
```
- 编辑 `/etc/profile` 添加环境变量（对所有用户生效）
- 编辑 `/etc/bashrc` 添加环境变量（对当前用户生效）
- 使环境变量生效

## 文件上传下载

### rz/sz
```bash
rz
sz
```
- 上传文件到服务器
- 从服务器下载文件

## 打包和压缩

### tar
```bash
tar [-c -f -v -x -z -C]
```
- `-c`：打包
- `-f`：指定打包文件名
- `-v`：显示过程
- `-x`：解压
- `-z`：gzip压缩
- `-C`：解压到指定目录

#### 示例
```bash
tar -cvf 打包文件名.tar 要打包的文件
tar -xvf 打包文件名.tar
tar -zxvf 打包文件名.tar.gz -C 指定目录
```
- 打包文件到 `tar` 文件
- 解压 `tar` 文件
- 解压 `gzip` 压缩文件到指定目录

### zip/unzip
```bash
zip [-r] 打包文件名.zip 要打包的文件
unzip 打包文件名.zip -d 指定目录
```
- `-r`：递归打包目录
- 解压文件到指定目录

---
