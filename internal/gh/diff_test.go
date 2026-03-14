package gh

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseUnifiedDiff(t *testing.T) {
	tests := []struct {
		name string
		diff string
		want []DiffFile
	}{
		{
			name: "status and counts",
			diff: `diff --git a/new.txt b/new.txt
new file mode 100644
index 0000000..1111111
--- /dev/null
+++ b/new.txt
@@ -0,0 +2 @@
+line1
+line2
diff --git a/old.txt b/old.txt
deleted file mode 100644
index 1111111..0000000
--- a/old.txt
+++ /dev/null
@@ -2,0 +0,0 @@
-line1
-line2
diff --git a/before.txt b/after.txt
similarity index 98%
rename from before.txt
rename to after.txt
@@ -1 +1 @@
-old
+new`,
			want: []DiffFile{
				{
					Path: "new.txt",
					Content: `diff --git a/new.txt b/new.txt
new file mode 100644
index 0000000..1111111
--- /dev/null
+++ b/new.txt
@@ -0,0 +2 @@
+line1
+line2`,
					Status:    DiffFileStatusAdded,
					Additions: 2,
					Deletions: 0,
				},
				{
					Path: "old.txt",
					Content: `diff --git a/old.txt b/old.txt
deleted file mode 100644
index 1111111..0000000
--- a/old.txt
+++ /dev/null
@@ -2,0 +0,0 @@
-line1
-line2`,
					Status:    DiffFileStatusDeleted,
					Additions: 0,
					Deletions: 2,
				},
				{
					Path: "after.txt",
					Content: `diff --git a/before.txt b/after.txt
similarity index 98%
rename from before.txt
rename to after.txt
@@ -1 +1 @@
-old
+new`,
					Status:    DiffFileStatusRenamed,
					Additions: 1,
					Deletions: 1,
				},
			},
		},
		{
			name: "deleted file path with spaces",
			diff: `diff --git a/gone file.txt b/gone file.txt
deleted file mode 100644
index 587be6b..0000000
--- a/gone file.txt	
+++ /dev/null
@@ -1 +0,0 @@
-x`,
			want: []DiffFile{
				{
					Path: "gone file.txt",
					Content: `diff --git a/gone file.txt b/gone file.txt
deleted file mode 100644
index 587be6b..0000000
--- a/gone file.txt	
+++ /dev/null
@@ -1 +0,0 @@
-x`,
					Status:    DiffFileStatusDeleted,
					Additions: 0,
					Deletions: 1,
				},
			},
		},
		{
			name: "mode only path with spaces",
			diff: `diff --git a/script file.sh b/script file.sh
old mode 100644
new mode 100755`,
			want: []DiffFile{
				{
					Path: "script file.sh",
					Content: `diff --git a/script file.sh b/script file.sh
old mode 100644
new mode 100755`,
					Status:    DiffFileStatusType,
					Additions: 0,
					Deletions: 0,
				},
			},
		},
		{
			name: "rename and copy target path preserves a b prefix",
			diff: `diff --git a/old.txt b/b/new.txt
similarity index 100%
rename from old.txt
rename to b/new.txt
diff --git a/source.txt b/a/copied.txt
similarity index 100%
copy from source.txt
copy to a/copied.txt`,
			want: []DiffFile{
				{
					Path: "b/new.txt",
					Content: `diff --git a/old.txt b/b/new.txt
similarity index 100%
rename from old.txt
rename to b/new.txt`,
					Status:    DiffFileStatusRenamed,
					Additions: 0,
					Deletions: 0,
				},
				{
					Path: "a/copied.txt",
					Content: `diff --git a/source.txt b/a/copied.txt
similarity index 100%
copy from source.txt
copy to a/copied.txt`,
					Status:    DiffFileStatusCopied,
					Additions: 0,
					Deletions: 0,
				},
			},
		},
		{
			name: "mode only path containing b slash",
			diff: `diff --git a/foo b/bar.sh b/foo b/bar.sh
old mode 100644
new mode 100755`,
			want: []DiffFile{
				{
					Path: "foo b/bar.sh",
					Content: `diff --git a/foo b/bar.sh b/foo b/bar.sh
old mode 100644
new mode 100755`,
					Status:    DiffFileStatusType,
					Additions: 0,
					Deletions: 0,
				},
			},
		},
		{
			name: "quoted path keeps escaped control sequence literal",
			diff: `diff --git "a/line\nbreak\033.txt" "b/line\nbreak\033.txt"
old mode 100644
new mode 100755`,
			want: []DiffFile{
				{
					Path: `line\nbreak\033.txt`,
					Content: `diff --git "a/line\nbreak\033.txt" "b/line\nbreak\033.txt"
old mode 100644
new mode 100755`,
					Status:    DiffFileStatusType,
					Additions: 0,
					Deletions: 0,
				},
			},
		},
		{
			name: "actual control chars are sanitized from path",
			diff: "diff --git \"a/\x1bbad.txt\" \"b/\x1bbad.txt\"\nold mode 100644\nnew mode 100755",
			want: []DiffFile{
				{
					Path:      "bad.txt",
					Content:   "diff --git \"a/\x1bbad.txt\" \"b/\x1bbad.txt\"\nold mode 100644\nnew mode 100755",
					Status:    DiffFileStatusType,
					Additions: 0,
					Deletions: 0,
				},
			},
		},
		{
			name: "hunk data starting with triple plus minus is counted",
			diff: `diff --git a/a.txt b/a.txt
--- a/a.txt
+++ b/a.txt
@@ -1,2 +1,2 @@
----bar
++++foo`,
			want: []DiffFile{
				{
					Path: "a.txt",
					Content: `diff --git a/a.txt b/a.txt
--- a/a.txt
+++ b/a.txt
@@ -1,2 +1,2 @@
----bar
++++foo`,
					Status:    DiffFileStatusModified,
					Additions: 1,
					Deletions: 1,
				},
			},
		},
		{
			name: "hunk data with triple plus minus and space keeps header parsing intact",
			diff: `diff --git a/a.txt b/a.txt
--- a/a.txt
+++ b/a.txt
@@ -1,3 +1,3 @@
--- foo
+++ bar
-old
+new`,
			want: []DiffFile{
				{
					Path: "a.txt",
					Content: `diff --git a/a.txt b/a.txt
--- a/a.txt
+++ b/a.txt
@@ -1,3 +1,3 @@
--- foo
+++ bar
-old
+new`,
					Status:    DiffFileStatusModified,
					Additions: 2,
					Deletions: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseUnifiedDiff(tt.diff)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("parse result mismatch (-want +got)\n%s", diff)
			}
		})
	}
}

func TestUnifiedDiffParserParse(t *testing.T) {
	diff := `diff --git "a/line\nbreak\033.txt" "b/line\nbreak\033.txt"
old mode 100644
new mode 100755
diff --git a/a.txt b/a.txt
--- a/a.txt
+++ b/a.txt
@@ -1,2 +1,2 @@
----bar
++++foo`

	parser := NewUnifiedDiffParser(diff)
	got := parser.Parse()
	want := ParseUnifiedDiff(diff)

	if d := cmp.Diff(want, got); d != "" {
		t.Fatalf("parse result mismatch (-want +got)\n%s", d)
	}
}
