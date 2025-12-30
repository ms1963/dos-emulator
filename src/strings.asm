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