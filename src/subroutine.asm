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