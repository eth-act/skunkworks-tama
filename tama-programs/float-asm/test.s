.text
.globl _start

_start:
    # Test floating-point with proper IEEE 754 numbers
    
    # Use RAM address for storing IEEE 754 values
    li x1, 0xa0020000  # RAM address that's safe
    
    # Store IEEE 754 single-precision representation of 3.14f (0x4048F5C3)
    li x2, 0x4048F5C3
    sw x2, 0(x1)
    
    # Store IEEE 754 single-precision representation of 2.71f (0x402D70A4)
    li x3, 0x402D70A4
    sw x3, 4(x1)
    
    # Load IEEE 754 values into floating-point registers
    flw f1, 0(x1)      # Load 3.14f into f1
    flw f2, 4(x1)      # Load 2.71f into f2
    
    # Test floating-point addition: 3.14 + 2.71 = 5.85
    fadd.s f3, f1, f2  # This should use our softfloat emulation
    
    # Store result back to memory for verification
    fsw f3, 8(x1)      # Store result at 0xa0020008
    
    # Load the result back and compare with expected value
    lw x4, 8(x1)       # Load the computed result
    
    # Expected IEEE 754 representation of 5.85f is 0x40BB3333
    li x5, 0x40BB3333  # Expected result
    
    # Compare actual vs expected
    bne x4, x5, test_failed
    
    # Test passed - write 'P' (Pass) to UART and exit
    li x6, 0xa0000200  # UART address
    li x7, 80          # ASCII 'P' for Pass
    sb x7, 0(x6)       # Write 'P' to UART
    li a7, 93          # Exit syscall number  
    ecall

test_failed:
    # Test failed - write 'F' (Fail) to UART and exit
    li x6, 0xa0000200  # UART address
    li x7, 70          # ASCII 'F' for Fail
    sb x7, 0(x6)       # Write 'F' to UART
    li a7, 93          # Exit syscall number
    ecall

    # Loop in case ecall doesn't immediately exit
exit_loop:
    j exit_loop