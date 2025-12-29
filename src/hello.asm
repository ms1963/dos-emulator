org 0x100

mov ah, 0x09
mov dx, msg
int 0x21

mov ah, 0x4C
mov al, 0
int 0x21

msg db 'Hello from DOS!$'
