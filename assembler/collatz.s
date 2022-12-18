
// Collatz orbit calculator. Expects an input number in mem[0]. Will keep it
// there, and fill mem[1], etc with the Collatz orbit for that number.

li 0
ld r7    // r7 = input number

li 0
xor r7   // will set Z flag if r7 == acc == 0
jnz proceed  // jump over the halt to the actual code if precondition was good
halt

proceed:
li 1
getacc r6    // hold the index in r6

//
// THE LOOP
//
begin_loop:
  li 1
  and r7
  jz  even_case       // jump if r7 was even

  // odd case:
  setacc r7
  shli 1
  add r7
  inc 1
  getacc r7
  setacc r6
  st r7

  inc 1  // index++
  getacc r6
  jmp  begin_loop  // loop back

even_case:
  setacc r7
  shri 1
  getacc r7
  setacc r6
  st r7
  inc 1  // index++
  getacc r6

  li 1
  xor r7
  jz after_loop   // jump to end of program

  jmp begin_loop  // loop back

after_loop:
halt
