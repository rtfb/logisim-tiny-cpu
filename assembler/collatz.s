
// Collatz orbit calculator. Expects an input number in mem[0]. Will keep it
// there, and fill mem[1], etc with the Collatz orbit for that number.

li 0
ld r7    // r7 = input number

li 0
xor r7   // will set Z flag if r7 == acc == 0
getst r6
getpc r5
li 6
add r5
setst r6
jnz       // jump over the halt to the actual code if precondition was good
halt

li 1
getacc r6    // hold the index in r6

//
// calculate jump addresses and store them in r0, etc
//

// set r0 to begin_loop
li 3
shli 3
getacc r0
li 7
or r0
getacc r0   // r0=31

// set r1 to even case
li 5
shli 3
getacc r1
li 6
or r1
getacc r1  // r1=46

// set r2 to after_loop
li 7
shli 3
getacc r2
li 3
or r2
getacc r2   // r2=59

//
// THE LOOP
//

// begin_loop:
li 1
and r7
setacc r1
jz        // jump if r7 was even

  // odd case:
  setacc r7
  shli 1
  add r7
  inc
  getacc r7
  setacc r6
  st r7

  inc    // index++
  getacc r6

  setacc r0
  jmp   // loop back

  // even case:
  setacc r7
  shri 1
  getacc r7
  setacc r6
  st r7
  inc    // index++
  getacc r6

  li 1
  xor r7
  setacc r2
  jz      // jump to end of program

  setacc r0
  jmp   // loop back

// after_loop:
halt
