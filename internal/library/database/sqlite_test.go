package database

import (
	"YoshinoGal/internal/library/types"
	"database/sql"
	"reflect"
	"testing"
)

func TestInitSQLiteDB(t *testing.T) {
	type args struct {
		dbPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestInitSQLiteDB",
			args: args{
				dbPath: "test.db",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InitSQLiteDB(tt.args.dbPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitSQLiteDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewSqliteGameLibrary(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *SqliteGameLibrary
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSqliteGameLibrary(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSqliteGameLibrary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqliteGameLibrary_GetGameDataByName(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.GalgameMetadata
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SqliteGameLibrary{
				db: tt.fields.db,
			}
			got, err := s.GetGameDataByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGameDataByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGameDataByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqliteGameLibrary_InsertGameMetadata(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		game types.GalgameMetadata
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SqliteGameLibrary{
				db: tt.fields.db,
			}
			if err := s.InsertGameMetadata(tt.args.game); (err != nil) != tt.wantErr {
				t.Errorf("InsertGameMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
