package sources

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestVNDBSearch(t *testing.T) {
	type args struct {
		gameName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Search Senren * Banka",
			args: args{
				gameName: "Senren * Banka",
			},
			wantErr: false,
		},
		{
			name: "Search ATRI",
			args: args{
				gameName: "ATRI -My Dear Moments-",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.TraceLevel)
		t.Run(tt.name, func(t *testing.T) {
			_, err := VNDBSearch(tt.args.gameName, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("VNDB搜索失败，失败原因： %v, wantErr %v", err, tt.wantErr)
				return
			}
		})

	}
}

func TestVNDBSearchWithEmptyName(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
	type args struct {
		gameName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "EmptyName",
			args: args{
				gameName: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VNDBSearch(tt.args.gameName, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("VNDBSearch 没有对空字符串进行处理，笨蛋！ error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestVNDBSearchWithInvalidName(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
	type args struct {
		gameName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "InvalidName",
			args: args{
				gameName: "我不知道写啥喵~",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VNDBSearch(tt.args.gameName, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("没有对无搜索结果的游戏进行处理捏！ error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})

	}
}
