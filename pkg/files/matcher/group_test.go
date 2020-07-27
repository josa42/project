package matcher

import "testing"

func TestGroup_String(t *testing.T) {
	type fields struct {
		str  string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group{
				str:  tt.fields.str,
				name: tt.fields.name,
			}
			if got := g.String(); got != tt.want {
				t.Errorf("Group.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Pattern(t *testing.T) {
	type fields struct {
		str  string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group{
				str:  tt.fields.str,
				name: tt.fields.name,
			}
			if got := g.Pattern(); got != tt.want {
				t.Errorf("Group.Pattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Name(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		// TODO: Add test cases.
		// {"empty", "", map[int]string{}},
		// {"named", "{**:path}/file.js", map[int]string{0: "path"}},
		// {"transformed", "{**|dashed}/file.js", map[int]string{0: "path"}},
		// {"named and transformed", "{**|dashed:path}/file.js", map[int]string{0: "path"}},
		// {"constant", "src/{controllers}/file.js", map[int]string{0: "path"}},
		// {"constant and named", "src/{controllers:type}/file.js", map[int]string{0: "type"}},
		// {"constant, transformed and named", "src/{controllers|dashed:type}/file.js", map[int]string{0: "type"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Group{str: tt.path}).Name(); got != tt.want {
				t.Errorf("Group.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Transform(t *testing.T) {
	type fields struct {
		str  string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group{
				str:  tt.fields.str,
				name: tt.fields.name,
			}
			if got := g.Transform(); got != tt.want {
				t.Errorf("Group.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_IsConstant(t *testing.T) {
	type fields struct {
		str  string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group{
				str:  tt.fields.str,
				name: tt.fields.name,
			}
			if got := g.IsConstant(); got != tt.want {
				t.Errorf("Group.IsConstant() = %v, want %v", got, tt.want)
			}
		})
	}
}
