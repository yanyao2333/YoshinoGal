package sources

import (
	"encoding/json"
	"os"
	"strconv"
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
	var num = 0

	for _, tt := range tests {
		num++
		t.Run(tt.name, func(t *testing.T) {
			game, err := SearchInVNDB(tt.args.gameName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchInVNDB 在执行测试：%s  时发生错误了喵, error: %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			var data, _ = json.MarshalIndent(game, "", "    ")
			var fileName = "vndb_test_" + strconv.Itoa(num) + ".json"
			os.WriteFile(fileName, data, 0777)
		})
	}
}
