# filebeat-plugin


> Filebeat轻量型日志采集器-自定义processors处理器


## 引入对beat的依赖
`go get github.com/elastic/beats/v7`

## 定义在filebeat中的配置文件
filebeat通常以配置文件的方式加载插件。让定义一下自定义配置。
```yaml
filebeat.inputs:
  - type: log
    paths:
      - example/example.log


processors:
  # 自定义处理器插件  
  - parse_text:
      file_has_suffix: example.log


output.console:
  pretty: true
```

## go文件中的配置
```go
package actions

type config struct {
	FileHasSuffix string `config:"file_has_suffix" validate:"required"`
}
```

## 初始化加载插件
```go
func init() {
	processors.RegisterPlugin("parse_text",
		checks.ConfigChecked(NewParseText,
			checks.RequireFields("file_has_suffix")),
	)
}
```
## Run接口
处理 **filebeat** 读取到的每行日志数据`message`, 这里我们就可以自定义一些处理解析逻辑, 下面逻辑是将日志数据按`,`切分,重新组装到`event.Fields`字段里.
```go
func (p parseText) Run(event *beat.Event) (*beat.Event, error) {
	if !p.isParseFile(event) {
		return event, nil
	}
	message, err := p.getMessage(event)
	if err != nil {
		return event, nil
	}
	// 按空格切分文本
	split := strings.Split(message, ",")
	p.result["split"] = split
	_, err = event.Fields.Put("parse_text", p.result)
	if err != nil {
		return nil, err
	}

	return event, nil

}
```
## main函数
```go
package main

import (
	"os"

	_ "filebeat-plugin/pkg/processors/actions"  // 这里需要将自定义插件注册
	"github.com/elastic/beats/v7/filebeat/cmd"
	inputs "github.com/elastic/beats/v7/filebeat/input/default-inputs"
)

func main() {
	if err := cmd.Filebeat(inputs.Init, cmd.FilebeatSettings()).Execute(); err != nil {
		os.Exit(1)
	}
}
```

## 打包 build
项目路径下执行 `go build -o filebeat.exe  ./cmd/filebeat/`

## 解析结果及源日志
源日志如下:
```text
I,am,Bat,I
```
解析后 `parse_text` 字段数据如下:
```json
"parse_text": {
    "split": [
      "I",
      "am",
      "Bat",
      "I"
    ]
  },
```
完整解析结果如下:
```text
{
  "@timestamp": "2023-02-23T10:26:05.154Z",
  "@metadata": {
    "beat": "filebeat",
    "type": "_doc",
    "version": "7.13.4"
  },
  "host": {
    "name": "chentong.local"
  },
  "agent": {
    "version": "7.13.4",
    "hostname": "chentong.local",
    "ephemeral_id": "02d53454-19dc-4fff-9182-8f04a402e399",
    "id": "f10cf835-535b-4714-8bc4-1c0b3636ae1d",
    "name": "chentong.local",
    "type": "filebeat"
  },
  "ecs": {
    "version": "1.8.0"
  },
  "parse_text": {
    "split": [
      "I",
      "am",
      "Bat",
      "I"
    ]
  },
  "log": {
    "offset": 187,
    "file": {
      "path": "filebeat-plugin/example/example.log"
    }
  },
  "message": "I,am,Bat,I",
  "input": {
    "type": "log"
  }
}

```



## 源码

[https://github.com/ch3nnn/filebeat-plugin](https://github.com/ch3nnn/filebeat-plugin)

