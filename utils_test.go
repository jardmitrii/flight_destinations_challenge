package main

import "testing"

func Test_getDestination(t *testing.T) {
	tests := []struct {
		given  string
		dest   string
		want   string
		wantOk bool
	}{
		{given: "LAX", dest: "LAXVIE", want: "VIE", wantOk: true},
		{given: "LAX", dest: "laxfra", want: "FRA", wantOk: true},
		{given: "LAX", dest: "NYCLAX", want: "NYC", wantOk: true},
		{given: "LAX", dest: "BOLAXS", want: "", wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.given, func(t *testing.T) {
			if got, ok := getDestination(tt.given, tt.dest); got != tt.want || ok != tt.wantOk {
				t.Errorf("checkDest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLetters(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{name: "All_lower", args: "abcdefghijklmnopqrstuvwxyz", want: true},
		{name: "All_upper", args: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", want: true},
		{name: "An_empty", args: "", want: true},
		{name: "A_number", args: "1", want: false},
		{name: "A_special", args: "_", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLetters(tt.args); got != tt.want {
				t.Errorf("isLetters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateAirportCode(t *testing.T) {
	tests := []struct {
		name    string
		origin  string
		wantErr bool
	}{
		{name: "Valid", origin: "LAX", wantErr: false},
		{name: "Too_long", origin: "LAXXX", wantErr: true},
		{name: "Too_short", origin: "LA", wantErr: true},
		{name: "Non_latin", origin: "LA1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAirportCode(tt.origin); (err != nil) != tt.wantErr {
				t.Errorf("validateAirportCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
