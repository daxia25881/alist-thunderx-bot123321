# Alist ThunderX Telegram Bot

这是一个Telegram机器人，用于将磁力链接添加到Alist的ThunderX离线下载队列。

## 功能

- 处理用户发送的磁力链接
- 根据番号搜索并获取磁力链接
- 将磁力链接添加到离线下载队列
- 查询离线下载目录内容

## 部署

### 环境变量

可以通过环境变量或config.json配置文件来配置机器人：

- `BOT_USERNAME`: Alist用户名
- `BOT_PASSWORD`: Alist密码
- `BOT_BASE_URL`: Alist实例的基础URL
- `BOT_SEARCH_URL`: 搜索API的URL
- `BOT_OFFLINE_DOWNLOAD_DIR`: 离线下载目录路径
- `BOT_TELEGRAM_TOKEN`: Telegram机器人Token

### 使用Docker部署

1. 构建Docker镜像：
   ```
   docker build -t alist-thunderx-bot .
   ```

2. 运行Docker容器：
   ```
   docker run -d --name alist-bot \
     -e BOT_USERNAME=your_alist_username \
     -e BOT_PASSWORD=your_alist_password \
     -e BOT_BASE_URL=http://your_alist_instance:port/ \
     -e BOT_SEARCH_URL=http://your_search_api_url/api/search?keyword= \
     -e BOT_OFFLINE_DOWNLOAD_DIR=/thunderx/我的云盘 \
     -e BOT_TELEGRAM_TOKEN=your_telegram_bot_token \
     alist-thunderx-bot
   ```

### 在Render上部署

1. Fork此仓库到您的GitHub账户
2. 在Render上创建一个新的Web Service
3. 连接您的GitHub仓库
4. Render会自动检测Dockerfile并构建镜像
5. 在环境变量部分设置上述配置项
6. 点击部署

## 使用方法

1. 向机器人发送 `/start` 启动
2. 发送磁力链接（以 magnet:? 开头）或番号
3. 机器人会自动处理并添加到离线下载队列 