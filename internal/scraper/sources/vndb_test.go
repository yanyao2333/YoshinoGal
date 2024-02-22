package sources

import (
	"testing"
)

func TestVNDBSearch(t *testing.T) {
	//logrus.SetReportCaller(true)
	//logrus.SetLevel(logrus.TraceLevel)

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
		{
			name: "EmptyName",
			args: args{
				gameName: "",
			},
			wantErr: true,
		},
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
			_, err := SearchInVNDB(tt.args.gameName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchInVNDB 在执行测试：%s  时发生错误了喵, error: %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
