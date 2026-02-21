package sorter

import "testing"

func Test_extractKey(t *testing.T) {
	type args struct {
		line      string
		delimiter string
		column    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"whole line w/o delimiter flag",
			args{
				"a\tb\tc",
				"",
				0,
			},
			"a\tb\tc",
		},
		{
			"no delimiter in line",
			args{
				"abc",
				"\t",
				2,
			},
			"",
		},
		{
			"empty delimiter second column",
			args{
				"a\tb\tc",
				"",
				2,
			},
			"b",
		},
		{
			"custom delimiter",
			args{
				"a::b::c",
				"::",
				2,
			},
			"b",
		},
		{
			"endline delimiter empty last column",
			args{
				"a::b::c::",
				"::",
				4,
			},
			"",
		},
		{
			"endline delimiter missing column",
			args{
				"a::b::c::",
				"::",
				5,
			},
			"",
		},
		{
			"only custom delimiter",
			args{
				"::",
				"::",
				2,
			},
			"",
		},
		{
			"negative column",
			args{
				"a\tb\tc",
				"\t",
				-1,
			},
			"a\tb\tc",
		},
		{
			"first column",
			args{
				"a\tb\tc",
				"\t",
				1,
			},
			"a",
		},
		{
			"second column",
			args{
				"a\tb\tc",
				"\t",
				2,
			},
			"b",
		},
		{
			"third column",
			args{
				"a\tb\tc",
				"\t",
				3,
			},
			"c",
		},
		{
			"missing column",
			args{
				"a\tb",
				"\t",
				3,
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractKey(tt.args.line, tt.args.delimiter, tt.args.column); got != tt.want {
				t.Fatalf("extractKey(%q, %q, %d) = %q, want %q",
					tt.args.line,
					tt.args.delimiter,
					tt.args.column,
					got,
					tt.want,
				)
			}
		})
	}
}
