package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Config 定义配置结构体
type Config struct {
    Username           string `json:"username"`
    Password           string `json:"password"`
    BaseURL            string `json:"base_url"`
    SearchURL          string `json:"search_url"`
    OfflineDownloadDir string `json:"offline_download_dir"`
    TelegramToken      string `json:"telegram_token"`
}

// 全局配置
var config Config

// 加载配置
func loadConfig() error {
    // 优先从环境变量加载
    config.Username = os.Getenv("BOT_USERNAME")
    config.Password = os.Getenv("BOT_PASSWORD")
    config.BaseURL = os.Getenv("BOT_BASE_URL")
    config.SearchURL = os.Getenv("BOT_SEARCH_URL")
    config.OfflineDownloadDir = os.Getenv("BOT_OFFLINE_DOWNLOAD_DIR")
    config.TelegramToken = os.Getenv("BOT_TELEGRAM_TOKEN")

    // 如果环境变量未设置，从配置文件加载
    if config.Username == "" {
        file, err := os.Open("config.json")
        if err != nil {
            return err
        }
        defer file.Close()

        decoder := json.NewDecoder(file)
        err = decoder.Decode(&config)
        if err != nil {
            return err
        }
    }

    return nil
}

// 全局 token 缓存
var globalToken string

// getMagnet 根据番号获取磁力链接
func getMagnet(fanhao string) (string, error) {
    url := config.SearchURL + fanhao
    log.Printf("正在搜索番号: %s...", fanhao)

    client := http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        log.Printf("获取磁力链接时出错: %v", err)
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("请求失败，状态码: %d", resp.StatusCode)
        return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("读取响应内容时出错: %v", err)
        return "", err
    }

    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        log.Printf("解析 JSON 数据时出错: %v", err)
        return "", err
    }

    dataList, ok := data["data"].([]interface{})
    if!ok || len(dataList) == 0 {
        log.Printf("错误: 未找到番号 %s 的磁力链接", fanhao)
        return "", fmt.Errorf("错误: 未找到番号 %s 的磁力链接", fanhao)
    }

    firstEntry := fmt.Sprintf("%v", dataList[0])
    parts := strings.Split(firstEntry, ",")
    magnet := strings.TrimSpace(strings.Trim(parts[0], "['"))
    log.Printf("成功获取磁力链接: %s", magnet)
    return magnet, nil
}

// getToken 获取 token
func getToken() (string, error) {
    if globalToken != "" {
        return globalToken, nil
    }

    url := config.BaseURL + "api/auth/login"
    log.Println("正在登录获取 token...")

    loginInfo := map[string]string{
        "username": config.Username,
        "password": config.Password,
    }
    payload, err := json.Marshal(loginInfo)
    if err != nil {
        log.Printf("登录信息 JSON 编码时出错: %v", err)
        return "", err
    }

    client := http.Client{Timeout: 10 * time.Second}
    resp, err := client.Post(url, "application/json", strings.NewReader(string(payload)))
    if err != nil {
        log.Printf("登录获取 token 时出错: %v", err)
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("登录失败，状态码: %d", resp.StatusCode)
        return "", fmt.Errorf("登录失败，状态码: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("读取响应内容时出错: %v", err)
        return "", err
    }

    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        log.Printf("解析 JSON 数据时出错: %v", err)
        return "", err
    }

    data, ok := result["data"].(map[string]interface{})
    if!ok || data["token"] == nil {
        message := "未知错误"
        if msg, ok := result["message"].(string); ok {
            message = msg
        }
        log.Printf("登录失败: %s", message)
        return "", fmt.Errorf("登录失败: %s", message)
    }

    token := fmt.Sprintf("%v", data["token"])
    log.Println("登录成功，已获取 token")
    globalToken = token
    return token, nil
}

// addMagnet 添加磁力链接到离线下载任务
func addMagnet(magnet string) bool {
    token, err := getToken()
    if err != nil || token == "" {
        log.Println("错误: token 为空，无法添加离线下载任务")
        return false
    }

    url := config.BaseURL + "api/fs/add_offline_download"
    log.Printf("正在添加离线下载任务到目录: %s", config.OfflineDownloadDir)

    postData := map[string]interface{}{
        "path":            config.OfflineDownloadDir,
        "urls":            []string{magnet},
        "tool":            "storage",
        "delete_policy":   "delete_on_upload_succeed",
    }
    payload, err := json.Marshal(postData)
    if err != nil {
        log.Printf("离线下载任务数据 JSON 编码时出错: %v", err)
        return false
    }

    client := http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
    if err != nil {
        log.Printf("创建请求时出错: %v", err)
        return false
    }
    req.Header.Set("Authorization", token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("添加离线下载任务时出错: %v", err)
        return false
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("添加离线下载任务失败，状态码: %d", resp.StatusCode)
        return false
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("读取响应内容时出错: %v", err)
        return false
    }

    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        log.Printf("解析 JSON 数据时出错: %v", err)
        return false
    }

    code, ok := result["code"].(float64)
    if ok && code == 200 {
        log.Println("离线下载任务添加成功!")
        return true
    }

    message := "未知错误"
    if msg, ok := result["message"].(string); ok {
        message = msg
    }
    log.Printf("添加离线下载任务失败: %s", message)
    return false
}

// listDownloadDir 列出离线下载目录内容
func listDownloadDir() error {
    token, err := getToken()
    if err != nil || token == "" {
        log.Println("错误: token 为空，无法列出目录内容")
        return err
    }

    url := config.BaseURL + "api/fs/list"
    log.Printf("正在获取目录内容: %s", config.OfflineDownloadDir)

    postData := map[string]interface{}{
        "path":     config.OfflineDownloadDir,
        "password": "",
        "page":     1,
        "per_page": 0,
        "refresh":  true,
    }
    payload, err := json.Marshal(postData)
    if err != nil {
        log.Printf("目录列表数据 JSON 编码时出错: %v", err)
        return err
    }

    client := http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
    if err != nil {
        log.Printf("创建请求时出错: %v", err)
        return err
    }
    req.Header.Set("Authorization", token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("获取目录列表时出错: %v", err)
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("获取目录列表失败，状态码: %d", resp.StatusCode)
        return fmt.Errorf("获取目录列表失败，状态码: %d", resp.StatusCode)
    }

    // 读取响应内容
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("读取目录列表响应内容时出错: %v", err)
        return err
    }

    // 解析JSON
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        log.Printf("解析目录列表JSON数据时出错: %v", err)
        return err
    }

    // 打印目录内容到日志 - 简化版，只显示文件名
    log.Println("目录内容:")
    data, ok := result["data"].(map[string]interface{})
    if ok {
        content, ok := data["content"].([]interface{})
        if ok {
            log.Printf("目录 %s 中共有 %d 个文件:", config.OfflineDownloadDir, len(content))
            for _, item := range content {
                fileInfo, ok := item.(map[string]interface{})
                if ok {
                    name, _ := fileInfo["name"].(string)
                    log.Printf("- %s", name)
                }
            }
        } else {
            log.Println("无法解析目录内容")
        }
    } else {
        log.Println("无法解析返回数据")
    }

    return nil
}

// 处理 /start 命令
func startCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, `欢迎使用thunderx离线下载机器人！
直接发送磁力链接，我会帮你添加到离线下载队列。`)
    bot.Send(msg)
}

// 处理 /help 命令
func helpCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, `使用方法：
1. 直接发送番号（如：ABC-123）
2. 直接发送磁力链接（以 magnet:? 开头）
机器人会自动处理并添加到离线下载队列。`)
    bot.Send(msg)
}

// 处理用户消息
func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    messageText := strings.TrimSpace(update.Message.Text)

    token, err := getToken()
    if err != nil || token == "" {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "错误: 获取 token 失败，无法处理请求")
        bot.Send(msg)
        return
    }

    var magnet string
    if strings.HasPrefix(messageText, "magnet:?") {
        magnet = messageText
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "收到磁力链接，正在添加到离线下载队列...")
        bot.Send(msg)
    } else {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("正在搜索番号: %s...", messageText))
        bot.Send(msg)

        var err error
        magnet, err = getMagnet(messageText)
        if err != nil {
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("错误: 未找到番号 %s 的磁力链接", messageText))
            bot.Send(msg)
            return
        }
        msg = tgbotapi.NewMessage(update.Message.Chat.ID, "已找到磁力链接，正在添加到离线下载队列...")
        bot.Send(msg)
    }

    success := addMagnet(magnet)
    if success {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ 离线下载任务添加成功！")
        bot.Send(msg)
        
        // 在成功添加后列出目录内容，并在接下来的9秒内每3秒刷新一次，共3次
        go func() {
            // 立即刷新一次
            if err := listDownloadDir(); err != nil {
                log.Printf("列出目录内容时出错: %v", err)
            }
            
            // 然后每隔3秒刷新一次，共刷新2次
            for i := 0; i < 2; i++ {
                time.Sleep(3 * time.Second)
                log.Printf("第 %d 次刷新目录...", i+2)
                if err := listDownloadDir(); err != nil {
                    log.Printf("列出目录内容时出错: %v", err)
                }
            }
        }()
    } else {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ 添加离线下载任务失败")
        bot.Send(msg)
    }
}

func main() {
    err := loadConfig()
    if err != nil {
        log.Fatalf("加载配置文件时出错: %v", err)
    }

    bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
    if err != nil {
        log.Fatalf("无法创建 Telegram 机器人: %v", err)
    }

    bot.Debug = false
    log.Printf("已授权的机器人: %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    // 添加健康检查端点，用于Render平台监控
    go func() {
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("Bot is running"))
        })
        
        // 获取端口，优先使用环境变量PORT（Render需要）
        port := os.Getenv("PORT")
        if port == "" {
            port = "8080"
        }
        
        log.Printf("启动健康检查服务器在端口 %s", port)
        if err := http.ListenAndServe(":"+port, nil); err != nil {
            log.Printf("健康检查服务器错误: %v", err)
        }
    }()

    go func() {
        for update := range updates {
            if update.Message == nil {
                continue
            }

            if update.Message.IsCommand() {
                switch update.Message.Command() {
                case "start":
                    startCommand(bot, update)
                case "help":
                    helpCommand(bot, update)
                }
            } else {
                processMessage(bot, update)
            }
        }
    }()

    log.Println("启动 Telegram 机器人...")

    // 处理中断信号
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
    log.Println("已退出!")
}

