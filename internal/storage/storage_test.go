package storage

import (
	"fmt"
	"sync"
	"testing"
)

func TestStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		gaugeHistory   map[string][]float64
		counterHistory map[string][]map[string]int64
	}
	type args struct {
		name string
		val  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "good",
			fields: fields{
				gaugeHistory:   map[string][]float64{},
				counterHistory: map[string][]map[string]int64{},
			},
			args: args{
				name: "GoodArgs",
				val:  0.00001,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mtx:            sync.RWMutex{},
				gaugeHistory:   tt.fields.gaugeHistory,
				counterHistory: tt.fields.counterHistory,
			}
			s.UpdateGauge(tt.args.name, tt.args.val)
		})
	}
}

func TestStorage_GetGaugesByName(t *testing.T) {
	type fields struct {
		gaugeHistory   map[string][]float64
		counterHistory map[string][]map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "good",
			fields: fields{
				gaugeHistory: map[string][]float64{
					"first": {0.00002, 0.00001}},
				counterHistory: map[string][]map[string]int64{},
			},
			args:    args{name: "first"},
			want:    0.00002,
			wantErr: false,
		},
		{
			name: "empty storage",
			fields: fields{
				gaugeHistory:   map[string][]float64{},
				counterHistory: map[string][]map[string]int64{},
			},
			args:    args{name: "first"},
			want:    0,
			wantErr: true,
		},
		{
			name: "missing name",
			fields: fields{
				gaugeHistory: map[string][]float64{
					"first": {0.00002, 0.00001},
				},
				counterHistory: map[string][]map[string]int64{},
			},
			args:    args{name: "second"},
			want:    0,
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mtx:            sync.RWMutex{},
				gaugeHistory:   tt.fields.gaugeHistory,
				counterHistory: tt.fields.counterHistory,
			}
			got, err := s.GetGaugesByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetGaugesByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.GetGaugesByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		gaugeHistory   map[string][]float64
		counterHistory map[string][]map[string]int64
	}
	type args struct {
		name string
		val  int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "good",
			fields: fields{
				gaugeHistory: map[string][]float64{},
				counterHistory: map[string][]map[string]int64{
					"first": {
						{
							"counter": 3,
							"value":   2,
						},
						{
							"counter": 2,
							"value":   1,
						},
						{
							"counter": 1,
							"value":   1,
						},
					}},
			},
			args: args{
				name: "first",
				val:  2,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mtx:            sync.RWMutex{},
				gaugeHistory:   tt.fields.gaugeHistory,
				counterHistory: tt.fields.counterHistory,
			}
			s.UpdateCounter(tt.args.name, tt.args.val)
			fmt.Println(s)
		})
	}
}

func TestStorage_GetCountersByName(t *testing.T) {
	type fields struct {
		gaugeHistory   map[string][]float64
		counterHistory map[string][]map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "good",
			fields: fields{
				gaugeHistory: map[string][]float64{},
				counterHistory: map[string][]map[string]int64{
					"first": {
						{
							"counter": 3,
							"value":   2,
						},
						{
							"counter": 2,
							"value":   1,
						},
						{
							"counter": 1,
							"value":   1,
						},
					},
				},
			},
			args: args{
				name: "first",
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "empty storage",
			fields: fields{
				gaugeHistory:   map[string][]float64{},
				counterHistory: map[string][]map[string]int64{},
			},
			args: args{
				name: "first",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "missing name",
			fields: fields{
				gaugeHistory: map[string][]float64{},
				counterHistory: map[string][]map[string]int64{
					"first": {
						{
							"counter": 3,
							"value":   2,
						},
						{
							"counter": 2,
							"value":   1,
						},
						{
							"counter": 1,
							"value":   1,
						},
					},
				},
			},
			args: args{
				name: "second",
			},
			want:    0,
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mtx:            sync.RWMutex{},
				gaugeHistory:   tt.fields.gaugeHistory,
				counterHistory: tt.fields.counterHistory,
			}
			got, err := s.GetCountersByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetCountersByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.GetCountersByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
