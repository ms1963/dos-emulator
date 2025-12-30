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