package main

import (
	"fmt"
	"log"

	"github.com/mailgun/groupcache/v2"
)

type field struct {
	key   string
	value any
}

type logger struct {
	fields []field
}

func (l *logger) Error() groupcache.Logger {
	return l
}

func (l *logger) Warn() groupcache.Logger {
	return l
}

func (l *logger) Info() groupcache.Logger {
	return l
}

func (l *logger) Debug() groupcache.Logger {
	return l
}

func (l *logger) ErrorField(label string, err error) groupcache.Logger {
	return &logger{
		fields: append(l.fields, field{label, err}),
	}
}

func (l *logger) StringField(label string, val string) groupcache.Logger {
	return &logger{
		fields: append(l.fields, field{label, val}),
	}
}

func (l *logger) WithFields(fields map[string]interface{}) groupcache.Logger {
	res := &logger{
		fields: l.fields,
	}

	for k, v := range fields {
		res.fields = append(res.fields, field{k, v})
	}

	return res
}

func (l *logger) Printf(format string, args ...interface{}) {
	if len(l.fields) > 0 {
		format += " {"
		for _, v := range l.fields {
			format = format + v.key + "=" + fmt.Sprint(v.value) + ","
		}
		format += "}"
	}

	log.Printf(format, args...)
}

func newLogger() *logger {
	return &logger{}
}
