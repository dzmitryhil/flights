package finder

import (
	"reflect"
	"testing"
)

func TestFlightPath_Find(t *testing.T) {
	type args struct {
		flights [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "positive_one_point",
			args: args{
				flights: [][]string{
					{"SFO", "EWR"},
				},
			},
			want:    []string{"SFO", "EWR"},
			wantErr: false,
		},
		{
			name: "positive_two_unordered_point",
			args: args{
				flights: [][]string{
					{"SFO", "ATL"},
					{"ATL", "EWR"},
				},
			},
			want:    []string{"SFO", "EWR"},
			wantErr: false,
		},
		{
			name: "positive_four_unordered_point",
			args: args{
				flights: [][]string{
					{"IND", "EWR"},
					{"SFO", "ATL"},
					{"GSO", "IND"},
					{"ATL", "GSO"},
				},
			},
			want:    []string{"SFO", "EWR"},
			wantErr: false,
		},
		{
			name: "positive_valid_sub_cycle_from_departure",
			args: args{
				flights: [][]string{
					{"A", "B"},
					{"B", "C"},
					{"C", "B"},
				},
			},
			want:    []string{"A", "B"},
			wantErr: false,
		},
		{
			name: "positive_valid_sub_cycle_from_arrival",
			args: args{
				flights: [][]string{
					{"A", "B"},
					{"B", "A"},
					{"A", "B"},
					{"B", "A"},
					{"A", "C"},
				},
			},
			want:    []string{"A", "C"},
			wantErr: false,
		},
		{
			name: "positive_valid_two_sub_cycles",
			args: args{
				flights: [][]string{
					{"A", "B"},
					{"B", "C"},
					{"C", "A"},
					{"A", "B"},
					{"B", "C"},
					{"C", "A"},
					{"A", "B"},
					{"B", "C"},
					{"C", "A"},
					{"A", "Z"},
				},
			},
			want:    []string{"A", "Z"},
			wantErr: false,
		},
		{
			name: "negative_duplicated_departure",
			args: args{
				flights: [][]string{
					{"A", "B"},
					{"B", "C"},
					{"A", "B"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative_two_arrivals",
			args: args{
				flights: [][]string{
					{"A", "B"},
					{"C", "B"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative_two_departures",
			args: args{
				flights: [][]string{
					{"A", "B"},
					{"B", "C"},
					{"B", "D"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative_cycle",
			args: args{
				flights: [][]string{
					{"A", "B"},
					{"B", "C"},
					{"C", "A"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := &FlightPath{}
			got, err := f.Find(tt.args.flights)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}
