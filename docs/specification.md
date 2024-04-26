# Golang 规范

## 前后端事件命名

1. 项目初始化期间（在`Startup`等钩子函数中发起的事件），事件根据具体含义命名
2. 项目运行期间发生的错误统一事件名称为`GlobalRuntimeError`，并在data中包含`errorMessage`、`errorName`两部分