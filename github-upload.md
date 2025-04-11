# GitHub上传指南

## 使用Git命令行上传代码

### 1. 安装Git（如果尚未安装）
- Windows: 从[Git官网](https://git-scm.com/download/win)下载安装程序
- Mac: 通过Terminal运行`brew install git`（需要Homebrew）或从[Git官网](https://git-scm.com/download/mac)下载
- Linux: 使用包管理器，例如`sudo apt-get install git`（Ubuntu/Debian）或`sudo yum install git`（CentOS/RHEL）

### 2. 配置Git
```bash
git config --global user.name "您的GitHub用户名"
git config --global user.email "您的GitHub电子邮件"
```

### 3. 初始化本地仓库
在您的项目文件夹中运行:
```bash
cd /path/to/alist_thunderx_bot
git init
```

### 4. 添加文件到暂存区
```bash
git add .
```

### 5. 提交更改
```bash
git commit -m "Initial commit"
```

### 6. 创建GitHub仓库
- 访问[GitHub](https://github.com)并登录
- 点击右上角"+"图标，选择"New repository"
- 输入仓库名称（如"alist-thunderx-bot"）
- 不要初始化README，.gitignore或license
- 点击"Create repository"

### 7. 连接本地仓库到GitHub
```bash
git remote add origin https://github.com/您的用户名/alist-thunderx-bot.git
```

### 8. 推送代码到GitHub
```bash
git push -u origin main
```
注意：如果您的默认分支是`master`而不是`main`，请使用：
```bash
git push -u origin master
```

### 9. 输入GitHub凭据
当提示时，输入您的GitHub用户名和密码或个人访问令牌

### 10. 验证上传
访问您的GitHub仓库页面，确认所有文件已成功上传 