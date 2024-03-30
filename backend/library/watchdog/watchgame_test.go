package watchdog

import (
	"YoshinoGal/backend/library/database"
	"testing"
)

func TestWatchGame(t *testing.T) {
	type args struct {
		gameDir         string
		lib             *database.SqliteGameLibrary
		scraperPriority []string
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = WatchGame(tt.args.gameDir, tt.args.lib, tt.args.scraperPriority)
		})
	}
}
