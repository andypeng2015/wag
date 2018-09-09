// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
)

const (
	ptr1 = (1 << 0)
	ptr2 = (1 << 1)
	ptr3 = (1 << 2)
	ptr4 = (1 << 3)
	ptr5 = (1 << 4)
)

var regs = []string{"DI", "SI", "DX", "R10", "R8", "R9"}

type call struct {
	name    string
	number  int
	params  int
	ptrMask int
}

func (sc call) titleName() string {
	return strings.Replace(strings.Title(strings.Replace(sc.name, "_", " ", -1)), " ", "", -1)
}

func main() {
	decl, err := os.Create("syscall.go")
	if err != nil {
		log.Panic(err)
	}
	defer decl.Close()

	impl, err := os.Create("syscall_amd64.s")
	if err != nil {
		log.Panic(err)
	}
	defer impl.Close()

	fmt.Fprintf(decl, "// Generated by internal/cmd/syscalls/generate.go\n\n")
	fmt.Fprintf(decl, "package main\n\n")

	fmt.Fprintf(impl, "// Generated by internal/cmd/syscalls/generate.go\n\n")
	fmt.Fprintf(impl, "#include \"textflag.h\"\n")

	for _, sc := range syscalls {
		fmt.Fprintf(decl, "func import%s() uint64\n", sc.titleName())

		fmt.Fprintf(impl, "\n// func import%s() uint64\n", sc.titleName())
		fmt.Fprintf(impl, "TEXT ·import%s(SB),$0-8\n", sc.titleName())
		fmt.Fprintf(impl, "\tLEAQ\tsys%s<>(SB), AX\n", sc.titleName())
		fmt.Fprintf(impl, "\tMOVQ\tAX, ret+0(FP)\n")
		fmt.Fprintf(impl, "\tRET\n\n")

		fmt.Fprintf(impl, "TEXT sys%s<>(SB),NOSPLIT,$0\n", sc.titleName())

		for i := 0; i < sc.params; i++ {
			r := regs[i]

			fmt.Fprintf(impl, "\tMOVQ\t%d(SP), %s\n", (sc.params-i)*8, r)

			if (sc.ptrMask & (1 << uint(i))) != 0 {
				fmt.Fprintf(impl, "\tANDL\t%s, %s\n", r, r) // zero-extend and test
				fmt.Fprintf(impl, "\tJZ\tnull%d\n", i+1)
				fmt.Fprintf(impl, "\tADDQ\tR14, %s\n", r)
				fmt.Fprintf(impl, "null%d:", i+1)
			}
		}

		fmt.Fprintf(impl, "\tMOVL\t$%d, AX\n", sc.number)
		fmt.Fprintf(impl, "\tSYSCALL\n")
		fmt.Fprintf(impl, "\tXORL\tDX, DX\n")
		fmt.Fprintf(impl, "\tRET\n")
	}

	fmt.Fprintf(decl, "\nfunc init() {\n")

	for _, sc := range syscalls {
		fmt.Fprintf(decl, "\timportFuncs[\"%s\"] = importFunc{import%s(), %d}\n", sc.name, sc.titleName(), sc.params)
	}

	fmt.Fprintf(decl, "}\n") // init()
}

var syscalls = []call{
	{"read", syscall.SYS_READ, 3, ptr2},
	{"write", syscall.SYS_WRITE, 3, ptr2},
	{"open", syscall.SYS_OPEN, 3, ptr1},
	{"close", syscall.SYS_CLOSE, 1, 0},
	{"lseek", syscall.SYS_LSEEK, 3, 0},
	{"pread", syscall.SYS_PREAD64, 4, ptr2},
	{"pwrite", syscall.SYS_PWRITE64, 4, ptr2},
	{"access", syscall.SYS_ACCESS, 2, ptr1},
	{"pipe", syscall.SYS_PIPE, 1, ptr1},
	{"dup", syscall.SYS_DUP, 1, 0},
	{"dup2", syscall.SYS_DUP2, 2, 0},
	{"getpid", syscall.SYS_GETPID, 0, 0},
	{"sendfile", syscall.SYS_SENDFILE, 4, ptr3},
	{"shutdown", syscall.SYS_SHUTDOWN, 2, 0},
	{"socketpair", syscall.SYS_SOCKETPAIR, 4, ptr4},
	{"flock", syscall.SYS_FLOCK, 2, 0},
	{"fsync", syscall.SYS_FSYNC, 1, 0},
	{"fdatasync", syscall.SYS_FDATASYNC, 1, 0},
	{"truncate", syscall.SYS_TRUNCATE, 2, ptr1},
	{"ftruncate", syscall.SYS_FTRUNCATE, 2, 0},
	{"getcwd", syscall.SYS_GETCWD, 2, ptr1},
	{"chdir", syscall.SYS_CHDIR, 1, ptr1},
	{"fchdir", syscall.SYS_FCHDIR, 1, 0},
	{"rename", syscall.SYS_RENAME, 2, ptr1 | ptr2},
	{"mkdir", syscall.SYS_MKDIR, 2, ptr1},
	{"rmdir", syscall.SYS_RMDIR, 1, ptr1},
	{"creat", syscall.SYS_CREAT, 2, ptr1},
	{"link", syscall.SYS_LINK, 2, ptr1 | ptr2},
	{"unlink", syscall.SYS_UNLINK, 1, ptr1},
	{"symlink", syscall.SYS_SYMLINK, 2, ptr1 | ptr2},
	{"readlink", syscall.SYS_READLINK, 3, ptr1 | ptr2},
	{"chmod", syscall.SYS_CHMOD, 2, ptr1},
	{"fchmod", syscall.SYS_FCHMOD, 2, 0},
	{"chown", syscall.SYS_CHOWN, 3, ptr1},
	{"fchown", syscall.SYS_FCHOWN, 3, 0},
	{"lchown", syscall.SYS_LCHOWN, 3, ptr1},
	{"umask", syscall.SYS_UMASK, 1, 0},
	{"getuid", syscall.SYS_GETUID, 0, 0},
	{"getgid", syscall.SYS_GETGID, 0, 0},
	{"vhangup", syscall.SYS_VHANGUP, 0, 0},
	{"sync", syscall.SYS_SYNC, 0, 0},
	{"gettid", syscall.SYS_GETTID, 0, 0},
	{"time", syscall.SYS_TIME, 1, ptr1},
	{"posix_fadvise", syscall.SYS_FADVISE64, 4, 0},
	{"_exit", syscall.SYS_EXIT_GROUP, 1, 0},
	{"inotify_init", syscall.SYS_INOTIFY_INIT, 0, 0},
	{"inotify_add_watch", syscall.SYS_INOTIFY_ADD_WATCH, 3, ptr2},
	{"inotify_rm_watch", syscall.SYS_INOTIFY_RM_WATCH, 2, 0},
	{"openat", syscall.SYS_OPENAT, 4, ptr2},
	{"mkdirat", syscall.SYS_MKDIRAT, 3, ptr2},
	{"fchownat", syscall.SYS_FCHOWNAT, 5, ptr2},
	{"unlinkat", syscall.SYS_UNLINKAT, 3, ptr2},
	{"renameat", syscall.SYS_RENAMEAT, 4, ptr2 | ptr4},
	{"linkat", syscall.SYS_LINKAT, 5, ptr2 | ptr4},
	{"symlinkat", syscall.SYS_SYMLINKAT, 3, ptr1 | ptr3},
	{"readlinkat", syscall.SYS_READLINKAT, 4, ptr2 | ptr3},
	{"fchmodat", syscall.SYS_FCHMODAT, 4, ptr2},
	{"faccessat", syscall.SYS_FACCESSAT, 4, ptr2},
	{"splice", syscall.SYS_SPLICE, 6, ptr2 | ptr4},
	{"tee", syscall.SYS_TEE, 4, 0},
	{"sync_file_range", syscall.SYS_SYNC_FILE_RANGE, 4, 0},
	{"fallocate", syscall.SYS_FALLOCATE, 4, 0},
	{"eventfd", syscall.SYS_EVENTFD2, 2, 0},
	{"dup3", syscall.SYS_DUP3, 3, 0},
	{"pipe2", syscall.SYS_PIPE2, 2, ptr1},
}
