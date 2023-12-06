package kazuyomi

import (
	"math"
	"testing"
)

func TestReadString(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "0", want: "ゼロ"},
		{in: "1", want: "イチ"},
		{in: "2", want: "ニ"},
		{in: "3", want: "サン"},
		{in: "4", want: "ヨン"},
		{in: "5", want: "ゴ"},
		{in: "6", want: "ロク"},
		{in: "7", want: "ナナ"},
		{in: "8", want: "ハチ"},
		{in: "9", want: "キュウ"},
		{in: "10", want: "ジュウ"},
		{in: "11", want: "ジュウイチ"},
		{in: "20", want: "ニジュウ"},
		{in: "42", want: "ヨンジュウニ"},
		{in: "100", want: "ヒャク"},
		{in: "101", want: "ヒャクイチ"},
		{in: "111", want: "ヒャクジュウイチ"},
		{in: "3300", want: "サンゼンサンビャク"},
		{in: "666", want: "ロッピャクロクジュウロク"},
		{in: "8888", want: "ハッセンハッピャクハチジュウハチ"},
		{in: "1_0000", want: "イチマン"},
		{in: "1234_5678", want: "センニヒャクサンジュウヨンマンゴセンロッピャクナナジュウハチ"},
		{in: "1_0000_0000", want: "イチオク"},
		{in: "1_1000_1000", want: "イチオクセンマンセン"},
		{in: "1_2345_6789", want: "イチオクニセンサンビャクヨンジュウゴマンロクセンナナヒャクハチジュウキュウ"},
		{in: "9876_0000_4321", want: "キュウセンハッピャクナナジュウロクオクヨンセンサンビャクニジュウイチ"},
		{in: "1_0000_0000_0000", want: "イッチョウ"},
		{in: "8_0000_0000_0000", want: "ハッチョウ"},
		{in: "10_0000_0000_0000", want: "ジッチョウ"},
		{in: "1111_0000_0000_0000", want: "センヒャクジュウイッチョウ"},
		{in: "8888_0000_0000_0000", want: "ハッセンハッピャクハチジュウハッチョウ"},
		{in: "1_0000_0000_0000_0000", want: "イッケイ"},
		{in: "8_0000_0000_0000_0000", want: "ハッケイ"},
		{in: "10_0000_0000_0000_0000", want: "ジッケイ"},
		{in: "12,3456,7890", want: "ジュウニオクサンゼンヨンヒャクゴジュウロクマンナナセンハッピャクキュウジュウ"},
		{in: "1,234,567,890", want: "ジュウニオクサンゼンヨンヒャクゴジュウロクマンナナセンハッピャクキュウジュウ"},
		{in: "0.10", want: "レイテンイチゼロ"},
		{in: "1.23", want: "イッテンニサン"},
		{in: "10.1", want: "ジッテンイチ"},
		{in: "18.88", want: "ジュウハッテンハチハチ"},
		{in: "42.195", want: "ヨンジュウニテンイチキュウゴ"},
		{in: ".1", want: "テンイチ"},
		{in: "0.", want: "ゼロ"},
		{in: "1.", want: "イチ"},
		{in: "+1", want: "プラスイチ"},
		{in: "-1", want: "マイナスイチ"},
		{in: "127.0.0.1", want: "イチニナナテンゼロテンゼロテンイチ"},
		{in: "0120", want: "ゼロイチニゼロ"},
		{in: "1_2345_6789_0123_4567_8901", want: "イチニサンヨンゴロクナナハチキュウゼロイチニサンヨンゴロクナナハチキュウゼロイチ"},
	}

	for _, tt := range tests {
		got, err := ReadString(tt.in)
		if err != nil {
			t.Errorf("ReadString(%s) returns error unexpectedly: %s", tt.in, err)
		}
		if got != tt.want {
			t.Errorf("ReadString(%s) = %s, want %s", tt.in, got, tt.want)
		}
	}
}

func TestReadString_error(t *testing.T) {
	tests := []string{
		"foobar",
		"1+2",
		"*1",
	}

	for _, tt := range tests {
		_, err := ReadString(tt)
		if err == nil {
			t.Errorf("ReadString(%s) should return error", tt)
		}
	}
}

func TestReadInt(t *testing.T) {
	tests := []struct {
		in   int
		want string
	}{
		{in: 0, want: "ゼロ"},
		{in: 1, want: "イチ"},
		{in: -1, want: "マイナスイチ"},
		// 9,223,372,036,854,775,807
		{in: math.MaxInt64, want: "キュウヒャクニジュウニケイサンゼンサンビャクナナジュウニチョウサンビャクロクジュウハチオクゴセンヨンヒャクナナジュウナナマンゴセンハッピャクナナ"},
	}

	for _, tt := range tests {
		got := ReadInt(tt.in)
		if got != tt.want {
			t.Errorf("ReadInt(%d) = %s, want %s", tt.in, got, tt.want)
		}
	}
}

func TestReadUint(t *testing.T) {
	tests := []struct {
		in   uint
		want string
	}{
		{in: 0, want: "ゼロ"},
		{in: 1, want: "イチ"},
		// 18,446,744,073,709,551,615
		{in: math.MaxUint64, want: "センハッピャクヨンジュウヨンケイロクセンナナヒャクヨンジュウヨンチョウナナヒャクサンジュウナナオクキュウヒャクゴジュウゴマンセンロッピャクジュウゴ"},
	}

	for _, tt := range tests {
		got := ReadUint(tt.in)
		if got != tt.want {
			t.Errorf("ReadUint(%d) = %s, want %s", tt.in, got, tt.want)
		}
	}
}

func TestReadFloat64(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{in: 0.0, want: "ゼロ"},
		{in: 0.01, want: "レイテンゼロイチ"},
		{in: 3.14, want: "サンテンイチヨン"},
		{in: -1.23, want: "マイナスイッテンニサン"},
	}

	for _, tt := range tests {
		got := ReadFloat64(tt.in)
		if got != tt.want {
			t.Errorf("ReadFloat64(%f) = %s, want %s", tt.in, got, tt.want)
		}
	}
}
