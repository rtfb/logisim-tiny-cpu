/*
 * Compute Fibonacci sequence. Fill in RAM with 1, 1, 2, 3, etc, until the sum
 * overflows 8 bits. Store the result after all the reserved locations.
 */

LI 1
GETACC r7     // fib(1)=1
GETACC r6     // fib(2)=1

LI 0
ST r6         // store [r6] to [0]
INC
ST r7         // store [r7] to [1]
INC
GETACC r5     // index to write result to (r5 = 2)

LI 7
GETACC r1
ADD r1
GETACC r1
LI 1
ADD r1
GETACC r1     // r1=15, the length of the loop, we’ll use it to jump out to avoid
              // writing an overflowed last value to RAM

GETPC r4
LI 6
ADD r4
GETACC r4     // r4 now contains address of loop
ADD r1
GETACC r1     // r1 now contains address of the end of the loop

// loop:
    SETACC r7
    GETACC r2 // preserve r7 in r2 temporarily
    SETACC r6
    ADD r7
    GETACC r7 // r7 now contains the larger Fib number
    SETACC r2
    GETACC r6 // r6 now contains the smaller Fib number

    SETACC r1 // load end of the loop to acc
    JO        // jump out of the loop if we’re done

    SETACC r5 // acc = index
    ST r7     // store the latest Fib number in mem
    INC
    GETACC r5

    SETACC r4
    JMP

HALT
