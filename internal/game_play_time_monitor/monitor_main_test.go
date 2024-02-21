package game_play_time_monitor

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func Test_monitorActiveWindows(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	type args struct {
		gameBaseFolder       string
		gamePlayTimeFilePath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestMonitorActiveWindows",
			args: args{
				gameBaseFolder:       "E:\\GalGames",
				gamePlayTimeFilePath: "E:\\GalGames\\playTime.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gamePlayTimeMonitor(tt.args.gameBaseFolder, tt.args.gamePlayTimeFilePath)
		})
	}
}
