package game_play_time_monitor

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func Test_monitorActiveWindows(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	type args struct {
		assumedGamesFolder string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestMonitorActiveWindows",
			args: args{
				assumedGamesFolder: "E:\\GalGames",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := genGamesFoldersSlice(tt.args.assumedGamesFolder)
			if err != nil {
				t.Errorf("genGamesFoldersSlice() 在执行测试：%s  时发生错误了喵, error: %v", tt.name, err)
				return
			}
			monitorActiveWindows(g)
		})
	}
}
