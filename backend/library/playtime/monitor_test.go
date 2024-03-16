package playtime

import "testing"
import "YoshinoGal/backend/library"

func Test_monitor(t *testing.T) {
	type args struct {
		gameBaseFolder string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestMonitor",
			args: args{
				gameBaseFolder: "E:\\GalGames",
			},
		},
	}
	for _, tt := range tests {
		//Monitor(tt.args.gameBaseFolder, tt.args.gamePlayTimeFilePath)
		t.Run(tt.name, func(t *testing.T) {
			db, err := library.InitGameLibrary(tt.args.gameBaseFolder)
			if err != nil {
				t.Errorf("InitGameLibrary() error = %v", err)
				return
			}
			StartMonitor(db)
			for {
				continue
			}
		})
	}
}
