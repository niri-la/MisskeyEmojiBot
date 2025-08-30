package handler

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v2"

	"MisskeyEmojiBot/pkg/entity"
)

const (
	JA_JP = "ja-jp"
)

var messages []entity.MessageKeyValue

var languages = []language.Tag{
	language.Japanese,
	language.English,
}

var messageJp string

type MessageHandler interface {
	GetMessage(acceptLanguage string, msg string, args ...interface{}) string
}

type messageHanmdler struct {
}

func NewMessageHandler() (MessageHandler, error) {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(messageJp), &m)
	if err != nil {
		return nil, err
	}
	processMap(m, "", &messages)
	return &messageHanmdler{}, nil
}

func processMap(m map[interface{}]interface{}, parentKey string, result *[]entity.MessageKeyValue) {
	for k, v := range m {
		key := fmt.Sprintf("%v", k)
		if parentKey != "" {
			key = parentKey + "." + key
		}
		switch value := v.(type) {
		case map[string]interface{}:
			newMap := make(map[interface{}]interface{})
			for k, v := range value {
				newMap[k] = v
			}
			processMap(newMap, key, result)
		case map[interface{}]interface{}:
			processMap(value, key, result)
		default:
			valueStr := fmt.Sprintf("%v", v)
			message.SetString(language.Japanese, key, valueStr) //nolint:errcheck
			*result = append(*result, entity.MessageKeyValue{Key: key, Value: valueStr})
		}
	}
}

func (h *messageHanmdler) GetMessage(acceptLanguage string, msg string, args ...interface{}) string {
	t, _, _ := language.ParseAcceptLanguage(acceptLanguage)
	matcher := language.NewMatcher(languages)
	tag, _, _ := matcher.Match(t...)
	p := message.NewPrinter(tag)
	return p.Sprintf(msg, args...)
}
