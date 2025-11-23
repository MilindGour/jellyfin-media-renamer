package filesystem

import (
	"io/fs"
	"reflect"
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/testdata"
)

func TestJmrFS_ScanDirectory(t *testing.T) {
	type fields struct {
		fs fs.FS
	}
	type args struct {
		dirpath            string
		includedExtensions []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []DirEntry
	}{
		{
			name: "Mock dir scan",
			fields: fields{
				fs: testdata.MockFSStructure,
			},
			args: args{
				dirpath:            "../testdata/fs-structure",
				includedExtensions: []string{".x"},
			},
			want: []DirEntry{
				{
					Name:        "testdir",
					Path:        "../testdata/fs-structure/testdir",
					Size:        3072,
					IsDirectory: true,
					Children: []DirEntry{
						{
							Name:        "testfile2.x",
							Path:        "../testdata/fs-structure/testdir/testfile2.x",
							Size:        3072,
							IsDirectory: false,
							Children:    nil,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JmrFS{
				fs: tt.fields.fs,
			}
			if got := j.ScanDirectory(tt.args.dirpath, tt.args.includedExtensions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JmrFS.ScanDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJmrFS_GetDirectorySize(t *testing.T) {
	type fields struct {
		fs fs.FS
	}
	type args struct {
		in DirEntry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "Non empty directory size",
			fields: fields{
				fs: nil,
			},
			args: args{
				in: DirEntry{
					Name:        "Test Dir1",
					Path:        "test/path/1",
					Size:        0,
					IsDirectory: true,
					Children: []DirEntry{
						{
							Name:        "Test File",
							Path:        "test/path/1/f1",
							Size:        42,
							IsDirectory: false,
							Children:    nil,
						},
						{
							Name:        "Test File 2",
							Path:        "test/path/1/f2",
							Size:        53,
							IsDirectory: false,
							Children:    nil,
						},
					},
				},
			},
			want: 95,
		},
		{
			name: "Empty directory size",
			fields: fields{
				fs: nil,
			},
			args: args{
				in: DirEntry{
					Name:        "Empty dir2",
					Path:        "test/path/2",
					Size:        12,
					IsDirectory: true,
					Children:    []DirEntry{},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JmrFS{
				fs: tt.fields.fs,
			}
			if got := j.GetDirectorySize(tt.args.in); got != tt.want {
				t.Errorf("JmrFS.GetDirectorySize() = %v, want %v", got, tt.want)
			}
		})
	}
}

// NOTE: These test cases are commented because they take time, and are not required for now.
// func TestJmrFS_MoveFile(t *testing.T) {
// 	jmrFS := NewJmrFS()
//
// 	progressChannel := make(chan FileTransferProgress)
// 	go jmrFS.MoveFile(
// 		"/Users/milindgour/Documents/workspace/personal/test-structure/rsync_fs/src/",
// 		"/Users/milindgour/Documents/workspace/personal/test-structure/rsync_fs/dest/",
// 		progressChannel,
// 	)
//
// 	for progress := range progressChannel {
// 		fmt.Printf("[TEST] %s\n", progress.ToString())
// 	}
// }

// func TestJmrFS_CreateDirectory(t *testing.T) {
// 	// create a test directory
// 	jfs := NewJmrFS()
//
// 	jfs.CreateDirectory("/Users/milindgour/Documents/workspace/personal/test-structure/rsync_fs/.jmr-renames/another/depth/of-dirs")
// }

func TestJmrFS_GetMountPointInfo(t *testing.T) {
	type fields struct {
		fs fs.FS
	}
	type args struct {
		mountPoint string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []MountPointInfo
	}{
		{
			name: "Test using 1 path",
			fields: fields{
				fs: NewJmrFS().fs,
			},
			args: args{
				mountPoint: "/",
			},
			want: []MountPointInfo{
				{
					TotalSizeKB: 123123,
					FreeSizeKB:  12312,
					UsedSizeKB:  1234,
					MountPoint:  "/",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JmrFS{
				fs: tt.fields.fs,
			}
			got := j.GetMountPointInfo(tt.args.mountPoint)
			if got.MountPoint != tt.args.mountPoint {
				t.Errorf("GetMountPointInfo() = %s, want %s", got.MountPoint, tt.args.mountPoint)
				return
			}
			if got.TotalSizeKB == 0 || got.UsedSizeKB == 0 || got.FreeSizeKB == 0 {
				t.Errorf("GetMountPointInfo() total: %d, used: %d, free: %d. Want non-zero values", got.TotalSizeKB, got.UsedSizeKB, got.FreeSizeKB)
				return
			}
		})
	}
}
