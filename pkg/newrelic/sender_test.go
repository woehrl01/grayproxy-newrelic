package newrelic

import (
	"reflect"
	"testing"
)

func Test_convertToNewRelicFormat(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name       string
		args       args
		wantOutput []byte
		wantErr    bool
	}{
		{name: "convert graylog to newrelic", args: struct{
			data []byte}{ data: []byte(`{"version":"1.1","host":"example.org","short_message":"A short message","full_message":"Backtrace here nnn","timestamp":1296000000.0,"level":5,"_some_info":"foo","_some_env_var":"bar"}`)}, 
			wantOutput: []byte(`{"_some_env_var":"bar","_some_info":"foo","logtype":"NOTICE","message":"A short message Backtrace here nnn","timestamp":1296000000}`), 
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, err := convertToNewRelicFormat(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToNewRelicFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotOutputAsString := string(gotOutput)
			wantOutputAsString := string(tt.wantOutput)

			if !reflect.DeepEqual(gotOutputAsString, wantOutputAsString) {
				t.Errorf("convertToNewRelicFormat() = %v, want %v", gotOutputAsString, wantOutputAsString)
			}
		})
	}
}
