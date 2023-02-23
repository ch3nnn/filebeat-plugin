package actions

import (
	"errors"
	"fmt"
	"github.com/elastic/beats/v7/filebeat/input/file"
	"github.com/elastic/beats/v7/libbeat/beat"
	"strings"
)

var (
	errFieldIsNotString = errors.New("field value is not a string")
)

type Parse interface {
	// 判断插件是否是当前解析日志文件
	isParseFile(event *beat.Event) bool
	// 获取文件日志文本
	getMessage(event *beat.Event) (string, error)
}

type defaultParse struct {
	result map[string]interface{}
	config config
}

func newDefaultParse(config config) *defaultParse {
	return &defaultParse{
		result: map[string]interface{}{},
		config: config,
	}
}

func (d defaultParse) isParseFile(event *beat.Event) bool {
	private, ok := event.Private.(file.State)
	if !ok {
		return false
	}
	if !strings.HasSuffix(private.Source, d.config.FileHasSuffix) {
		return false
	}
	return true
}

func (d defaultParse) getMessage(event *beat.Event) (string, error) {
	message, _ := event.Fields.GetValue("message")
	str, ok := message.(string)
	if !ok {
		err := fmt.Errorf("failed in message on the %q field: %w", "message", errFieldIsNotString)
		return "", err
	}
	return str, nil
}
