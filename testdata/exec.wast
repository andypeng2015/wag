(module
  (import "spectest" "print" (func $print_i32 (param i32)))
  (import "spectest" "print" (func $print_f32_f64 (param f32 f64)))

  (memory 0)

  (func $main
    (local i32)
    (local i32)
    (drop
      (i32.add
        (i32.const 3)
        (i32.add
          (i32.const 1)
          (i32.const 2))))
    (call $print_i32 (i32.const 10))
    (call $print_f32_f64 (f32.const 10.0) (f64.const 10.0))
    (drop (i32.add (i32.const 1) (i32.const 2)))
    (set_local 0
      (i32.add
        (i32.const 12000)
        (i32.const 345)))
    (call $nop)
    (set_local 1
      (i32.add
        (i32.const 30)
        (i32.const 40)))
    (call $print_i32
      (call $get_result
        (i32.add
          (i32.const 30)
          (i32.const 40))
        (get_local 0)
        (i32.add
          (i32.const 10000)
          (i32.const 2345))))
    (call $stack_overflow
      (call $get (i32.const  0))
      (call $get (i32.const  1))
      (i32.const  2)
      (call $get (i32.const  3))
      (call $get (i32.const  4))
      (call $get (i32.const  5))
      (call $get (i32.const  6))
      (call $get (i32.const  7))
      (call $get (i32.const  8))
      (call $get (i32.const  9))
      (call $get (i32.const 10))
      (call $get (i32.const 11))
      (call $get (i32.const 12))
      (call $get (i32.const 13))
      (call $get (i32.const 14))
      (call $get (i32.const 15))))

  (func $get (param i32) (result i32)
    (get_local 0))

  (func $nop)

  (func)

  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))
  (func (unreachable)) (func (unreachable)) (func (unreachable))

  (func $get_result (param i32) (param i32) (param i32) (result i32)
    (call $print_i32 (i32.const 20))
    (get_local 2))

  (func $stack_overflow
    (param i32) (param i32) (param i32) (param i32)
    (param i32) (param i32) (param i32) (param i32)
    (param i32) (param i32) (param i32) (param i32)
    (param i32) (param i32) (param i32) (param i32)
    (call $print_i32 (get_local 15))
    (call $print_i32 (get_local 14))
    (call $print_i32 (get_local 13))
    (call $print_i32 (get_local 12))
    (call $print_i32 (get_local 11))
    (call $print_i32 (get_local 10))
    (call $print_i32 (get_local 9))
    (call $print_i32 (get_local 8))
    (call $print_i32 (get_local 7))
    (call $print_i32 (get_local 6))
    (call $print_i32 (get_local 5))
    (call $print_i32 (get_local 4))
    (call $print_i32 (get_local 3))
    (call $print_i32 (get_local 2))
    (call $print_i32 (get_local 1))
    (call $print_i32 (get_local 0)))

  (start $main)
)
