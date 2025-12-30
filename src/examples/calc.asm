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
