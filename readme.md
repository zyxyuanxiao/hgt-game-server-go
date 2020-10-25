# 服务端框架

> 使用 golang 高性能语言实现

---

### 目录结构

```
├─api                      # 接口逻辑
|   ├─example.go           # 示例
|   ├─...                  # 更多
├─app                      # 框架核心
|   ├─app.go               # 框架初始化
|   ├─cron.go              # 自动任务管理
|   ├─db.go                # 数据库加载
|   ├─logger.go            # 日志处理
|   ├─middleware.go        # 中间件管理
├─config                   # 配置
|   ├─develop              # 开发环境配置
|   |   ├─db.json          # 数据库配置
|   ├─prod                 # 线上配置
├─console                  # 自动任务管理
├─exception                # 异常管理
|   ├─LogicException.go    # 逻辑异常
|   ├─...                  # 更多异常
├─logs                     # 日志管理
├─model                    # 模型层
|   ├─mysql                # 数据库模型目录
├─router                   # 路由管理
|   ├─config.go            # 路由初始化配置
|   ├─User.go              # 用户模块路由初始化
|   ├─...                  # 更多模块
├─servers                  # 存放服务
├─main.go                  # 入口文件
├─.gitignore               # git提交忽略配置
```

----
