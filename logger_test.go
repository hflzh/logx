package logx

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	testCases := []struct {
		out          io.Writer
		level        Level
		useLocalTime bool
		wantNil      bool // If the New() function should return nil.
	}{
		{nil, Fine, true, true},
		{&buf, Off, true, true},
		{&buf, Fine - 1, true, true},
		{&buf, Error + 1, true, true},
		{&buf, Fine, true, false},
		{&buf, Fine, false, false},
		{&buf, Debug, true, false},
		{&buf, Debug, false, false},
		{&buf, Info, true, false},
		{&buf, Info, false, false},
		{&buf, Warn, true, false},
		{&buf, Warn, false, false},
		{&buf, Error, true, false},
		{&buf, Error, false, false},
	}

	var lgr *Logger
	for _, ts := range testCases {
		lgr = New(ts.out, ts.level, ts.useLocalTime)
		if ts.wantNil {
			if lgr != nil {
				t.Errorf("New(out=%v, level=%v, useLocalTime=%v) should have failed",
					ts.out, ts.level, ts.useLocalTime)
			}
			continue
		}

		// wantNil == false
		if lgr == nil {
			t.Errorf("New(out=%v, level=%v, useLocalTime=%v) should have succeeded",
				ts.out, ts.level, ts.useLocalTime)
			continue
		}

		if lgr.level != ts.level {
			t.Errorf("unexpected logging level, got %v, want %v", lgr.level, ts.level)
		}

		wantTimeSuffix := ""
		if !ts.useLocalTime {
			wantTimeSuffix = "UTC "
		}

		if lgr.timeSuffix != wantTimeSuffix {
			t.Errorf("unexpected timezone string, got %s, want %s", lgr.timeSuffix, wantTimeSuffix)
		}
	}
}

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	const msg string = "Log a message at a specific level."
	allLevels := [...]Level{Fine, Debug, Info, Warn, Error, Off}
	for _, level := range allLevels {
		lgr := New(&buf, level, false)
		for _, msgLevel := range allLevels {
			buf.Reset()

			switch msgLevel {
			case Fine:
				lgr.Fine(msg)
			case Debug:
				lgr.Debug(msg)
			case Info:
				lgr.Info(msg)
			case Warn:
				lgr.Warn(msg)
			case Error:
				lgr.Error(msg)
			case Off:
				lgr.Log(Off, msg)
			}

			msgPrefix := fmt.Sprintf("LoggerLevel=%s, MessageLevel=%s: ",
				level, msgLevel)
			logEntry := buf.String()
			if level == Off || msgLevel == Off || msgLevel < level {
				if len(logEntry) > 0 {
					t.Errorf(msgPrefix + "log message should have been empty")
				}
				continue
			}

			if !strings.Contains(logEntry, "UTC "+label(msgLevel)) || !strings.Contains(logEntry, msg) {
				t.Errorf(msgPrefix + "wrong log message")
			}
		}
	}
}
