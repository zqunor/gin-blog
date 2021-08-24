教程来源：https://segmentfault.com/a/1190000013297683

1、目录说明

```
gin-blog/
├── conf
├── middleware
├── models
├── pkg
├── routers
└── runtime
```

- conf：用于存储配置文件
- middleware：应用中间件
- models：应用数据库模型
- pkg：第三方包
- router：路由逻辑处理
- runtime：应用运行时数据


## 六、文件操作

1、涉及的库

- os

相关开发文档：[Golang os](https://cloud.tencent.com/developer/section/1143829)

2、方法

- `os.Getwd()` 获取当前目录的根路径名
- `os.MkdirAll(fullPath, os.ModePerm)`: 创建对应的目录以及所需的子目录，若成功则返回nil，否则返回error
  - os.ModePerm：const定义ModePerm FileMode = 0777
- `os.Stat(filePath)` 获取文件的FileInfo
- `os.OpenFile(filePath, mode, perm)`  以指定的模式和权限打开文件

