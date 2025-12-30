
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