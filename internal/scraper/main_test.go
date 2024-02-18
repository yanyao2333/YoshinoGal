package scraper

import (
	"YoshinoGal/internal/scraper/types"
	"testing"
)

func TestScanGamesAndScrape(t *testing.T) {
	//logrus.SetReportCaller(true)
	//logrus.SetLevel(logrus.TraceLevel)

	type args struct {
		directory string
		priority  []types.GalgameMetadataSource
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "EmptyDirectory",
			args: args{
				directory: "",
			},
			wantErr: true,
		},
		{
			name: "EmptyPriority",
			args: args{
				directory: "test",
			},
			wantErr: true,
		},
		{
			name: "RealTest",
			args: args{
				directory: "E:\\GalGames",
				priority:  []types.GalgameMetadataSource{"VNDB"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ScanGamesAndScrape(tt.args.directory, tt.args.priority); (err != nil) != tt.wantErr {
				t.Errorf("ScanGamesAndScrape 在执行测试：%s  时发生错误了喵, error: %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
