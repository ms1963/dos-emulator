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
