MS-DOS Emulator v4.0 - User Manual
Table of Contents

Introduction
Installation
Getting Started
Shell Commands
File Operations
System Commands
Emulator Commands
Running Programs
Debugging Features
Technical Specifications
Examples
Troubleshooting
FAQ
Appendices


Introduction
MS-DOS Emulator v4.0 is a complete implementation of an MS-DOS compatible environment written in Go. It emulates an Intel 8086 CPU with full BIOS and DOS interrupt support.
Key Features

Complete 8086 CPU with all 16 registers
Full flags register (CF, PF, AF, ZF, SF, TF, IF, DF, OF)
Over 100 instructions implemented
BIOS interrupt support (INT 10h, 13h, 16h, 1Ah, 33h)
DOS interrupt support (INT 21h with 50+ functions)
Interactive shell with file system operations
Advanced debugging tools (breakpoints, step mode, disassembler)
Memory dump and analysis tools
Performance statistics


Installation
Prerequisites

Go 1.16 or higher
Terminal or command prompt

Building from Source
# Navigate to the project directory
cd dos-emulator

# Initialize Go module
go mod init dos-emulator

# Build the emulator
go build -o dos-emulator main.go

# Run the emulator
./dos-emulator

Windows Build
go build -o dos-emulator.exe main.go
dos-emulator.exe


Getting Started
Starting the Emulator
Run the executable to start the interactive shell:
./dos-emulator

You will see:
╔══════════════════════════════════════════════════════════════╗
║     MS-DOS Emulator v4.0 - COMPLETE IMPLEMENTATION          ║
║     Full 8086 CPU + BIOS + DOS + 100+ Instructions          ║
╚══════════════════════════════════════════════════════════════╝

Type 'HELP' for available commands

A:\>

Command Line Usage
# Interactive shell
./dos-emulator

# Run a COM file directly
./dos-emulator program.com

# Run with debug mode enabled
./dos-emulator -d program.com

# Show help
./dos-emulator -h
./dos-emulator --help


Shell Commands
HELP - Display Available Commands
Shows all available commands with descriptions.
Usage:
HELP
?

Output:
╔══════════════════════════════════════════════════════════════╗
║                    AVAILABLE COMMANDS                        ║
╠══════════════════════════════════════════════════════════════╣
║ File: DIR, CD, MD, RD, DEL, TYPE, COPY, REN                  ║
║ System: CLS, VER, DATE, TIME, MEM, ECHO                      ║
║ Emulator: RUN, DEBUG, STEP, TRACE, REGS, BREAK, DUMP        ║
║           STACK, STATS, DISASM, EXIT                         ║
╚══════════════════════════════════════════════════════════════╝


File Operations
DIR - List Directory Contents
Lists all files and subdirectories in the current directory.
Usage:
DIR

Example:
A:\> DIR

 Volume in drive A is EMULATOR
 Directory of A:\

main.go      15234 12-27-24  08:30p
README.md     8456 12-27-24  08:35p
test.com       256 12-27-24  09:00p
programs  <DIR>         12-27-24  09:15p
    3 File(s) 23946 bytes
    1 Dir(s)

A:\>

CD - Change Directory
Changes the current working directory or displays the current directory path.
Usage:
CD [path]
CD          (displays current directory)

Examples:
A:\> CD programs
A:\> CD
A:\programs

A:\> CD ..
A:\> CD
A:\

A:\> CD \programs\games

MD / MKDIR - Make Directory
Creates a new directory.
Usage:
MD <directory>
MKDIR <directory>

Examples:
A:\> MD programs
A:\> MKDIR games
A:\> MD \temp\data

RD / RMDIR - Remove Directory
Removes an empty directory.
Usage:
RD <directory>
RMDIR <directory>

Examples:
A:\> RD olddir
A:\> RMDIR \temp\old

Note: Directory must be empty to be removed.
DEL / ERASE - Delete File
Deletes one or more files.
Usage:
DEL <filename>
ERASE <filename>

Examples:
A:\> DEL oldfile.txt
A:\> ERASE temp.dat
A:\> DEL *.bak

TYPE - Display File Contents
Displays the contents of a text file to the screen.
Usage:
TYPE <filename>

Example:
A:\> TYPE readme.txt
This is a sample text file.
It contains multiple lines.
End of file.

A:\>

COPY - Copy Files
Copies a file from source to destination.
Usage:
COPY <source> <destination>

Examples:
A:\> COPY file1.txt file2.txt
        1 file(s) copied

A:\> COPY readme.txt \backup\readme.txt
        1 file(s) copied

REN / RENAME - Rename File
Renames a file or directory.
Usage:
REN <oldname> <newname>
RENAME <oldname> <newname>

Examples:
A:\> REN old.txt new.txt
A:\> RENAME data.dat backup.dat


System Commands
CLS - Clear Screen
Clears the terminal screen.
Usage:
CLS

Example:
A:\> CLS
(screen clears)

A:\>

VER - Show Version
Displays emulator version information and execution statistics.
Usage:
VER

Example:
A:\> VER
MS-DOS Emulator Version 4.0 - Complete Implementation
Instructions executed: 125678

A:\>

DATE - Display Date
Shows the current system date.
Usage:
DATE

Example:
A:\> DATE
Current date: Sat 12/27/2024

A:\>

TIME - Display Time
Shows the current system time.
Usage:
TIME

Example:
A:\> TIME
Current time: 20:45:32

A:\>

MEM - Memory Information
Displays memory usage and statistics.
Usage:
MEM

Example:
A:\> MEM

Memory: 1024 KB total, 640 KB conventional

A:\>

ECHO - Display Message
Displays text to the screen.
Usage:
ECHO <message>

Examples:
A:\> ECHO Hello World!
Hello World!

A:\> ECHO System is ready
System is ready

A:\>


Emulator Commands
RUN / EXEC - Execute Program
Loads and executes a COM or EXE file.
Usage:
RUN <filename>
EXEC <filename>
<filename>          (direct execution)

Examples:
A:\> RUN test.com
Running program...
Hello from test program!
Program exited with code: 0

A:\> EXEC game.com
Running program...

A:\> test.com
Running program...

DEBUG - Toggle Debug Mode
Enables or disables debug mode for detailed execution information.
Usage:
DEBUG

Example:
A:\> DEBUG
Debug mode: true

A:\> RUN test.com
1000:0100  MOV AH, 0x09  AX=0000 BX=0000 CX=0000 DX=0000
1000:0102  MOV DX, 0x010E  AX=0900 BX=0000 CX=0000 DX=0000
...

A:\> DEBUG
Debug mode: false

When Debug Mode is Active:

Each instruction is displayed before execution
Register values are shown
Memory accesses are logged
Interrupt calls are detailed

STEP - Toggle Step Mode
Enables step-by-step execution with manual control.
Usage:
STEP

Example:
A:\> STEP
Step mode: true

A:\> RUN test.com
1000:0100  MOV AH, 0x09  AX=0000 BX=0000 CX=0000 DX=0000
Press Enter (c=continue, q=quit, r=registers)> 

Step Mode Controls:

Enter - Execute next instruction
c - Continue execution without stepping
q - Quit program execution
r - Show register dump

Example Session:
Press Enter (c=continue, q=quit, r=registers)> [Enter]
1000:0102  MOV DX, 0x010E  AX=0900 BX=0000 CX=0000 DX=0000
Press Enter (c=continue, q=quit, r=registers)> r

╔══════════════════════════════════════════════════════════════╗
║                      CPU REGISTERS                           ║
╠══════════════════════════════════════════════════════════════╣
║ AX=0900  BX=0000  CX=0000  DX=0000                        ║
...

Press Enter (c=continue, q=quit, r=registers)> c
Running program...

TRACE - Toggle Trace Mode
Enables trace mode to display each instruction during execution.
Usage:
TRACE

Example:
A:\> TRACE
Trace mode: true

A:\> RUN test.com
1000:0100  MOV AH, 0x09  AX=0000 BX=0000 CX=0000 DX=0000
1000:0102  MOV DX, 0x010E  AX=0900 BX=0000 CX=0000 DX=0000
1000:0105  INT 0x21  AX=0900 BX=0000 CX=0000 DX=010E
Hello World!
...

Difference from Debug Mode:

Trace mode shows instructions but doesn't pause
Debug mode provides more detailed information
Trace is useful for following program flow

REGS - Show CPU Registers
Displays all CPU registers and flags.
Usage:
REGS

Example:
A:\> REGS

╔══════════════════════════════════════════════════════════════╗
║                      CPU REGISTERS                           ║
╠══════════════════════════════════════════════════════════════╣
║ AX=1234  BX=5678  CX=9ABC  DX=DEF0                        ║
║ SI=0000  DI=0000  BP=0000  SP=FFFE                        ║
║ CS=1000  DS=1000  ES=1000  SS=1000                        ║
║ IP=0100  FLAGS=0002                                        ║
║ Flags: ZF                                                  ║
╚══════════════════════════════════════════════════════════════╝

A:\>

Flag Indicators:

CF - Carry Flag
ZF - Zero Flag
SF - Sign Flag
OF - Overflow Flag
PF - Parity Flag
AF - Auxiliary Carry Flag
IF - Interrupt Flag
DF - Direction Flag
TF - Trap Flag

BREAK - Set/List Breakpoints
Sets breakpoints at specific memory addresses or lists all active breakpoints.
Usage:
BREAK [address]
BREAK               (list all breakpoints)

Examples:
A:\> BREAK 1100
Breakpoint set at 00001100

A:\> BREAK 1200
Breakpoint set at 00001200

A:\> BREAK

Active breakpoints:
  00001100
  00001200

A:\>

Using Breakpoints:
A:\> BREAK 1150
Breakpoint set at 00001150

A:\> RUN test.com
Running program...

*** Breakpoint hit at 1000:0150 ***
Press Enter (c=continue, q=quit, r=registers)>

DUMP - Memory Dump
Displays raw memory contents in hexadecimal and ASCII format.
Usage:
DUMP [address] [length]
DUMP [address]      (dumps 256 bytes)
DUMP                (dumps from address 0)

Examples:
A:\> DUMP 100

Memory dump from 00000100:
00000100: B4 09 BA 0E 01 CD 21 B4 4C CD 21 48 65 6C 6C 6F  | ......!.L.!Hello
00000110: 20 57 6F 72 6C 64 21 0D 0A 24 00 00 00 00 00 00  |  World!..$......
00000120: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  | ................
...

A:\> DUMP 100 64

Memory dump from 00000100:
00000100: B4 09 BA 0E 01 CD 21 B4 4C CD 21 48 65 6C 6C 6F  | ......!.L.!Hello
00000110: 20 57 6F 72 6C 64 21 0D 0A 24 00 00 00 00 00 00  |  World!..$......
00000120: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  | ................
00000130: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  | ................

A:\>

Format:

First column: Memory address in hexadecimal
Middle columns: Hexadecimal byte values
Last column: ASCII representation (. for non-printable)

STACK - Show Stack Contents
Displays the top entries of the stack.
Usage:
STACK

Example:
A:\> STACK

Stack (top 10):
  [00] 0103
  [01] 1000
  [02] 0246
  [03] 5678
  [04] 1234

A:\>

Stack Information:

Shows up to 10 most recent stack entries
Index [00] is the top of stack
Values are displayed in hexadecimal
Useful for debugging CALL/RET sequences

STATS - Show Statistics
Displays emulator performance and execution statistics.
Usage:
STATS

Example:
A:\> STATS

╔══════════════════════════════════════════════════════════════╗
║                   EMULATOR STATISTICS                        ║
╠══════════════════════════════════════════════════════════════╣
║ Instructions: 1234567                                        ║
║ Running time: 2.456s                                         ║
║ IPS:          502500                                         ║
╚══════════════════════════════════════════════════════════════╝

A:\>

Statistics Explained:

Instructions - Total number of instructions executed
Running time - Total execution time
IPS - Instructions Per Second (performance metric)

DISASM - Disassemble Code
Disassembles machine code into assembly language instructions.
Usage:
DISASM [address] [count]
DISASM [address]        (disassembles 20 instructions)
DISASM                  (disassembles from current IP)

Examples:
A:\> DISASM

Disassembly from 00010100:
00010100: MOV AH, 0x09
00010102: MOV DX, 0x010E
00010105: INT 0x21
00010107: MOV AH, 0x4C
00010109: INT 0x21
0001010B: UNKNOWN (0x48)
...

A:\> DISASM 100 10

Disassembly from 00000100:
00000100: MOV AH, 0x09
00000102: MOV DX, 0x010E
00000105: INT 0x21
00000107: MOV AH, 0x4C
00000109: INT 0x21
0000010B: NOP
0000010C: NOP
0000010D: NOP
0000010E: UNKNOWN (0x48)
0000010F: UNKNOWN (0x65)

A:\>

Use Cases:

Analyzing compiled programs
Verifying code generation
Understanding program flow
Debugging unknown code

EXIT / QUIT - Exit Emulator
Exits the emulator and returns to the operating system.
Usage:
EXIT
QUIT

Example:
A:\> EXIT
Exiting emulator...
$


Running Programs
Creating a Simple Program
Here's how to create a basic "Hello World" program:
hello.asm:
; Hello World program for MS-DOS
; Assemble with: nasm -f bin -o hello.com hello.asm

ORG 100h                ; COM files start at offset 100h

section .text
    mov ah, 09h         ; DOS function: print string
    mov dx, message     ; Address of message
    int 21h             ; Call DOS

    mov ah, 4Ch         ; DOS function: exit program
    mov al, 0           ; Return code 0
    int 21h             ; Call DOS

section .data
    message db 'Hello, World!', 0Dh, 0Ah, '$'

Assembling the Program
Using NASM assembler:
nasm -f bin -o hello.com hello.asm

Running in the Emulator
Method 1: Direct execution
A:\> hello.com
Running program...
Hello, World!
Program exited with code: 0

A:\>

Method 2: Using RUN command
A:\> RUN hello.com
Running program...
Hello, World!
Program exited with code: 0

A:\>

Method 3: Command line
./dos-emulator hello.com

Program Structure
COM File Format:

Starts at offset 100h (256 bytes)
Maximum size: approximately 64KB
No header, pure binary code
CS, DS, ES, SS all point to PSP segment
Entry point is at 100h

Basic Program Template:
ORG 100h

; Your code here
mov ah, 09h
mov dx, message
int 21h

; Exit
mov ah, 4Ch
int 21h

message db 'Text here$'


Debugging Features
Complete Debugging Workflow
Step 1: Enable Debug Mode
A:\> DEBUG
Debug mode: true

Step 2: Set Breakpoints
A:\> BREAK 1100
Breakpoint set at 00001100

A:\> BREAK 1150
Breakpoint set at 00001150

Step 3: Run Program
A:\> RUN test.com
1000:0100  MOV AH, 0x09  AX=0000 BX=0000 CX=0000 DX=0000
1000:0102  MOV DX, 0x010E  AX=0900 BX=0000 CX=0000 DX=0000
...
*** Breakpoint hit at 1000:0100 ***

Step 4: Examine State
A:\> REGS
(shows all registers)

A:\> STACK
(shows stack contents)

A:\> DUMP 100 128
(shows memory contents)

Step 5: Continue or Step
A:\> STEP
Step mode: true

A:\> RUN test.com
(step through each instruction)

Advanced Debugging Techniques
Technique 1: Trace Execution Flow
A:\> TRACE
Trace mode: true

A:\> RUN program.com
(watch instruction flow)

Technique 2: Memory Analysis
A:\> DUMP 100
(examine program code)

A:\> DUMP 200
(examine data area)

A:\> DISASM 100 50
(disassemble code)

Technique 3: Conditional Debugging
A:\> BREAK 1200
A:\> RUN program.com
(program runs until breakpoint)

Press Enter (c=continue, q=quit, r=registers)> r
(examine registers)

Press Enter (c=continue, q=quit, r=registers)> c
(continue execution)

Debugging Example Session
A:\> DEBUG
Debug mode: true

A:\> STEP
Step mode: true

A:\> RUN calc.com
1000:0100  MOV AX, 0x0005  AX=0000 BX=0000 CX=0000 DX=0000
Press Enter (c=continue, q=quit, r=registers)> [Enter]

1000:0103  MOV BX, 0x0003  AX=0005 BX=0000 CX=0000 DX=0000
Press Enter (c=continue, q=quit, r=registers)> [Enter]

1000:0106  ADD AX, BX  AX=0005 BX=0003 CX=0000 DX=0000
Press Enter (c=continue, q=quit, r=registers)> [Enter]

1000:0108  ADD AL, 0x30  AX=0008 BX=0003 CX=0000 DX=0000
Press Enter (c=continue, q=quit, r=registers)> r

╔══════════════════════════════════════════════════════════════╗
║                      CPU REGISTERS                           ║
╠══════════════════════════════════════════════════════════════╣
║ AX=0008  BX=0003  CX=0000  DX=0000                        ║
║ SI=0000  DI=0000  BP=0000  SP=FFFE                        ║
║ CS=1000  DS=1000  ES=1000  SS=1000                        ║
║ IP=010A  FLAGS=0002                                        ║
║ Flags:                                                     ║
╚══════════════════════════════════════════════════════════════╝

Press Enter (c=continue, q=quit, r=registers)> c
Running program...
8
Program exited with code: 0

A:\> STATS

╔══════════════════════════════════════════════════════════════╗
║                   EMULATOR STATISTICS                        ║
╠══════════════════════════════════════════════════════════════╣
║ Instructions: 12                                             ║
║ Running time: 45.123s                                        ║
║ IPS:          0                                              ║
╚══════════════════════════════════════════════════════════════╝

A:\>


Technical Specifications
CPU Architecture
Register Set:



Register
Size
Description



AX
16-bit
Accumulator (AH:AL)


BX
16-bit
Base (BH:BL)


CX
16-bit
Counter (CH:CL)


DX
16-bit
Data (DH:DL)


SI
16-bit
Source Index


DI
16-bit
Destination Index


BP
16-bit
Base Pointer


SP
16-bit
Stack Pointer


CS
16-bit
Code Segment


DS
16-bit
Data Segment


ES
16-bit
Extra Segment


SS
16-bit
Stack Segment


IP
16-bit
Instruction Pointer


FLAGS
16-bit
Flags Register


Flags Register (16-bit):
Bit  15 14 13 12 11 10  9  8  7  6  5  4  3  2  1  0
     -  -  -  - OF DF IF TF SF ZF  - AF  - PF  - CF




Flag
Bit
Name
Description



CF
0
Carry
Arithmetic carry/borrow


PF
2
Parity
Even parity of low byte


AF
4
Auxiliary
BCD carry/borrow


ZF
6
Zero
Result is zero


SF
7
Sign
Sign of result


TF
8
Trap
Single-step mode


IF
9
Interrupt
Interrupt enable


DF
10
Direction
String operation direction


OF
11
Overflow
Signed overflow


Memory Architecture
Memory Map:



Address Range
Size
Description



00000-003FF
1 KB
Interrupt Vector Table


00400-004FF
256 B
BIOS Data Area


00500-9FFFF
638 KB
Conventional Memory


A0000-BFFFF
128 KB
Video Memory


C0000-FFFFF
256 KB
ROM BIOS


Total Addressable Memory: 1 MB (0x00000 - 0xFFFFF)
Segment:Offset Addressing:
Physical Address = (Segment × 16) + Offset
Example: CS:IP = 1000:0100 → 10000 + 0100 = 10100h

Instruction Set
Data Transfer Instructions:

MOV - Move data
PUSH - Push onto stack
POP - Pop from stack
XCHG - Exchange
LEA - Load effective address
LDS - Load pointer using DS
LES - Load pointer using ES
LAHF - Load AH from flags
SAHF - Store AH to flags
PUSHF - Push flags
POPF - Pop flags

Arithmetic Instructions:

ADD - Addition
ADC - Add with carry
SUB - Subtraction
SBB - Subtract with borrow
MUL - Unsigned multiply
IMUL - Signed multiply
DIV - Unsigned divide
IDIV - Signed divide
INC - Increment
DEC - Decrement
NEG - Negate
CMP - Compare
AAA - ASCII adjust after addition
AAS - ASCII adjust after subtraction
AAM - ASCII adjust after multiply
AAD - ASCII adjust before division
DAA - Decimal adjust after addition
DAS - Decimal adjust after subtraction
CBW - Convert byte to word
CWD - Convert word to doubleword

Logical Instructions:

AND - Logical AND
OR - Logical OR
XOR - Logical XOR
NOT - Logical NOT
TEST - Test bits

Shift and Rotate Instructions:

SHL/SAL - Shift left
SHR - Shift right logical
SAR - Shift right arithmetic
ROL - Rotate left
ROR - Rotate right
RCL - Rotate through carry left
RCR - Rotate through carry right

String Instructions:

MOVS - Move string
CMPS - Compare string
SCAS - Scan string
LODS - Load string
STOS - Store string
REP - Repeat prefix
REPE/REPZ - Repeat while equal/zero
REPNE/REPNZ - Repeat while not equal/not zero

Control Transfer Instructions:

JMP - Unconditional jump
CALL - Call procedure
RET - Return from procedure
RETF - Return far
JZ/JE - Jump if zero/equal
JNZ/JNE - Jump if not zero/not equal
JS - Jump if sign
JNS - Jump if not sign
JO - Jump if overflow
JNO - Jump if not overflow
JB/JC - Jump if below/carry
JNB/JNC - Jump if not below/not carry
JBE - Jump if below or equal
JA - Jump if above
JL - Jump if less
JGE - Jump if greater or equal
JLE - Jump if less or equal
JG - Jump if greater
JP/JPE - Jump if parity/parity even
JNP/JPO - Jump if not parity/parity odd
JCXZ - Jump if CX is zero
LOOP - Loop
LOOPE/LOOPZ - Loop while equal/zero
LOOPNE/LOOPNZ - Loop while not equal/not zero

Processor Control Instructions:

CLC - Clear carry flag
STC - Set carry flag
CMC - Complement carry flag
CLD - Clear direction flag
STD - Set direction flag
CLI - Clear interrupt flag
STI - Set interrupt flag
HLT - Halt
WAIT - Wait
NOP - No operation

Interrupt Instructions:

INT - Software interrupt
INTO - Interrupt on overflow
IRET - Return from interrupt

I/O Instructions:

IN - Input from port
OUT - Output to port

BIOS Interrupts
INT 10h - Video Services:



AH
Function
Description



00h
Set Video Mode
AL = mode number


01h
Set Cursor Type
CH = start line, CL = end line


02h
Set Cursor Position
DH = row, DL = column, BH = page


03h
Get Cursor Position
BH = page → DH = row, DL = column


05h
Select Active Page
AL = page number


06h
Scroll Up
AL = lines, BH = attribute


07h
Scroll Down
AL = lines, BH = attribute


08h
Read Char/Attr
BH = page → AL = char, AH = attr


09h
Write Char/Attr
AL = char, BH = page, BL = attr, CX = count


0Ah
Write Char Only
AL = char, BH = page, CX = count


0Eh
Teletype Output
AL = character, BH = page


0Fh
Get Video Mode
→ AL = mode, AH = columns, BH = page


13h
Write String
ES:BP = string, CX = length


INT 13h - Disk Services:



AH
Function
Description



00h
Reset Disk
DL = drive


01h
Get Status
DL = drive → AH = status


02h
Read Sectors
AL = count, CH = cylinder, CL = sector, DH = head, DL = drive, ES:BX = buffer


03h
Write Sectors
AL = count, CH = cylinder, CL = sector, DH = head, DL = drive, ES:BX = buffer


04h
Verify Sectors
AL = count, CH = cylinder, CL = sector, DH = head, DL = drive


05h
Format Track
AL = count, CH = cylinder, DH = head, DL = drive, ES:BX = format data


08h
Get Drive Params
DL = drive → CH = cylinders, CL = sectors, DH = heads, DL = drives


15h
Get Disk Type
DL = drive → AH = type


16h
Detect Disk Change
DL = drive → AH = status


INT 16h - Keyboard Services:



AH
Function
Description



00h
Read Keystroke
→ AH = scan code, AL = ASCII


01h
Check Keystroke
→ ZF = 1 if no key, else AH = scan, AL = ASCII


02h
Get Shift Flags
→ AL = shift flags


03h
Set Typematic
AL = sub-function


05h
Store Keystroke
CH = scan code, CL = ASCII


10h
Extended Read
→ AH = scan code, AL = ASCII


11h
Extended Check
→ ZF = 1 if no key


12h
Extended Shift
→ AL = shift flags, AH = extended flags


INT 1Ah - Time Services:



AH
Function
Description



00h
Get System Time
→ CX:DX = tick count, AL = midnight flag


01h
Set System Time
CX:DX = tick count


02h
Get RTC Time
→ CH = hours, CL = minutes, DH = seconds, CF = error


03h
Set RTC Time
CH = hours, CL = minutes, DH = seconds


04h
Get RTC Date
→ CH = century, CL = year, DH = month, DL = day, CF = error


05h
Set RTC Date
CH = century, CL = year, DH = month, DL = day


06h
Set Alarm
CH = hours, CL = minutes, DH = seconds


07h
Reset Alarm
-


INT 33h - Mouse Services:



AX
Function
Description



00h
Reset Mouse
→ AX = FFFFh if present, BX = buttons


01h
Show Cursor
-


02h
Hide Cursor
-


03h
Get Position
→ BX = buttons, CX = X, DX = Y


04h
Set Position
CX = X, DX = Y


05h
Get Button Press
BX = button → AX = status, BX = count, CX = X, DX = Y


06h
Get Button Release
BX = button → AX = status, BX = count, CX = X, DX = Y


07h
Set Horizontal Limits
CX = min, DX = max


08h
Set Vertical Limits
CX = min, DX = max


DOS Interrupts
INT 20h - Program Terminate:
Terminates the current program and returns to DOS.
INT 21h - DOS Services:



AH
Function
Description



01h
Read Character
→ AL = character (with echo)


02h
Write Character
DL = character


03h
Auxiliary Input
→ AL = character


04h
Auxiliary Output
DL = character


05h
Printer Output
DL = character


06h
Direct Console I/O
DL = FFh for input, else output → AL = char (if input)


07h
Direct Input
→ AL = character (no echo)


08h
Input Without Echo
→ AL = character


09h
Write String
DS:DX = string address ($ terminated)


0Ah
Buffered Input
DS:DX = buffer


0Bh
Check Input Status
→ AL = FFh if char available


0Ch
Flush and Read
AL = input function


0Dh
Disk Reset
-


0Eh
Select Disk
DL = drive → AL = number of drives


19h
Get Current Drive
→ AL = drive (0=A, 1=B, ...)


1Ah
Set DTA Address
DS:DX = DTA address


25h
Set Interrupt Vector
AL = interrupt, DS:DX = handler


2Ah
Get Date
→ CX = year, DH = month, DL = day, AL = day of week


2Bh
Set Date
CX = year, DH = month, DL = day


2Ch
Get Time
→ CH = hours, CL = minutes, DH = seconds, DL = hundredths


2Dh
Set Time
CH = hours, CL = minutes, DH = seconds, DL = hundredths


2Eh
Set Verify Flag
AL = 0 (off) or 1 (on)


2Fh
Get DTA Address
→ ES:BX = DTA address


30h
Get DOS Version
→ AL = major, AH = minor, BX:CX = OEM


31h
TSR
AL = return code, DX = paragraphs


33h
Get/Set Break
AL = 0 (get) or 1 (set), DL = flag


34h
Get InDOS Flag
→ ES:BX = flag address


35h
Get Interrupt Vector
AL = interrupt → ES:BX = handler


36h
Get Disk Space
DL = drive → AX = sectors/cluster, BX = free clusters, CX = bytes/sector, DX = total clusters


38h
Get/Set Country
AL = function


39h
Create Directory
DS:DX = path


3Ah
Remove Directory
DS:DX = path


3Bh
Change Directory
DS:DX = path


3Ch
Create File
CX = attributes, DS:DX = filename → AX = handle


3Dh
Open File
AL = mode, DS:DX = filename → AX = handle


3Eh
Close File
BX = handle


3Fh
Read File
BX = handle, CX = bytes, DS:DX = buffer → AX = bytes read


40h
Write File
BX = handle, CX = bytes, DS:DX = buffer → AX = bytes written


41h
Delete File
DS:DX = filename


42h
Seek File
AL = method, BX = handle, CX:DX = offset → DX:AX = new position


43h
Get/Set Attributes
AL = function, DS:DX = filename, CX = attributes


44h
IOCTL
AL = function, BX = handle


45h
Duplicate Handle
BX = handle → AX = new handle


46h
Force Duplicate
BX = source, CX = target


47h
Get Current Dir
DL = drive, DS:SI = buffer


48h
Allocate Memory
BX = paragraphs → AX = segment


49h
Free Memory
ES = segment


4Ah
Resize Memory
BX = paragraphs, ES = segment


4Bh
Execute Program
AL = function, DS:DX = program, ES:BX = parameter block


4Ch
Exit Program
AL = return code


4Dh
Get Return Code
→ AX = return code


4Eh
Find First
CX = attributes, DS:DX = filespec


4Fh
Find Next
-


50h
Set PSP
BX = segment


51h
Get PSP
→ BX = segment


54h
Get Verify Flag
→ AL = flag


56h
Rename File
DS:DX = old name, ES:DI = new name


57h
Get/Set File Time
AL = function, BX = handle, CX = time, DX = date


59h
Get Extended Error
→ AX = error code, BH = class, BL = action, CH = locus


5Ah
Create Temp File
CX = attributes, DS:DX = path → AX = handle


5Bh
Create New File
CX = attributes, DS:DX = filename → AX = handle


5Ch
Lock/Unlock File
AL = function, BX = handle, CX:DX = offset, SI:DI = length


62h
Get PSP
→ BX = segment



Examples
Example 1: Hello World
hello.asm:
; Simple Hello World program
ORG 100h

    mov ah, 09h         ; Function: Display string
    mov dx, message     ; DS:DX points to message
    int 21h             ; Call DOS

    mov ah, 4Ch         ; Function: Exit program
    mov al, 0           ; Return code 0
    int 21h             ; Call DOS

message:
    db 'Hello, World!', 0Dh, 0Ah, '$'

Assemble and Run:
nasm -f bin -o hello.com hello.asm
./dos-emulator hello.com

Output:
Running program...
Hello, World!
Program exited with code: 0

Example 2: User Input
input.asm:
; Program that reads user input
ORG 100h

    ; Display prompt
    mov ah, 09h
    mov dx, prompt
    int 21h

    ; Read character
    mov ah, 01h
    int 21h
    mov bl, al          ; Save character in BL

    ; Display newline
    mov ah, 02h
    mov dl, 0Dh
    int 21h
    mov dl, 0Ah
    int 21h

    ; Display message
    mov ah, 09h
    mov dx, message
    int 21h

    ; Display the character
    mov ah, 02h
    mov dl, bl
    int 21h

    ; Exit
    mov ah, 4Ch
    int 21h

prompt:
    db 'Enter a character: $'
message:
    db 'You entered: $'

Example 3: Simple Calculator
calc.asm:
; Add two single-digit numbers
ORG 100h

    mov ax, 5           ; First number
    mov bx, 3           ; Second number
    add ax, bx          ; Add them

    ; Convert to ASCII
    add al, '0'

    ; Display result
    mov dl, al
    mov ah, 02h
    int 21h

    ; Display newline
    mov dl, 0Dh
    mov ah, 02h
    int 21h
    mov dl, 0Ah
    int 21h

    ; Exit
    mov ah, 4Ch
    int 21h

Output:
8

Example 4: Loop Example
loop.asm:
; Display numbers 1 to 5
ORG 100h

    mov cx, 5           ; Counter
    mov bl, '1'         ; Starting digit

print_loop:
    mov dl, bl
    mov ah, 02h
    int 21h

    ; Print space
    mov dl, ' '
    int 21h

    inc bl              ; Next digit
    loop print_loop

    ; Newline
    mov dl, 0Dh
    mov ah, 02h
    int 21h
    mov dl, 0Ah
    int 21h

    ; Exit
    mov ah, 4Ch
    int 21h

Output:
1 2 3 4 5

Example 5: File Creation
file.asm:
; Create a file and write to it
ORG 100h

    ; Create file
    mov ah, 3Ch         ; Function: Create file
    mov cx, 0           ; Normal attributes
    mov dx, filename
    int 21h
    jc error            ; Jump if error
    mov bx, ax          ; Save file handle

    ; Write to file
    mov ah, 40h         ; Function: Write to file
    mov cx, msg_len     ; Number of bytes
    mov dx, message
    int 21h
    jc error

    ; Close file
    mov ah, 3Eh         ; Function: Close file
    int 21h

    ; Display success message
    mov ah, 09h
    mov dx, success
    int 21h

    ; Exit
    mov ah, 4Ch
    int 21h

error:
    mov ah, 09h
    mov dx, err_msg
    int 21h
    mov ah, 4Ch
    mov al, 1
    int 21h

filename:
    db 'output.txt', 0
message:
    db 'Hello from file!', 0Dh, 0Ah
msg_len equ $ - message
success:
    db 'File created successfully!', 0Dh, 0Ah, '$'
err_msg:
    db 'Error creating file!', 0Dh, 0Ah, '$'

Example 6: String Operations
string.asm:
; String copy using MOVSB
ORG 100h

    ; Setup segments
    mov ax, cs
    mov ds, ax
    mov es, ax

    ; Setup pointers
    mov si, source
    mov di, dest
    mov cx, src_len
    cld                 ; Clear direction flag

    ; Copy string
    rep movsb

    ; Display destination
    mov ah, 09h
    mov dx, dest
    int 21h

    ; Exit
    mov ah, 4Ch
    int 21h

source:
    db 'Hello String!$'
src_len equ $ - source
dest:
    times 20 db 0

Example 7: Conditional Jumps
compare.asm:
; Compare two numbers
ORG 100h

    mov ax, 10
    mov bx, 5

    cmp ax, bx
    jg greater          ; Jump if AX > BX
    je equal            ; Jump if AX = BX
    jl less             ; Jump if AX < BX

greater:
    mov dx, msg_greater
    jmp display

equal:
    mov dx, msg_equal
    jmp display

less:
    mov dx, msg_less

display:
    mov ah, 09h
    int 21h

    mov ah, 4Ch
    int 21h

msg_greater:
    db 'First is greater', 0Dh, 0Ah, '$'
msg_equal:
    db 'Numbers are equal', 0Dh, 0Ah, '$'
msg_less:
    db 'First is less', 0Dh, 0Ah, '$'

Example 8: Subroutine Call
subroutine.asm:
; Using CALL and RET
ORG 100h

    call print_hello
    call print_hello
    call print_hello

    mov ah, 4Ch
    int 21h

print_hello:
    mov ah, 09h
    mov dx, message
    int 21h
    ret

message:
    db 'Hello!', 0Dh, 0Ah, '$'

Output:
Hello!
Hello!
Hello!


Troubleshooting
Common Issues and Solutions
Issue 1: File Not Found
Symptom:
A:\> RUN program.com
Error: open program.com: no such file or directory

Solutions:

Check if file exists:
A:\> DIR


Verify current directory:
A:\> CD


Use full path:
A:\> RUN C:\programs\test.com


Check file permissions:
ls -l program.com
chmod +r program.com



Issue 2: Program Hangs
Symptom:
Program appears to hang or loop infinitely.
Solutions:

Enable debug mode to see what's happening:
A:\> DEBUG
A:\> RUN program.com


Use step mode to trace execution:
A:\> STEP
A:\> RUN program.com


Set breakpoints to find infinite loops:
A:\> BREAK 1100
A:\> RUN program.com


Check for:

Missing exit code (INT 21h/4Ch)
Infinite loops in code
Incorrect jump addresses



Issue 3: Incorrect Output
Symptom:
Program produces unexpected output or crashes.
Solutions:

Examine registers during execution:
A:\> STEP
A:\> RUN program.com
(press 'r' at each step)


Check memory contents:
A:\> DUMP 100 256


Disassemble code to verify:
A:\> DISASM 100 50


Common causes:

Wrong register values
Incorrect memory addresses
Missing $ terminator in strings
Wrong interrupt numbers



Issue 4: Assembly Errors
Symptom:
NASM reports errors during assembly.
Common Errors:

Missing ORG directive:
; Wrong
mov ah, 09h

; Correct
ORG 100h
mov ah, 09h


Missing $ terminator:
; Wrong
message db 'Hello'

; Correct
message db 'Hello$'


Incorrect syntax:
; Wrong
mov [100], ax

; Correct
mov word [100], ax



Issue 5: Memory Access Errors
Symptom:
Program crashes or produces garbage output.
Solutions:

Verify segment registers:
A:\> REGS


Check memory addresses are valid:
A:\> DUMP address


Ensure proper segment setup:
mov ax, cs
mov ds, ax
mov es, ax



Issue 6: Stack Overflow
Symptom:
Program crashes after many PUSH operations or recursive calls.
Solutions:

Check stack pointer:
A:\> REGS
(look at SP value)


View stack contents:
A:\> STACK


Ensure balanced PUSH/POP:
push ax
push bx
; ... code ...
pop bx
pop ax



Issue 7: Interrupt Not Working
Symptom:
INT 21h or other interrupts don't produce expected results.
Solutions:

Verify function number in AH:
mov ah, 09h    ; Correct function number
int 21h


Check required registers:
; For INT 21h/09h
mov ah, 09h
mov dx, message    ; DS:DX must point to string
int 21h


Enable debug mode to see interrupt calls:
A:\> DEBUG
A:\> RUN program.com



Debugging Checklist
When troubleshooting a program:

 Verify file exists and is readable
 Check ORG 100h directive is present
 Ensure strings end with '$'
 Verify correct interrupt numbers
 Check all registers are initialized
 Ensure balanced PUSH/POP operations
 Verify segment registers are set correctly
 Check for proper program exit (INT 21h/4Ch)
 Use DEBUG mode to trace execution
 Examine memory with DUMP
 Disassemble code with DISASM


FAQ
General Questions
Q: What is MS-DOS Emulator?
A: MS-DOS Emulator is a software implementation of an Intel 8086 CPU with BIOS and DOS support, allowing you to run classic DOS programs on modern systems.
Q: Can I run any DOS program?
A: The emulator supports many simple DOS programs, especially those using standard INT 21h functions. Complex programs using advanced features, graphics modes, or direct hardware access may not work.
Q: Is this a complete DOS implementation?
A: The emulator implements the most commonly used DOS and BIOS functions. Some advanced or rarely-used features may not be available.
Compatibility Questions
Q: Can I run MS-DOS games?
A: Simple text-based games may work. Graphics-intensive games requiring VGA graphics, sound cards, or precise timing will not work.
Q: What file formats are supported?
A: COM files are fully supported. EXE files have partial support. COM format is recommended for best compatibility.
Q: Can I run Windows programs?
A: No. This emulator only supports DOS programs, not Windows applications.
Q: Will my old DOS programs work?
A: Programs that use standard DOS services (INT 21h) and don't require special hardware will likely work. Test your specific program to be sure.
Programming Questions
Q: How do I create COM files?
A: Use an assembler like NASM:
nasm -f bin -o program.com program.asm

Q: What assembler should I use?
A: NASM (Netwide Assembler) is recommended. It's free, cross-platform, and well-documented.
Q: Can I use C or other languages?
A: You can use any language that can produce DOS COM files. However, assembly language provides the most control and compatibility.
Q: Where can I learn 8086 assembly?
A: Resources include:

"The Art of Assembly Language" by Randall Hyde
"Assembly Language Step-by-Step" by Jeff Duntemann
Online tutorials for 8086 assembly
NASM documentation

Technical Questions
Q: How accurate is the emulation?
A: The emulator provides functional compatibility for most common operations. Timing and some hardware-specific features differ from real hardware.
Q: What's the maximum program size?
A: COM files can be up to approximately 64KB (65,280 bytes).
Q: How much memory is available?
A: The emulator provides 1MB of addressable memory, with 640KB of conventional memory available to programs.
Q: Are interrupts fully implemented?
A: Common BIOS (INT 10h, 13h, 16h, 1Ah) and DOS (INT 21h) interrupts are implemented. Some advanced functions may not be available.
Usage Questions
Q: How do I exit a running program?
A: Programs should exit using INT 21h function 4Ch. If a program hangs, you may need to use Ctrl+C (terminal dependent) or close the terminal.
Q: Can I access my real file system?
A: Yes. The emulator uses your current directory as drive A:. You can read and write real files.
Q: How do I debug my programs?
A: Use the built-in debugging features:

DEBUG mode for detailed execution info
STEP mode for instruction-by-instruction execution
BREAK to set breakpoints
DUMP to examine memory
DISASM to disassemble code

Q: Can I run multiple programs?
A: Programs run sequentially. After one program exits, you can run another.
Performance Questions
Q: Why is execution slow?
A: The emulator interprets each instruction. For better performance:

Disable DEBUG mode when not needed
Avoid excessive breakpoints
Use efficient code

Q: How can I measure performance?
A: Use the STATS command to see execution statistics including instructions per second.
Q: Is there a speed limit?
A: The emulator runs as fast as your system allows. There's no artificial speed limiting.
Error Messages
Q: What does "Bad command or file name" mean?
A: The command or file you entered doesn't exist or isn't recognized. Check spelling and use DIR to list files.
Q: What does "File not found" mean?
A: The specified file doesn't exist in the current directory. Use DIR to verify the filename.
Q: What does "Invalid directory" mean?
A: The directory path you specified doesn't exist. Use DIR to see available directories.
Q: What does "Unhandled interrupt" mean?
A: Your program called an interrupt function that isn't implemented in the emulator. Check your code or use a different approach.

Appendices
Appendix A: Complete Interrupt Reference
INT 10h - Video Services
Function 00h - Set Video Mode
Input:  AH = 00h
        AL = video mode
Output: None

Function 02h - Set Cursor Position
Input:  AH = 02h
        BH = page number
        DH = row
        DL = column
Output: None

Function 03h - Get Cursor Position
Input:  AH = 03h
        BH = page number
Output: CH = cursor start line
        CL = cursor end line
        DH = row
        DL = column

Function 06h - Scroll Up Window
Input:  AH = 06h
        AL = number of lines (0 = clear window)
        BH = attribute for blank lines
        CH = top row
        CL = left column
        DH = bottom row
        DL = right column
Output: None

Function 0Eh - Teletype Output
Input:  AH = 0Eh
        AL = character
        BH = page number
        BL = foreground color (graphics mode)
Output: None

Function 0Fh - Get Current Video Mode
Input:  AH = 0Fh
Output: AL = video mode
        AH = number of columns
        BH = active page

INT 13h - Disk Services
Function 00h - Reset Disk System
Input:  AH = 00h
        DL = drive number
Output: CF = 0 if successful
        AH = status

Function 02h - Read Sectors
Input:  AH = 02h
        AL = number of sectors
        CH = cylinder
        CL = sector
        DH = head
        DL = drive
        ES:BX = buffer address
Output: CF = 0 if successful
        AH = status
        AL = sectors read

Function 08h - Get Drive Parameters
Input:  AH = 08h
        DL = drive number
Output: CF = 0 if successful
        CH = maximum cylinder number
        CL = maximum sector number
        DH = maximum head number
        DL = number of drives

INT 16h - Keyboard Services
Function 00h - Read Keystroke
Input:  AH = 00h
Output: AH = BIOS scan code
        AL = ASCII character

Function 01h - Check for Keystroke
Input:  AH = 01h
Output: ZF = 1 if no keystroke available
        ZF = 0 if keystroke available
        AH = BIOS scan code (if available)
        AL = ASCII character (if available)

Function 02h - Get Shift Flags
Input:  AH = 02h
Output: AL = shift flags
        Bit 0: Right Shift pressed
        Bit 1: Left Shift pressed
        Bit 2: Ctrl pressed
        Bit 3: Alt pressed
        Bit 4: Scroll Lock on
        Bit 5: Num Lock on
        Bit 6: Caps Lock on
        Bit 7: Insert mode on

INT 1Ah - Time Services
Function 00h - Get System Time
Input:  AH = 00h
Output: CX = high word of tick count
        DX = low word of tick count
        AL = midnight flag

Function 02h - Get Real-Time Clock Time
Input:  AH = 02h
Output: CF = 0 if successful
        CH = hours (BCD)
        CL = minutes (BCD)
        DH = seconds (BCD)

Function 04h - Get Real-Time Clock Date
Input:  AH = 04h
Output: CF = 0 if successful
        CH = century (BCD)
        CL = year (BCD)
        DH = month (BCD)
        DL = day (BCD)

INT 21h - DOS Services (Selected Functions)
Function 01h - Character Input with Echo
Input:  AH = 01h
Output: AL = character read

Function 02h - Character Output
Input:  AH = 02h
        DL = character to output
Output: None

Function 09h - Display String
Input:  AH = 09h
        DS:DX = address of string ($ terminated)
Output: None

Function 0Ah - Buffered Input
Input:  AH = 0Ah
        DS:DX = address of input buffer
        Buffer format:
          Byte 0: Maximum characters
          Byte 1: Actual characters read (returned)
          Byte 2+: Input characters
Output: Buffer filled with input

Function 3Ch - Create File
Input:  AH = 3Ch
        CX = file attributes
        DS:DX = filename address (ASCIIZ)
Output: CF = 0 if successful
        AX = file handle (if successful)
        AX = error code (if failed)

Function 3Dh - Open File
Input:  AH = 3Dh
        AL = access mode
          0 = read only
          1 = write only
          2 = read/write
        DS:DX = filename address (ASCIIZ)
Output: CF = 0 if successful
        AX = file handle (if successful)
        AX = error code (if failed)

Function 3Eh - Close File
Input:  AH = 3Eh
        BX = file handle
Output: CF = 0 if successful
        AX = error code (if failed)

Function 3Fh - Read from File
Input:  AH = 3Fh
        BX = file handle
        CX = number of bytes to read
        DS:DX = buffer address
Output: CF = 0 if successful
        AX = number of bytes read
        AX = error code (if failed)

Function 40h - Write to File
Input:  AH = 40h
        BX = file handle
        CX = number of bytes to write
        DS:DX = buffer address
Output: CF = 0 if successful
        AX = number of bytes written
        AX = error code (if failed)

Function 4Ch - Exit Program
Input:  AH = 4Ch
        AL = return code
Output: Does not return

Appendix B: ASCII Table
Dec Hex Char | Dec Hex Char | Dec Hex Char | Dec Hex Char
  0  00 NUL |  32  20 SP  |  64  40  @  |  96  60  `
  1  01 SOH |  33  21  !  |  65  41  A  |  97  61  a
  2  02 STX |  34  22  "  |  66  42  B  |  98  62  b
  3  03 ETX |  35  23  #  |  67  43  C  |  99  63  c
  4  04 EOT |  36  24  $  |  68  44  D  | 100  64  d
  5  05 ENQ |  37  25  %  |  69  45  E  | 101  65  e
  6  06 ACK |  38  26  &amp;  |  70  46  F  | 102  66  f
  7  07 BEL |  39  27  '  |  71  47  G  | 103  67  g
  8  08 BS  |  40  28  (  |  72  48  H  | 104  68  h
  9  09 TAB |  41  29  )  |  73  49  I  | 105  69  i
 10  0A LF  |  42  2A  *  |  74  4A  J  | 106  6A  j
 11  0B VT  |  43  2B  +  |  75  4B  K  | 107  6B  k
 12  0C FF  |  44  2C  ,  |  76  4C  L  | 108  6C  l
 13  0D CR  |  45  2D  -  |  77  4D  M  | 109  6D  m
 14  0E SO  |  46  2E  .  |  78  4E  N  | 110  6E  n
 15  0F SI  |  47  2F  /  |  79  4F  O  | 111  6F  o
 16  10 DLE |  48  30  0  |  80  50  P  | 112  70  p
 17  11 DC1 |  49  31  1  |  81  51  Q  | 113  71  q
 18  12 DC2 |  50  32  2  |  82  52  R  | 114  72  r
 19  13 DC3 |  51  33  3  |  83  53  S  | 115  73  s
 20  14 DC4 |  52  34  4  |  84  54  T  | 116  74  t
 21  15 NAK |  53  35  5  |  85  55  U  | 117  75  u
 22  16 SYN |  54  36  6  |  86  56  V  | 118  76  v
 23  17 ETB |  55  37  7  |  87  57  W  | 119  77  w
 24  18 CAN |  56  38  8  |  88  58  X  | 120  78  x
 25  19 EM  |  57  39  9  |  89  59  Y  | 121  79  y
 26  1A SUB |  58  3A  :  |  90  5A  Z  | 122  7A  z
 27  1B ESC |  59  3B  ;  |  91  5B  [  | 123  7B  {
 28  1C FS  |  60  3C  <  |  92  5C  \  | 124  7C  |
 29  1D GS  |  61  3D  =  |  93  5D  ]  | 125  7D  }
 30  1E RS  |  62  3E  >  |  94  5E  ^  | 126  7E  ~
 31  1F US  |  63  3F  ?  |  95  5F  _  | 127  7F DEL

Appendix C: Error Codes
DOS Error Codes (INT 21h):
1   Invalid function number
2   File not found
3   Path not found
4   Too many open files
5   Access denied
6   Invalid handle
7   Memory control blocks destroyed
8   Insufficient memory
9   Invalid memory block address
10  Invalid environment
11  Invalid format
12  Invalid access code
13  Invalid data
14  Reserved
15  Invalid drive
16  Attempt to remove current directory
17  Not same device
18  No more files

Appendix D: Keyboard Scan Codes
Common Scan Codes:
Key         Scan Code
ESC         01h
1           02h
2           03h
...
9           0Ah
0           0Bh
-           0Ch
=           0Dh
Backspace   0Eh
Tab         0Fh
Q           10h
W           11h
...
Enter       1Ch
Ctrl        1Dh
A           1Eh
...
Space       39h
F1          3Bh
F2          3Ch
...
F10         44h

Appendix E: File Attributes
File Attribute Bits:
Bit  Hex  Description
0    01h  Read-only
1    02h  Hidden
2    04h  System
3    08h  Volume label
4    10h  Directory
5    20h  Archive
6    40h  Reserved
7    80h  Reserved

Appendix F: Video Modes
Common Video Modes:
Mode  Type        Resolution  Colors
00h   Text        40x25       16 (gray)
01h   Text        40x25       16
02h   Text        80x25       16 (gray)
03h   Text        80x25       16
04h   Graphics    320x200     4
05h   Graphics    320x200     4 (gray)
06h   Graphics    640x200     2
07h   Text        80x25       Mono
0Dh   Graphics    320x200     16
0Eh   Graphics    640x200     16
0Fh   Graphics    640x350     Mono
10h   Graphics    640x350     16
11h   Graphics    640x480     2
12h   Graphics    640x480     16
13h   Graphics    320x200     256

Note: This emulator primarily supports text mode 03h (80x25, 16 colors).

End of User Manual
MS-DOS Emulator v4.0 - Complete Implementation
For updates and support, please refer to the project repository.
