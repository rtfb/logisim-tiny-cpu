Tiny 8-bit ISA
==============

Fixed-width 8-bit instructions. 5 bits per opcode, 8 general purpose registers,
and a dedicated accumulator register to be used as one parameter for some
operations. This gives me enough of opcode space, but there's an obvious problem
with loading immediate numbers, only 3 bits remain for that. This should be
workable, though, as most programs initialize their variables with a 0 or 1.
Larger values will have to be constructed with math or bitshifts.

The PC register and data memory access will also be 8-bit. This will limit the
programs to be 256 instructions long (and RAM to 256 bytes), but this ISA is not
meant to be useful for large programs, just a toy. This will be big enough to
write a few demo programs.

5 upper bits for opcode: total 32 opcodes.

3 lower bits for argument: up to 3-bit immediate values, OR address 8 general
purpose registers. All two-argument operations will have to use an accumulator
register.

## Registers
* PC - program counter
* ACC - accumulator (implicitly used as specified in opcodes)
* R0-R7 - general purpose registers
* STAT- status flag register (zero, overflow, etc, TBD)

## Opcodes

Notes:
* All instructions are 5-bit opcode + 3-bit parameter
* rX means "a specified general purpose register"

| Dec | TopHex |  Bin  | Mnemonic | Description                             |
| --- | ------ | ----- | -------- | --------------------------------------- |
|   0 |      0 | 00000 |     HALT | stops the CPU                           |
|   1 |      0 | 00001 |       LI | load immediate 3 bits into acc          |
|   2 |      1 | 00010 |       LD | load memory to rX, from address stored at acc |
|   3 |      1 | 00011 |       ST | store contents of rX to an address stored at acc |
|   4 |      2 | 00100 |    GETPC | rX = pc                                 |
|   5 |      2 | 00101 |    GETST | rX = STAT                               |
|   6 |      3 | 00110 |    SETST | STAT = rX                               |
|   7 |      3 | 00111 |     SHLI | shift left acc by up to 3 bits          |
|   8 |      4 | 01000 |     SHRI | shift right acc by up to 3 bits         |
|   9 |      4 | 01001 |   GETACC | rX = acc                                |
|  10 |      5 | 01010 |   SETACC | acc = rX                                |
|  11 |      5 | 01011 |          |                                         |
|  12 |      6 | 01100 |       OR | bitwise OR acc and rX, store result in acc |
|  13 |      6 | 01101 |      AND | bitwise AND acc and rX, store result in acc |
|  14 |      7 | 01110 |      XOR | bitwise XOR acc and rX, store result in acc |
|  15 |      7 | 01111 |      ADD | acc = acc + rX                          |
|  16 |      8 | 10000 |      SUB | acc = acc - rX                          |
|  17 |      8 | 10001 |      INC | acc += immediate                        |
|  18 |      9 | 10010 |      DEC | acc -= immediate                        |
|  19 |      9 | 10011 |          |                                         |
|  20 |      a | 10100 |          |                                         |
|  21 |      a | 10101 |          |                                         |
|  22 |      b | 10110 |          |                                         |
|  23 |      b | 10111 |          |                                         |
|  24 |      c | 11000 |          |                                         |
|  25 |      c | 11001 |          |                                         |
|  26 |      d | 11010 |      LI0 | load immediate 4 bits into acc (0b0xxx) |
|  27 |      d | 11011 |      LI1 | load immediate 4 bits into acc (0b1xxx) |
|  28 |      e | 11100 |      SJF | set jump flavor: set the combination of flags to be respected by the next jump. SJF 0 means unconditional jump |
|  29 |      e | 11101 |     SJFN | same as SJF, but the flags will be negated before considering them |
|  30 |      f | 11110 |    JMPLO | jump to an absolute address 0bAAAA0XXX, where AAAA are the four bits in ACC, 0 is a literal zero bit, and XXX is the three-bit immediate param |
|  31 |      f | 11111 |    JMPHI | jump to an absolute address 0bAAAA1XXX, where AAAA are the four bits in ACC, 1 is a literal one bit, and XXX is the three-bit immediate param |

## Status Register

| 7 | 6 | 5 | 4 | 3 |     2    |    1     |   0  |
| - | - | - | - | - | -------- | -------- | ---- |
| x | x | x | x | x | Negative | Overflow | Zero |

NOTE: the jump instructions can be combined into a single opcode with 3 bits of
parameter for the flavor. (But my opcode space is not saturated yet, so no need
for that yet).

## Proof of concept program

Compute Fibonacci sequence. Fill in RAM with 2, 1, 2, 3, etc, until the sum
overflows 8 bits. Store the result after all the reserved locations.

```
LI 1
GETACC r7    // fib 1
GETACC r6    // fib 2

LI 0
ST r6  	// store [r6] to [0]
INC
ST r7  	// store [r7] to [1]
INC
GETACC r5    // index to write result to (r5 = 2)

LI 7
GETACC r1
ADD r1
GETACC r1
LI 3
ADD r1
GETACC r1    // r1=17, the length of the loop, we’ll use it to jump out to avoid
             // writing an overflowed last value to RAM

GETPC r4
LI 6
ADD r4
GETACC r4   // r4 now contains address of loop
ADD r1
GETACC r1   // r1 now contains address of the end of the loop

loop:
	SETACC r7
	GETACC r2 // preserve r7 in r2 temporarily
	SETACC r6
	ADD r7
	SETACC r1 // load end of the loop to acc
	JO        // jump out of the loop if we’re done
	GETST r3  // save status for use before jump

	GETACC r7 // r7 now contains the larger Fib number
	SETACC r2
	GETACC r6 // r6 now contains the smaller Fib number

	SETACC r5 // acc = index
	ST r7	// store the latest Fib number in mem
	INC
	GETACC r5

	SETACC r4
	SETST r3
	JNO

HALT
```
