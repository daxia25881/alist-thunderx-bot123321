# 使用GitHub Desktop上传项目

GitHub Desktop是一个用户友好的图形界面工具，非常适合不熟悉命令行的用户使用。

## 步骤指南

### 1. 下载并安装GitHub Desktop
- 访问[GitHub Desktop官网](https://desktop.github.com/)
- 下载适合您操作系统的版本（Windows/macOS）
- 安装应用程序

### 2. 登录您的GitHub账户
- 打开GitHub Desktop
- 使用您的GitHub账户登录

### 3. 创建新仓库
- 点击"File" > "New Repository"（或使用快捷键Ctrl+N）
- 填写仓库名称（如"alist-thunderx-bot"）
- 选择本地路径（即您项目文件所在的文件夹）
- 点击"Create Repository"

### 4. 提交您的代码
- GitHub Desktop会自动显示所有可以添加的文件
- 在左下角添加提交消息，如"Initial commit"
- 点击"Commit to main"（或master，取决于您的默认分支）

### 5. 发布到GitHub
- 提交后，点击右上角的"Publish repository"
- 确认仓库名称和描述
- 确保"Keep this code private"选项是否勾选（视您需要而定）
- 点击"Publish Repository"

### 6. 验证上传
- 完成发布后，您可以在GitHub Desktop中点击"View on GitHub"
- 或直接访问您的GitHub账户页面查看新仓库
- 确认所有文件都已成功上传

## 注意事项
- 确保上传了所有必要的文件：main.go, Dockerfile, go.mod, go.sum, config.json等
- 如果以后需要更新代码，只需在GitHub Desktop中做出更改，提交并点击"Push origin" 