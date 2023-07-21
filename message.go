package main

import (
	_ "embed"
	"fmt"
	debug "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

type MessageKeyValue struct {
	Key   string
	Value interface{}
}

const (
	JA_JP = "ja-jp"
)

var messages []MessageKeyValue

var languages = []language.Tag{
	language.Japanese,
	language.English,
}

func init() {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(messageJp), &m)

	if err != nil {
		logger.WithFields(debug.Fields{
			"event": "init",
		}).Error(err)
	}
	processMap(m, "", &messages)
}

func GetMessage(acceptLanguage string, msg string, args ...interface{}) string {
	t, _, _ := language.ParseAcceptLanguage(acceptLanguage)
	matcher := language.NewMatcher(languages)
	tag, _, _ := matcher.Match(t...)
	p := message.NewPrinter(tag)
	return p.Sprintf(msg, args...)
}

func processMap(m map[interface{}]interface{}, parentKey string, result *[]MessageKeyValue) {
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
			message.SetString(language.Japanese, key, valueStr)
			*result = append(*result, MessageKeyValue{Key: key, Value: valueStr})
		}
	}

	logger.WithFields(debug.Fields{
		"event":  "message",
		"length": len(*result),
	}).Debug("complete.")

}
