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
