# Render部署指南

将您的应用部署到Render非常简单，尤其是当您已经准备好了render.yaml和Dockerfile。

## 步骤指南

### 1. 创建Render账户
- 访问[Render官网](https://render.com/)
- 注册一个新账户或登录已有账户
- 您可以直接使用GitHub账户登录，这样会更方便连接仓库

### 2. 创建新的Web Service
- 登录后，点击Dashboard中的"New +"按钮
- 选择"Web Service"

### 3. 连接GitHub仓库
- 选择"Connect GitHub"选项
- 如果您尚未授权Render访问您的GitHub账户，请按照提示进行授权
- 在仓库列表中找到并选择您的"alist-thunderx-bot"仓库

### 4. 配置Web Service
Render会自动检测到您的Dockerfile和render.yaml，并预填充大部分设置。确认以下配置：

- **名称**: 保持默认或自定义一个名称
- **环境**: Docker
- **分支**: main（或master，取决于您使用的默认分支）
- **区域**: 选择离您或目标用户最近的区域
- **实例类型**: 免费
- **健康检查路径**: / （应该已自动填写）

### 5. 设置环境变量
点击"Advanced"展开高级设置，找到"Environment Variables"部分，添加所有必要的环境变量：

- `BOT_USERNAME`: 您的Alist用户名
- `BOT_PASSWORD`: 您的Alist密码
- `BOT_BASE_URL`: 您的Alist实例URL，例如 http://111.119.215.161:10021/
- `BOT_SEARCH_URL`: 搜索API的URL
- `BOT_OFFLINE_DOWNLOAD_DIR`: 离线下载目录路径，例如 /thunderx/我的云盘
- `BOT_TELEGRAM_TOKEN`: 您的Telegram机器人Token

### 6. 创建并部署
- 点击"Create Web Service"按钮
- Render将开始构建和部署您的应用
- 构建过程可能需要几分钟时间

### 7. 验证部署
- 部署完成后，Render会提供一个URL（形如 https://your-app-name.onrender.com）
- 您可以访问这个URL来验证健康检查（应该显示"Bot is running"）
- 检查Render提供的日志，确认机器人已经成功启动并连接到Telegram API

### 8. 设置Telegram Webhook（可选，但推荐）
如果您希望提高机器人的响应速度，可以设置Webhook：

```
https://api.telegram.org/bot您的BOT_TOKEN/setWebhook?url=https://your-app-name.onrender.com/您的SECRET_PATH
```

## 注意事项
- Render免费计划有一些限制，包括每月750小时的运行时间
- 如果长时间不活动，免费服务会进入休眠状态，第一次访问时可能需要一点时间来启动
- 定期检查日志，确保机器人正常运行 