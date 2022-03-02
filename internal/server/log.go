package server

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

func (s *ShortlyLog) ZapError(message string, error error) {
	s.Message = message
	s.Error = error
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
