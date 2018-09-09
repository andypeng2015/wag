// Generated by internal/cmd/syscalls/generate.go

package main

func importRead() uint64
func importWrite() uint64
func importOpen() uint64
func importClose() uint64
func importLseek() uint64
func importPread() uint64
func importPwrite() uint64
func importAccess() uint64
func importPipe() uint64
func importDup() uint64
func importDup2() uint64
func importGetpid() uint64
func importSendfile() uint64
func importShutdown() uint64
func importSocketpair() uint64
func importFlock() uint64
func importFsync() uint64
func importFdatasync() uint64
func importTruncate() uint64
func importFtruncate() uint64
func importGetcwd() uint64
func importChdir() uint64
func importFchdir() uint64
func importRename() uint64
func importMkdir() uint64
func importRmdir() uint64
func importCreat() uint64
func importLink() uint64
func importUnlink() uint64
func importSymlink() uint64
func importReadlink() uint64
func importChmod() uint64
func importFchmod() uint64
func importChown() uint64
func importFchown() uint64
func importLchown() uint64
func importUmask() uint64
func importGetuid() uint64
func importGetgid() uint64
func importVhangup() uint64
func importSync() uint64
func importGettid() uint64
func importTime() uint64
func importPosixFadvise() uint64
func importExit() uint64
func importInotifyInit() uint64
func importInotifyAddWatch() uint64
func importInotifyRmWatch() uint64
func importOpenat() uint64
func importMkdirat() uint64
func importFchownat() uint64
func importUnlinkat() uint64
func importRenameat() uint64
func importLinkat() uint64
func importSymlinkat() uint64
func importReadlinkat() uint64
func importFchmodat() uint64
func importFaccessat() uint64
func importSplice() uint64
func importTee() uint64
func importSyncFileRange() uint64
func importFallocate() uint64
func importEventfd() uint64
func importDup3() uint64
func importPipe2() uint64

func init() {
	importFuncs["read"] = importFunc{importRead(), 3}
	importFuncs["write"] = importFunc{importWrite(), 3}
	importFuncs["open"] = importFunc{importOpen(), 3}
	importFuncs["close"] = importFunc{importClose(), 1}
	importFuncs["lseek"] = importFunc{importLseek(), 3}
	importFuncs["pread"] = importFunc{importPread(), 4}
	importFuncs["pwrite"] = importFunc{importPwrite(), 4}
	importFuncs["access"] = importFunc{importAccess(), 2}
	importFuncs["pipe"] = importFunc{importPipe(), 1}
	importFuncs["dup"] = importFunc{importDup(), 1}
	importFuncs["dup2"] = importFunc{importDup2(), 2}
	importFuncs["getpid"] = importFunc{importGetpid(), 0}
	importFuncs["sendfile"] = importFunc{importSendfile(), 4}
	importFuncs["shutdown"] = importFunc{importShutdown(), 2}
	importFuncs["socketpair"] = importFunc{importSocketpair(), 4}
	importFuncs["flock"] = importFunc{importFlock(), 2}
	importFuncs["fsync"] = importFunc{importFsync(), 1}
	importFuncs["fdatasync"] = importFunc{importFdatasync(), 1}
	importFuncs["truncate"] = importFunc{importTruncate(), 2}
	importFuncs["ftruncate"] = importFunc{importFtruncate(), 2}
	importFuncs["getcwd"] = importFunc{importGetcwd(), 2}
	importFuncs["chdir"] = importFunc{importChdir(), 1}
	importFuncs["fchdir"] = importFunc{importFchdir(), 1}
	importFuncs["rename"] = importFunc{importRename(), 2}
	importFuncs["mkdir"] = importFunc{importMkdir(), 2}
	importFuncs["rmdir"] = importFunc{importRmdir(), 1}
	importFuncs["creat"] = importFunc{importCreat(), 2}
	importFuncs["link"] = importFunc{importLink(), 2}
	importFuncs["unlink"] = importFunc{importUnlink(), 1}
	importFuncs["symlink"] = importFunc{importSymlink(), 2}
	importFuncs["readlink"] = importFunc{importReadlink(), 3}
	importFuncs["chmod"] = importFunc{importChmod(), 2}
	importFuncs["fchmod"] = importFunc{importFchmod(), 2}
	importFuncs["chown"] = importFunc{importChown(), 3}
	importFuncs["fchown"] = importFunc{importFchown(), 3}
	importFuncs["lchown"] = importFunc{importLchown(), 3}
	importFuncs["umask"] = importFunc{importUmask(), 1}
	importFuncs["getuid"] = importFunc{importGetuid(), 0}
	importFuncs["getgid"] = importFunc{importGetgid(), 0}
	importFuncs["vhangup"] = importFunc{importVhangup(), 0}
	importFuncs["sync"] = importFunc{importSync(), 0}
	importFuncs["gettid"] = importFunc{importGettid(), 0}
	importFuncs["time"] = importFunc{importTime(), 1}
	importFuncs["posix_fadvise"] = importFunc{importPosixFadvise(), 4}
	importFuncs["_exit"] = importFunc{importExit(), 1}
	importFuncs["inotify_init"] = importFunc{importInotifyInit(), 0}
	importFuncs["inotify_add_watch"] = importFunc{importInotifyAddWatch(), 3}
	importFuncs["inotify_rm_watch"] = importFunc{importInotifyRmWatch(), 2}
	importFuncs["openat"] = importFunc{importOpenat(), 4}
	importFuncs["mkdirat"] = importFunc{importMkdirat(), 3}
	importFuncs["fchownat"] = importFunc{importFchownat(), 5}
	importFuncs["unlinkat"] = importFunc{importUnlinkat(), 3}
	importFuncs["renameat"] = importFunc{importRenameat(), 4}
	importFuncs["linkat"] = importFunc{importLinkat(), 5}
	importFuncs["symlinkat"] = importFunc{importSymlinkat(), 3}
	importFuncs["readlinkat"] = importFunc{importReadlinkat(), 4}
	importFuncs["fchmodat"] = importFunc{importFchmodat(), 4}
	importFuncs["faccessat"] = importFunc{importFaccessat(), 4}
	importFuncs["splice"] = importFunc{importSplice(), 6}
	importFuncs["tee"] = importFunc{importTee(), 4}
	importFuncs["sync_file_range"] = importFunc{importSyncFileRange(), 4}
	importFuncs["fallocate"] = importFunc{importFallocate(), 4}
	importFuncs["eventfd"] = importFunc{importEventfd(), 2}
	importFuncs["dup3"] = importFunc{importDup3(), 3}
	importFuncs["pipe2"] = importFunc{importPipe2(), 2}
}
