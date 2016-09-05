;; i32 operations

(module
  (func $add (param $x i32) (param $y i32) (result i32) (i32.add (get_local $x) (get_local $y)))
  (func $sub (param $x i32) (param $y i32) (result i32) (i32.sub (get_local $x) (get_local $y)))
  ; (func $mul (param $x i32) (param $y i32) (result i32) (i32.mul (get_local $x) (get_local $y)))
  ; (func $div_s (param $x i32) (param $y i32) (result i32) (i32.div_s (get_local $x) (get_local $y)))
  ; (func $div_u (param $x i32) (param $y i32) (result i32) (i32.div_u (get_local $x) (get_local $y)))
  ; (func $rem_s (param $x i32) (param $y i32) (result i32) (i32.rem_s (get_local $x) (get_local $y)))
  ; (func $rem_u (param $x i32) (param $y i32) (result i32) (i32.rem_u (get_local $x) (get_local $y)))
  (func $and (param $x i32) (param $y i32) (result i32) (i32.and (get_local $x) (get_local $y)))
  (func $or (param $x i32) (param $y i32) (result i32) (i32.or (get_local $x) (get_local $y)))
  (func $xor (param $x i32) (param $y i32) (result i32) (i32.xor (get_local $x) (get_local $y)))
  ; (func $shl (param $x i32) (param $y i32) (result i32) (i32.shl (get_local $x) (get_local $y)))
  ; (func $shr_s (param $x i32) (param $y i32) (result i32) (i32.shr_s (get_local $x) (get_local $y)))
  ; (func $shr_u (param $x i32) (param $y i32) (result i32) (i32.shr_u (get_local $x) (get_local $y)))
  ; (func $rotl (param $x i32) (param $y i32) (result i32) (i32.rotl (get_local $x) (get_local $y)))
  ; (func $rotr (param $x i32) (param $y i32) (result i32) (i32.rotr (get_local $x) (get_local $y)))
  ; (func $clz (param $x i32) (result i32) (i32.clz (get_local $x)))
  ; (func $ctz (param $x i32) (result i32) (i32.ctz (get_local $x)))
  ; (func $popcnt (param $x i32) (result i32) (i32.popcnt (get_local $x)))
  (func $eqz (param $x i32) (result i32) (i32.eqz (get_local $x)))
  ; (func $eq (param $x i32) (param $y i32) (result i32) (i32.eq (get_local $x) (get_local $y)))
  (func $ne (param $x i32) (param $y i32) (result i32) (i32.ne (get_local $x) (get_local $y)))
  ; (func $lt_s (param $x i32) (param $y i32) (result i32) (i32.lt_s (get_local $x) (get_local $y)))
  ; (func $lt_u (param $x i32) (param $y i32) (result i32) (i32.lt_u (get_local $x) (get_local $y)))
  ; (func $le_s (param $x i32) (param $y i32) (result i32) (i32.le_s (get_local $x) (get_local $y)))
  ; (func $le_u (param $x i32) (param $y i32) (result i32) (i32.le_u (get_local $x) (get_local $y)))
  ; (func $gt_s (param $x i32) (param $y i32) (result i32) (i32.gt_s (get_local $x) (get_local $y)))
  ; (func $gt_u (param $x i32) (param $y i32) (result i32) (i32.gt_u (get_local $x) (get_local $y)))
  ; (func $ge_s (param $x i32) (param $y i32) (result i32) (i32.ge_s (get_local $x) (get_local $y)))
  ; (func $ge_u (param $x i32) (param $y i32) (result i32) (i32.ge_u (get_local $x) (get_local $y)))

  (export "add" $add)
  (export "sub" $sub)
  ; (export "mul" $mul)
  ; (export "div_s" $div_s)
  ; (export "div_u" $div_u)
  ; (export "rem_s" $rem_s)
  ; (export "rem_u" $rem_u)
  (export "and" $and)
  (export "or" $or)
  (export "xor" $xor)
  ; (export "shl" $shl)
  ; (export "shr_s" $shr_s)
  ; (export "shr_u" $shr_u)
  ; (export "rotl" $rotl)
  ; (export "rotr" $rotr)
  ; (export "clz" $clz)
  ; (export "ctz" $ctz)
  ; (export "popcnt" $popcnt)
  (export "eqz" $eqz)
  ; (export "eq" $eq)
  (export "ne" $ne)
  ; (export "lt_s" $lt_s)
  ; (export "lt_u" $lt_u)
  ; (export "le_s" $le_s)
  ; (export "le_u" $le_u)
  ; (export "gt_s" $gt_s)
  ; (export "gt_u" $gt_u)
  ; (export "ge_s" $ge_s)
  ; (export "ge_u" $ge_u)
)

(assert_return (invoke "add" (i32.const 1) (i32.const 1)) (i32.const 2))
(assert_return (invoke "add" (i32.const 1) (i32.const 0)) (i32.const 1))
(assert_return (invoke "add" (i32.const -1) (i32.const -1)) (i32.const -2))
(assert_return (invoke "add" (i32.const -1) (i32.const 1)) (i32.const 0))
(assert_return (invoke "add" (i32.const 0x7fffffff) (i32.const 1)) (i32.const 0x80000000))
(assert_return (invoke "add" (i32.const 0x80000000) (i32.const -1)) (i32.const 0x7fffffff))
(assert_return (invoke "add" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 0))
(assert_return (invoke "add" (i32.const 0x3fffffff) (i32.const 1)) (i32.const 0x40000000))

(assert_return (invoke "sub" (i32.const 1) (i32.const 1)) (i32.const 0))
(assert_return (invoke "sub" (i32.const 1) (i32.const 0)) (i32.const 1))
(assert_return (invoke "sub" (i32.const -1) (i32.const -1)) (i32.const 0))
(assert_return (invoke "sub" (i32.const 0x7fffffff) (i32.const -1)) (i32.const 0x80000000))
(assert_return (invoke "sub" (i32.const 0x80000000) (i32.const 1)) (i32.const 0x7fffffff))
(assert_return (invoke "sub" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 0))
(assert_return (invoke "sub" (i32.const 0x3fffffff) (i32.const -1)) (i32.const 0x40000000))

; (assert_return (invoke "mul" (i32.const 1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "mul" (i32.const 1) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "mul" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "mul" (i32.const 0x10000000) (i32.const 4096)) (i32.const 0))
; (assert_return (invoke "mul" (i32.const 0x80000000) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "mul" (i32.const 0x80000000) (i32.const -1)) (i32.const 0x80000000))
; (assert_return (invoke "mul" (i32.const 0x7fffffff) (i32.const -1)) (i32.const 0x80000001))
; (assert_return (invoke "mul" (i32.const 0x01234567) (i32.const 0x76543210)) (i32.const 0x358e7470))

; (assert_trap (invoke "div_s" (i32.const 1) (i32.const 0)) "integer divide by zero")
; (assert_trap (invoke "div_s" (i32.const 0) (i32.const 0)) "integer divide by zero")
; (assert_trap (invoke "div_s" (i32.const 0x80000000) (i32.const -1)) "integer overflow")
; (assert_return (invoke "div_s" (i32.const 1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "div_s" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "div_s" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "div_s" (i32.const 0x80000000) (i32.const 2)) (i32.const 0xc0000000))
; (assert_return (invoke "div_s" (i32.const 0x80000001) (i32.const 1000)) (i32.const 0xffdf3b65))
; (assert_return (invoke "div_s" (i32.const 5) (i32.const 2)) (i32.const 2))
; (assert_return (invoke "div_s" (i32.const -5) (i32.const 2)) (i32.const -2))
; (assert_return (invoke "div_s" (i32.const 5) (i32.const -2)) (i32.const -2))
; (assert_return (invoke "div_s" (i32.const -5) (i32.const -2)) (i32.const 2))
; (assert_return (invoke "div_s" (i32.const 7) (i32.const 3)) (i32.const 2))
; (assert_return (invoke "div_s" (i32.const -7) (i32.const 3)) (i32.const -2))
; (assert_return (invoke "div_s" (i32.const 7) (i32.const -3)) (i32.const -2))
; (assert_return (invoke "div_s" (i32.const -7) (i32.const -3)) (i32.const 2))
; (assert_return (invoke "div_s" (i32.const 11) (i32.const 5)) (i32.const 2))
; (assert_return (invoke "div_s" (i32.const 17) (i32.const 7)) (i32.const 2))

; (assert_trap (invoke "div_u" (i32.const 1) (i32.const 0)) "integer divide by zero")
; (assert_trap (invoke "div_u" (i32.const 0) (i32.const 0)) "integer divide by zero")
; (assert_return (invoke "div_u" (i32.const 1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "div_u" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "div_u" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "div_u" (i32.const 0x80000000) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "div_u" (i32.const 0x80000000) (i32.const 2)) (i32.const 0x40000000))
; (assert_return (invoke "div_u" (i32.const 0x8ff00ff0) (i32.const 0x10001)) (i32.const 0x8fef))
; (assert_return (invoke "div_u" (i32.const 0x80000001) (i32.const 1000)) (i32.const 0x20c49b))
; (assert_return (invoke "div_u" (i32.const 5) (i32.const 2)) (i32.const 2))
; (assert_return (invoke "div_u" (i32.const -5) (i32.const 2)) (i32.const 0x7ffffffd))
; (assert_return (invoke "div_u" (i32.const 5) (i32.const -2)) (i32.const 0))
; (assert_return (invoke "div_u" (i32.const -5) (i32.const -2)) (i32.const 0))
; (assert_return (invoke "div_u" (i32.const 7) (i32.const 3)) (i32.const 2))
; (assert_return (invoke "div_u" (i32.const 11) (i32.const 5)) (i32.const 2))
; (assert_return (invoke "div_u" (i32.const 17) (i32.const 7)) (i32.const 2))

; (assert_trap (invoke "rem_s" (i32.const 1) (i32.const 0)) "integer divide by zero")
; (assert_trap (invoke "rem_s" (i32.const 0) (i32.const 0)) "integer divide by zero")
; (assert_return (invoke "rem_s" (i32.const 0x7fffffff) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "rem_s" (i32.const 1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "rem_s" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "rem_s" (i32.const -1) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "rem_s" (i32.const 0x80000000) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "rem_s" (i32.const 0x80000000) (i32.const 2)) (i32.const 0))
; (assert_return (invoke "rem_s" (i32.const 0x80000001) (i32.const 1000)) (i32.const -647))
; (assert_return (invoke "rem_s" (i32.const 5) (i32.const 2)) (i32.const 1))
; (assert_return (invoke "rem_s" (i32.const -5) (i32.const 2)) (i32.const -1))
; (assert_return (invoke "rem_s" (i32.const 5) (i32.const -2)) (i32.const 1))
; (assert_return (invoke "rem_s" (i32.const -5) (i32.const -2)) (i32.const -1))
; (assert_return (invoke "rem_s" (i32.const 7) (i32.const 3)) (i32.const 1))
; (assert_return (invoke "rem_s" (i32.const -7) (i32.const 3)) (i32.const -1))
; (assert_return (invoke "rem_s" (i32.const 7) (i32.const -3)) (i32.const 1))
; (assert_return (invoke "rem_s" (i32.const -7) (i32.const -3)) (i32.const -1))
; (assert_return (invoke "rem_s" (i32.const 11) (i32.const 5)) (i32.const 1))
; (assert_return (invoke "rem_s" (i32.const 17) (i32.const 7)) (i32.const 3))

; (assert_trap (invoke "rem_u" (i32.const 1) (i32.const 0)) "integer divide by zero")
; (assert_trap (invoke "rem_u" (i32.const 0) (i32.const 0)) "integer divide by zero")
; (assert_return (invoke "rem_u" (i32.const 1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "rem_u" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "rem_u" (i32.const -1) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "rem_u" (i32.const 0x80000000) (i32.const -1)) (i32.const 0x80000000))
; (assert_return (invoke "rem_u" (i32.const 0x80000000) (i32.const 2)) (i32.const 0))
; (assert_return (invoke "rem_u" (i32.const 0x8ff00ff0) (i32.const 0x10001)) (i32.const 0x8001))
; (assert_return (invoke "rem_u" (i32.const 0x80000001) (i32.const 1000)) (i32.const 649))
; (assert_return (invoke "rem_u" (i32.const 5) (i32.const 2)) (i32.const 1))
; (assert_return (invoke "rem_u" (i32.const -5) (i32.const 2)) (i32.const 1))
; (assert_return (invoke "rem_u" (i32.const 5) (i32.const -2)) (i32.const 5))
; (assert_return (invoke "rem_u" (i32.const -5) (i32.const -2)) (i32.const -5))
; (assert_return (invoke "rem_u" (i32.const 7) (i32.const 3)) (i32.const 1))
; (assert_return (invoke "rem_u" (i32.const 11) (i32.const 5)) (i32.const 1))
; (assert_return (invoke "rem_u" (i32.const 17) (i32.const 7)) (i32.const 3))

(assert_return (invoke "and" (i32.const 1) (i32.const 0)) (i32.const 0))
(assert_return (invoke "and" (i32.const 0) (i32.const 1)) (i32.const 0))
(assert_return (invoke "and" (i32.const 1) (i32.const 1)) (i32.const 1))
(assert_return (invoke "and" (i32.const 0) (i32.const 0)) (i32.const 0))
(assert_return (invoke "and" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 0))
(assert_return (invoke "and" (i32.const 0x7fffffff) (i32.const -1)) (i32.const 0x7fffffff))
(assert_return (invoke "and" (i32.const 0xf0f0ffff) (i32.const 0xfffff0f0)) (i32.const 0xf0f0f0f0))
(assert_return (invoke "and" (i32.const 0xffffffff) (i32.const 0xffffffff)) (i32.const 0xffffffff))

(assert_return (invoke "or" (i32.const 1) (i32.const 0)) (i32.const 1))
(assert_return (invoke "or" (i32.const 0) (i32.const 1)) (i32.const 1))
(assert_return (invoke "or" (i32.const 1) (i32.const 1)) (i32.const 1))
(assert_return (invoke "or" (i32.const 0) (i32.const 0)) (i32.const 0))
(assert_return (invoke "or" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const -1))
(assert_return (invoke "or" (i32.const 0x80000000) (i32.const 0)) (i32.const 0x80000000))
(assert_return (invoke "or" (i32.const 0xf0f0ffff) (i32.const 0xfffff0f0)) (i32.const 0xffffffff))
(assert_return (invoke "or" (i32.const 0xffffffff) (i32.const 0xffffffff)) (i32.const 0xffffffff))

(assert_return (invoke "xor" (i32.const 1) (i32.const 0)) (i32.const 1))
(assert_return (invoke "xor" (i32.const 0) (i32.const 1)) (i32.const 1))
(assert_return (invoke "xor" (i32.const 1) (i32.const 1)) (i32.const 0))
(assert_return (invoke "xor" (i32.const 0) (i32.const 0)) (i32.const 0))
(assert_return (invoke "xor" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const -1))
(assert_return (invoke "xor" (i32.const 0x80000000) (i32.const 0)) (i32.const 0x80000000))
(assert_return (invoke "xor" (i32.const -1) (i32.const 0x80000000)) (i32.const 0x7fffffff))
(assert_return (invoke "xor" (i32.const -1) (i32.const 0x7fffffff)) (i32.const 0x80000000))
(assert_return (invoke "xor" (i32.const 0xf0f0ffff) (i32.const 0xfffff0f0)) (i32.const 0x0f0f0f0f))
(assert_return (invoke "xor" (i32.const 0xffffffff) (i32.const 0xffffffff)) (i32.const 0))

; (assert_return (invoke "shl" (i32.const 1) (i32.const 1)) (i32.const 2))
; (assert_return (invoke "shl" (i32.const 1) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "shl" (i32.const 0x7fffffff) (i32.const 1)) (i32.const 0xfffffffe))
; (assert_return (invoke "shl" (i32.const 0xffffffff) (i32.const 1)) (i32.const 0xfffffffe))
; (assert_return (invoke "shl" (i32.const 0x80000000) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "shl" (i32.const 0x40000000) (i32.const 1)) (i32.const 0x80000000))
; (assert_return (invoke "shl" (i32.const 1) (i32.const 31)) (i32.const 0x80000000))
; (assert_return (invoke "shl" (i32.const 1) (i32.const 32)) (i32.const 1))
; (assert_return (invoke "shl" (i32.const 1) (i32.const 33)) (i32.const 2))
; (assert_return (invoke "shl" (i32.const 1) (i32.const -1)) (i32.const 0x80000000))
; (assert_return (invoke "shl" (i32.const 1) (i32.const 0x7fffffff)) (i32.const 0x80000000))

; (assert_return (invoke "shr_s" (i32.const 1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "shr_s" (i32.const 1) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "shr_s" (i32.const -1) (i32.const 1)) (i32.const -1))
; (assert_return (invoke "shr_s" (i32.const 0x7fffffff) (i32.const 1)) (i32.const 0x3fffffff))
; (assert_return (invoke "shr_s" (i32.const 0x80000000) (i32.const 1)) (i32.const 0xc0000000))
; (assert_return (invoke "shr_s" (i32.const 0x40000000) (i32.const 1)) (i32.const 0x20000000))
; (assert_return (invoke "shr_s" (i32.const 1) (i32.const 32)) (i32.const 1))
; (assert_return (invoke "shr_s" (i32.const 1) (i32.const 33)) (i32.const 0))
; (assert_return (invoke "shr_s" (i32.const 1) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "shr_s" (i32.const 1) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "shr_s" (i32.const 1) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "shr_s" (i32.const 0x80000000) (i32.const 31)) (i32.const -1))
; (assert_return (invoke "shr_s" (i32.const -1) (i32.const 32)) (i32.const -1))
; (assert_return (invoke "shr_s" (i32.const -1) (i32.const 33)) (i32.const -1))
; (assert_return (invoke "shr_s" (i32.const -1) (i32.const -1)) (i32.const -1))
; (assert_return (invoke "shr_s" (i32.const -1) (i32.const 0x7fffffff)) (i32.const -1))
; (assert_return (invoke "shr_s" (i32.const -1) (i32.const 0x80000000)) (i32.const -1))

; (assert_return (invoke "shr_u" (i32.const 1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "shr_u" (i32.const 1) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "shr_u" (i32.const -1) (i32.const 1)) (i32.const 0x7fffffff))
; (assert_return (invoke "shr_u" (i32.const 0x7fffffff) (i32.const 1)) (i32.const 0x3fffffff))
; (assert_return (invoke "shr_u" (i32.const 0x80000000) (i32.const 1)) (i32.const 0x40000000))
; (assert_return (invoke "shr_u" (i32.const 0x40000000) (i32.const 1)) (i32.const 0x20000000))
; (assert_return (invoke "shr_u" (i32.const 1) (i32.const 32)) (i32.const 1))
; (assert_return (invoke "shr_u" (i32.const 1) (i32.const 33)) (i32.const 0))
; (assert_return (invoke "shr_u" (i32.const 1) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "shr_u" (i32.const 1) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "shr_u" (i32.const 1) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "shr_u" (i32.const 0x80000000) (i32.const 31)) (i32.const 1))
; (assert_return (invoke "shr_u" (i32.const -1) (i32.const 32)) (i32.const -1))
; (assert_return (invoke "shr_u" (i32.const -1) (i32.const 33)) (i32.const 0x7fffffff))
; (assert_return (invoke "shr_u" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "shr_u" (i32.const -1) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "shr_u" (i32.const -1) (i32.const 0x80000000)) (i32.const -1))

; (assert_return (invoke "rotl" (i32.const 0xfe00dc00) (i32.const 4)) (i32.const 0xe00dc00f))
; (assert_return (invoke "rotl" (i32.const 0xabcd9876) (i32.const 1)) (i32.const 0x579b30ed))
; (assert_return (invoke "rotl" (i32.const 0x00008000) (i32.const 37)) (i32.const 0x00100000))
; (assert_return (invoke "rotl" (i32.const 0x769abcdf) (i32.const 0x8000000d)) (i32.const 0x579beed3))
; (assert_return (invoke "rotl" (i32.const 1) (i32.const 31)) (i32.const 0x80000000))
; (assert_return (invoke "rotl" (i32.const 0x80000000) (i32.const 1)) (i32.const 1))

; (assert_return (invoke "rotr" (i32.const 0xb0c1d2e3) (i32.const 0x0005)) (i32.const 0x1d860e97))
; (assert_return (invoke "rotr" (i32.const 0xb0c1d2e3) (i32.const 0xff05)) (i32.const 0x1d860e97))
; (assert_return (invoke "rotr" (i32.const 0xff00cc00) (i32.const 1)) (i32.const 0x7f806600))
; (assert_return (invoke "rotr" (i32.const 0x00080000) (i32.const 4)) (i32.const 0x00008000))
; (assert_return (invoke "rotr" (i32.const 0x769abcdf) (i32.const 0xffffffed)) (i32.const 0xe6fbb4d5))
; (assert_return (invoke "rotr" (i32.const 1) (i32.const 1)) (i32.const 0x80000000))
; (assert_return (invoke "rotr" (i32.const 0x80000000) (i32.const 31)) (i32.const 1))

; (assert_return (invoke "clz" (i32.const 0xffffffff)) (i32.const 0))
; (assert_return (invoke "clz" (i32.const 0)) (i32.const 32))
; (assert_return (invoke "clz" (i32.const 0x00008000)) (i32.const 16))
; (assert_return (invoke "clz" (i32.const 0xff)) (i32.const 24))
; (assert_return (invoke "clz" (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "clz" (i32.const 1)) (i32.const 31))
; (assert_return (invoke "clz" (i32.const 2)) (i32.const 30))
; (assert_return (invoke "clz" (i32.const 0x7fffffff)) (i32.const 1))

; (assert_return (invoke "ctz" (i32.const -1)) (i32.const 0))
; (assert_return (invoke "ctz" (i32.const 0)) (i32.const 32))
; (assert_return (invoke "ctz" (i32.const 0x00008000)) (i32.const 15))
; (assert_return (invoke "ctz" (i32.const 0x00010000)) (i32.const 16))
; (assert_return (invoke "ctz" (i32.const 0x80000000)) (i32.const 31))
; (assert_return (invoke "ctz" (i32.const 0x7fffffff)) (i32.const 0))

; (assert_return (invoke "popcnt" (i32.const -1)) (i32.const 32))
; (assert_return (invoke "popcnt" (i32.const 0)) (i32.const 0))
; (assert_return (invoke "popcnt" (i32.const 0x00008000)) (i32.const 1))
; (assert_return (invoke "popcnt" (i32.const 0x80008000)) (i32.const 2))
; (assert_return (invoke "popcnt" (i32.const 0x7fffffff)) (i32.const 31))
; (assert_return (invoke "popcnt" (i32.const 0xAAAAAAAA)) (i32.const 16))
; (assert_return (invoke "popcnt" (i32.const 0x55555555)) (i32.const 16))
; (assert_return (invoke "popcnt" (i32.const 0xDEADBEEF)) (i32.const 24))

(assert_return (invoke "eqz" (i32.const 0)) (i32.const 1))
(assert_return (invoke "eqz" (i32.const 1)) (i32.const 0))
(assert_return (invoke "eqz" (i32.const 0x80000000)) (i32.const 0))
(assert_return (invoke "eqz" (i32.const 0x7fffffff)) (i32.const 0))

; (assert_return (invoke "eq" (i32.const 0) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "eq" (i32.const 1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "eq" (i32.const -1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "eq" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "eq" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "eq" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "eq" (i32.const 1) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "eq" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "eq" (i32.const 0x80000000) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "eq" (i32.const 0) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "eq" (i32.const 0x80000000) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "eq" (i32.const -1) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "eq" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "eq" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 0))

(assert_return (invoke "ne" (i32.const 0) (i32.const 0)) (i32.const 0))
(assert_return (invoke "ne" (i32.const 1) (i32.const 1)) (i32.const 0))
(assert_return (invoke "ne" (i32.const -1) (i32.const 1)) (i32.const 1))
(assert_return (invoke "ne" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 0))
(assert_return (invoke "ne" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 0))
(assert_return (invoke "ne" (i32.const -1) (i32.const -1)) (i32.const 0))
(assert_return (invoke "ne" (i32.const 1) (i32.const 0)) (i32.const 1))
(assert_return (invoke "ne" (i32.const 0) (i32.const 1)) (i32.const 1))
(assert_return (invoke "ne" (i32.const 0x80000000) (i32.const 0)) (i32.const 1))
(assert_return (invoke "ne" (i32.const 0) (i32.const 0x80000000)) (i32.const 1))
(assert_return (invoke "ne" (i32.const 0x80000000) (i32.const -1)) (i32.const 1))
(assert_return (invoke "ne" (i32.const -1) (i32.const 0x80000000)) (i32.const 1))
(assert_return (invoke "ne" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 1))
(assert_return (invoke "ne" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 1))

; (assert_return (invoke "lt_s" (i32.const 0) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "lt_s" (i32.const 1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "lt_s" (i32.const -1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "lt_s" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "lt_s" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "lt_s" (i32.const -1) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "lt_s" (i32.const 1) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "lt_s" (i32.const 0) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "lt_s" (i32.const 0x80000000) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "lt_s" (i32.const 0) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "lt_s" (i32.const 0x80000000) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "lt_s" (i32.const -1) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "lt_s" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "lt_s" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 0))

; (assert_return (invoke "lt_u" (i32.const 0) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const 1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const -1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const -1) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const 1) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const 0) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "lt_u" (i32.const 0x80000000) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const 0) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "lt_u" (i32.const 0x80000000) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "lt_u" (i32.const -1) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "lt_u" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 1))

; (assert_return (invoke "le_s" (i32.const 0) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const 1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const -1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const 1) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "le_s" (i32.const 0) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const 0x80000000) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const 0) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "le_s" (i32.const 0x80000000) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const -1) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "le_s" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "le_s" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 0))

; (assert_return (invoke "le_u" (i32.const 0) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "le_u" (i32.const 1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "le_u" (i32.const -1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "le_u" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "le_u" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "le_u" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "le_u" (i32.const 1) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "le_u" (i32.const 0) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "le_u" (i32.const 0x80000000) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "le_u" (i32.const 0) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "le_u" (i32.const 0x80000000) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "le_u" (i32.const -1) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "le_u" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "le_u" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 1))

; (assert_return (invoke "gt_s" (i32.const 0) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const 1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const -1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const -1) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const 1) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "gt_s" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const 0x80000000) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const 0) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "gt_s" (i32.const 0x80000000) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const -1) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "gt_s" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "gt_s" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 1))

; (assert_return (invoke "gt_u" (i32.const 0) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "gt_u" (i32.const 1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "gt_u" (i32.const -1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "gt_u" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "gt_u" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "gt_u" (i32.const -1) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "gt_u" (i32.const 1) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "gt_u" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "gt_u" (i32.const 0x80000000) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "gt_u" (i32.const 0) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "gt_u" (i32.const 0x80000000) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "gt_u" (i32.const -1) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "gt_u" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "gt_u" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 0))

; (assert_return (invoke "ge_s" (i32.const 0) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "ge_s" (i32.const 1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "ge_s" (i32.const -1) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "ge_s" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "ge_s" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "ge_s" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "ge_s" (i32.const 1) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "ge_s" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "ge_s" (i32.const 0x80000000) (i32.const 0)) (i32.const 0))
; (assert_return (invoke "ge_s" (i32.const 0) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "ge_s" (i32.const 0x80000000) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "ge_s" (i32.const -1) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "ge_s" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 0))
; (assert_return (invoke "ge_s" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 1))

; (assert_return (invoke "ge_u" (i32.const 0) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const 1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const -1) (i32.const 1)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const 0x80000000) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const 0x7fffffff) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const -1) (i32.const -1)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const 1) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const 0) (i32.const 1)) (i32.const 0))
; (assert_return (invoke "ge_u" (i32.const 0x80000000) (i32.const 0)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const 0) (i32.const 0x80000000)) (i32.const 0))
; (assert_return (invoke "ge_u" (i32.const 0x80000000) (i32.const -1)) (i32.const 0))
; (assert_return (invoke "ge_u" (i32.const -1) (i32.const 0x80000000)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const 0x80000000) (i32.const 0x7fffffff)) (i32.const 1))
; (assert_return (invoke "ge_u" (i32.const 0x7fffffff) (i32.const 0x80000000)) (i32.const 0))
