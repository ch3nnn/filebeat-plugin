package actions

import (
	"fmt"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/processors"
	"github.com/elastic/beats/v7/libbeat/processors/checks"
	"strings"
)

func init() {
	processors.RegisterPlugin("parse_text",
		checks.ConfigChecked(NewParseText,
			checks.RequireFields("file_has_suffix")),
	)
}

type parseText struct {
	*defaultParse
}

func NewParseText(c *common.Config) (processors.Processor, error) {
	// 读取配置字段
	config := new(config)
	err := c.Unpack(config)
	if err != nil {
		return nil, fmt.Errorf("fail to unpack the add_fields configuration: %s", err)
	}

	p := parseText{newDefaultParse(*config)}
	return p, nil
}

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

func (p parseText) String() string {
	return fmt.Sprintf("parse_text=%v", p.result)

}
