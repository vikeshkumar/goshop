package db

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	type args struct {
		context context.Context
		uri     string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		timesConnect int
	}{
		{"should be able to connect",
			args{context.Background(), "postgres://goshop:goshop@localhost:5432/goshop"},
			false,
			1,
		},
		{"should not able to connect",
			args{context.Background(), "postgres://goshop:goshop@localhost:54321/goshop"},
			true,
			1,
		},
		{"multiple connect should fails",
			args{context.Background(), "postgres://goshop:goshop@localhost:54321/goshop"},
			true,
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var connectError error
			var p *pgxpool.Pool
			for i := 0; i < tt.timesConnect; i++ {
				pool, err := Connect(tt.args.context, tt.args.uri)
				connectError = err
				p = pool
			}
			if tt.wantErr && connectError == nil && tt.timesConnect == 1 {
				t.Errorf("connected to database")
			}
			if !tt.wantErr && connectError != nil && tt.timesConnect == 1 {
				t.Errorf("not able to connect to database")
			}
			if tt.timesConnect > 1 {
				if p != nil && connectError != nil {
					t.Errorf("was able to connect")
				}
			}
			if p != nil {
				defer p.Close()
			}
		})
	}
}

func TestQuery(t *testing.T) {
	type args struct {
		ctx  context.Context
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    pgx.Rows
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.args.ctx, tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
}
