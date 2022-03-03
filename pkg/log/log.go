package log

import (
	"fmt"
	"log"
)

type ShortlyLog struct {
	Error   error
	Message string
	Tag     string
	format  string
}

func (s *ShortlyLog) ZapError(message string, err error) {
	s.Message = message
	s.Error = err
	s.format = fmt.Sprintf("[ERROR] %s, %v", s.Message, s.Error)
	log.Println(s.format)
}

func (s *ShortlyLog) Zap(message string) {
	s.Message = message

	if s.Tag == "" {
		s.Tag = "INFO"
	}

	s.format = fmt.Sprintf("[%s] %s", s.Tag, s.Message)
	log.Println(s.format)
}

func (s *ShortlyLog) ZapFatal(message string, err error) {
	s.Message = message

	s.format = fmt.Sprintf("[FATAL] %s", s.Message)
	log.Fatalln(s.format, err)
}
