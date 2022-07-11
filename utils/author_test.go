package utils

import (
	"reflect"
	"sync"
	"testing"
)

func Test_authors_Add(t *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		email string
		name  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty", fields{data: map[string]string{}}, args{}, true},
		{"add_no_email", fields{data: map[string]string{}}, args{"", "test"}, true},
		{"add_no_name", fields{data: map[string]string{}}, args{"test", ""}, true},
		{"success", fields{data: map[string]string{}}, args{"test", "test"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authors{
				data:  tt.fields.data,
				mutex: sync.RWMutex{},
			}
			if err := a.Add(tt.args.email, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("authors.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authors_Find(t *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantEmail *string
		wantErr   bool
	}{
		{"empty", fields{data: map[string]string{}}, args{""}, nil, true},
		{"not_found", fields{data: map[string]string{}}, args{"test@test.com"}, nil, true},
		{"success", fields{data: map[string]string{"test@test.com": "test"}}, args{"test@test.com"}, String("test"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authors{
				data:  tt.fields.data,
				mutex: sync.RWMutex{},
			}
			gotEmail, err := a.Find(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("authors.Check() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if (gotEmail == nil || tt.wantEmail == nil) && gotEmail != tt.wantEmail {
				t.Errorf("authors.Check() = %v, want %v", *gotEmail, *tt.wantEmail)
			}
			if (gotEmail != nil && tt.wantEmail != nil) && *gotEmail != *tt.wantEmail {
				t.Errorf("authors.Check() = %v, want %v", *gotEmail, *tt.wantEmail)
			}
		})
	}
}

func Test_authors_Details(t *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{"empty", fields{data: map[string]string{}}, args{""}, nil, true},
		{"data", fields{data: map[string]string{"test@email.com": "test"}}, args{"test"}, []string{"test@email.com"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authors{
				data:  tt.fields.data,
				mutex: sync.RWMutex{},
			}
			got, err := a.Details(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("authors.Details() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authors.Details() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authors_Check(t *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		email string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"empty", fields{data: map[string]string{}}, args{""}, false},
		{"not_found", fields{data: map[string]string{}}, args{"test@test.com"}, false},
		{"success", fields{data: map[string]string{"test@test.com": "test"}}, args{"test@test.com"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authors{
				data:  tt.fields.data,
				mutex: sync.RWMutex{},
			}
			if got := a.Check(tt.args.email); got != tt.want {
				t.Errorf("authors.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
