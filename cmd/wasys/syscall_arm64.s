// Generated by internal/cmd/syscalls/generate.go

#include "textflag.h"

// func importRead() uint64
TEXT ·importRead(SB),$0-8
	BL	after

	MOVD	16(R29), R0
	MOVD	$0, R1
	MOVW	8(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	0(R29), R2
	MOVD	$63, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importWrite() uint64
TEXT ·importWrite(SB),$0-8
	BL	after

	MOVD	16(R29), R0
	MOVD	$0, R1
	MOVW	8(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	0(R29), R2
	MOVD	$64, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importClose() uint64
TEXT ·importClose(SB),$0-8
	BL	after

	MOVD	0(R29), R0
	MOVD	$57, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importLseek() uint64
TEXT ·importLseek(SB),$0-8
	BL	after

	MOVD	16(R29), R0
	MOVD	8(R29), R1
	MOVD	0(R29), R2
	MOVD	$62, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importPread() uint64
TEXT ·importPread(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	$0, R1
	MOVW	16(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$67, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importPwrite() uint64
TEXT ·importPwrite(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	$0, R1
	MOVW	16(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$68, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importDup() uint64
TEXT ·importDup(SB),$0-8
	BL	after

	MOVD	0(R29), R0
	MOVD	$23, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importGetpid() uint64
TEXT ·importGetpid(SB),$0-8
	BL	after

	MOVD	$172, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importSendfile() uint64
TEXT ·importSendfile(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	16(R29), R1
	MOVD	$0, R2
	MOVW	8(R29), R2
	CMPW	$0, R2
	BEQ	null3
	ADD	R26, R2
null3:	MOVD	0(R29), R3
	MOVD	$71, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importShutdown() uint64
TEXT ·importShutdown(SB),$0-8
	BL	after

	MOVD	8(R29), R0
	MOVD	0(R29), R1
	MOVD	$210, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importSocketpair() uint64
TEXT ·importSocketpair(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	16(R29), R1
	MOVD	8(R29), R2
	MOVD	$0, R3
	MOVW	0(R29), R3
	CMPW	$0, R3
	BEQ	null4
	ADD	R26, R3
null4:	MOVD	$199, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFlock() uint64
TEXT ·importFlock(SB),$0-8
	BL	after

	MOVD	8(R29), R0
	MOVD	0(R29), R1
	MOVD	$32, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFsync() uint64
TEXT ·importFsync(SB),$0-8
	BL	after

	MOVD	0(R29), R0
	MOVD	$82, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFdatasync() uint64
TEXT ·importFdatasync(SB),$0-8
	BL	after

	MOVD	0(R29), R0
	MOVD	$83, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importTruncate() uint64
TEXT ·importTruncate(SB),$0-8
	BL	after

	MOVD	$0, R0
	MOVW	8(R29), R0
	CMPW	$0, R0
	BEQ	null1
	ADD	R26, R0
null1:	MOVD	0(R29), R1
	MOVD	$45, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFtruncate() uint64
TEXT ·importFtruncate(SB),$0-8
	BL	after

	MOVD	8(R29), R0
	MOVD	0(R29), R1
	MOVD	$46, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importGetcwd() uint64
TEXT ·importGetcwd(SB),$0-8
	BL	after

	MOVD	$0, R0
	MOVW	8(R29), R0
	CMPW	$0, R0
	BEQ	null1
	ADD	R26, R0
null1:	MOVD	0(R29), R1
	MOVD	$17, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importChdir() uint64
TEXT ·importChdir(SB),$0-8
	BL	after

	MOVD	$0, R0
	MOVW	0(R29), R0
	CMPW	$0, R0
	BEQ	null1
	ADD	R26, R0
null1:	MOVD	$49, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFchdir() uint64
TEXT ·importFchdir(SB),$0-8
	BL	after

	MOVD	0(R29), R0
	MOVD	$50, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFchmod() uint64
TEXT ·importFchmod(SB),$0-8
	BL	after

	MOVD	8(R29), R0
	MOVD	0(R29), R1
	MOVD	$52, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFchown() uint64
TEXT ·importFchown(SB),$0-8
	BL	after

	MOVD	16(R29), R0
	MOVD	8(R29), R1
	MOVD	0(R29), R2
	MOVD	$55, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importLchown() uint64
TEXT ·importLchown(SB),$0-8
	BL	after

	MOVD	$0, R0
	MOVW	16(R29), R0
	CMPW	$0, R0
	BEQ	null1
	ADD	R26, R0
null1:	MOVD	8(R29), R1
	MOVD	0(R29), R2
	MOVD	$1032, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importUmask() uint64
TEXT ·importUmask(SB),$0-8
	BL	after

	MOVD	0(R29), R0
	MOVD	$166, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importGetuid() uint64
TEXT ·importGetuid(SB),$0-8
	BL	after

	MOVD	$174, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importGetgid() uint64
TEXT ·importGetgid(SB),$0-8
	BL	after

	MOVD	$176, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importVhangup() uint64
TEXT ·importVhangup(SB),$0-8
	BL	after

	MOVD	$58, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importSync() uint64
TEXT ·importSync(SB),$0-8
	BL	after

	MOVD	$81, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importGettid() uint64
TEXT ·importGettid(SB),$0-8
	BL	after

	MOVD	$178, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importTime() uint64
TEXT ·importTime(SB),$0-8
	BL	after

	MOVD	$0, R0
	MOVW	0(R29), R0
	CMPW	$0, R0
	BEQ	null1
	ADD	R26, R0
null1:	MOVD	$1062, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importPosixFadvise() uint64
TEXT ·importPosixFadvise(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	16(R29), R1
	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$223, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importExit() uint64
TEXT ·importExit(SB),$0-8
	BL	after

	MOVD	0(R29), R0
	MOVD	$94, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importInotifyInit1() uint64
TEXT ·importInotifyInit1(SB),$0-8
	BL	after

	MOVD	$26, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importInotifyAddWatch() uint64
TEXT ·importInotifyAddWatch(SB),$0-8
	BL	after

	MOVD	16(R29), R0
	MOVD	$0, R1
	MOVW	8(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	0(R29), R2
	MOVD	$27, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importInotifyRmWatch() uint64
TEXT ·importInotifyRmWatch(SB),$0-8
	BL	after

	MOVD	8(R29), R0
	MOVD	0(R29), R1
	MOVD	$28, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importOpenat() uint64
TEXT ·importOpenat(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	$0, R1
	MOVW	16(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$56, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importMkdirat() uint64
TEXT ·importMkdirat(SB),$0-8
	BL	after

	MOVD	16(R29), R0
	MOVD	$0, R1
	MOVW	8(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	0(R29), R2
	MOVD	$34, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFchownat() uint64
TEXT ·importFchownat(SB),$0-8
	BL	after

	MOVD	32(R29), R0
	MOVD	$0, R1
	MOVW	24(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	16(R29), R2
	MOVD	8(R29), R3
	MOVD	0(R29), R4
	MOVD	$54, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importUnlinkat() uint64
TEXT ·importUnlinkat(SB),$0-8
	BL	after

	MOVD	16(R29), R0
	MOVD	$0, R1
	MOVW	8(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	0(R29), R2
	MOVD	$35, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importRenameat() uint64
TEXT ·importRenameat(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	$0, R1
	MOVW	16(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	8(R29), R2
	MOVD	$0, R3
	MOVW	0(R29), R3
	CMPW	$0, R3
	BEQ	null4
	ADD	R26, R3
null4:	MOVD	$38, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importLinkat() uint64
TEXT ·importLinkat(SB),$0-8
	BL	after

	MOVD	32(R29), R0
	MOVD	$0, R1
	MOVW	24(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	16(R29), R2
	MOVD	$0, R3
	MOVW	8(R29), R3
	CMPW	$0, R3
	BEQ	null4
	ADD	R26, R3
null4:	MOVD	0(R29), R4
	MOVD	$37, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importSymlinkat() uint64
TEXT ·importSymlinkat(SB),$0-8
	BL	after

	MOVD	$0, R0
	MOVW	16(R29), R0
	CMPW	$0, R0
	BEQ	null1
	ADD	R26, R0
null1:	MOVD	8(R29), R1
	MOVD	$0, R2
	MOVW	0(R29), R2
	CMPW	$0, R2
	BEQ	null3
	ADD	R26, R2
null3:	MOVD	$36, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importReadlinkat() uint64
TEXT ·importReadlinkat(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	$0, R1
	MOVW	16(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	$0, R2
	MOVW	8(R29), R2
	CMPW	$0, R2
	BEQ	null3
	ADD	R26, R2
null3:	MOVD	0(R29), R3
	MOVD	$78, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFchmodat() uint64
TEXT ·importFchmodat(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	$0, R1
	MOVW	16(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$53, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFaccessat() uint64
TEXT ·importFaccessat(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	$0, R1
	MOVW	16(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$48, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importSplice() uint64
TEXT ·importSplice(SB),$0-8
	BL	after

	MOVD	40(R29), R0
	MOVD	$0, R1
	MOVW	32(R29), R1
	CMPW	$0, R1
	BEQ	null2
	ADD	R26, R1
null2:	MOVD	24(R29), R2
	MOVD	$0, R3
	MOVW	16(R29), R3
	CMPW	$0, R3
	BEQ	null4
	ADD	R26, R3
null4:	MOVD	8(R29), R4
	MOVD	0(R29), R5
	MOVD	$76, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importTee() uint64
TEXT ·importTee(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	16(R29), R1
	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$77, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importSyncFileRange() uint64
TEXT ·importSyncFileRange(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	16(R29), R1
	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$84, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importFallocate() uint64
TEXT ·importFallocate(SB),$0-8
	BL	after

	MOVD	24(R29), R0
	MOVD	16(R29), R1
	MOVD	8(R29), R2
	MOVD	0(R29), R3
	MOVD	$47, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importEventfd() uint64
TEXT ·importEventfd(SB),$0-8
	BL	after

	MOVD	8(R29), R0
	MOVD	0(R29), R1
	MOVD	$19, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importDup3() uint64
TEXT ·importDup3(SB),$0-8
	BL	after

	MOVD	16(R29), R0
	MOVD	8(R29), R1
	MOVD	0(R29), R2
	MOVD	$24, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET

// func importPipe2() uint64
TEXT ·importPipe2(SB),$0-8
	BL	after

	MOVD	$0, R0
	MOVW	8(R29), R0
	CMPW	$0, R0
	BEQ	null1
	ADD	R26, R0
null1:	MOVD	0(R29), R1
	MOVD	$59, R8
	SVC
	MOVD	R27, R1
	ADD	$16, R1
	B	(R1)

after:	MOVD	LR, ret+0(FP)
	RET
