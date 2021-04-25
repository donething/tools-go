package textcoding

import (
	"github.com/saintfish/chardet"
	"reflect"
	"testing"
)

func TestDetectFileCoding(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *chardet.Result
		wantErr bool
	}{
		{
			"UTF-8",
			args{"D:/影视/连续剧/《金牌冰人》国语外挂GOTV_720P_TS_800M/01.srt"},
			&chardet.Result{"UTF-8", "", 100},
			false,
		},
		{
			"GB-18030",
			args{"D:/影视/连续剧/《金牌冰人》国语外挂GOTV_720P_TS_800M/20.srt"},
			&chardet.Result{"GB-18030", "zh", 100},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectFileCoding(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFileCoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectFileCoding() got = %v, want %v", got, tt.want)
			}
		})
	}
}
