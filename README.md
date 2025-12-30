<img width="1024" height="1024" alt="image" src="https://github.com/user-attachments/assets/4cea2842-8225-4612-a743-65acd386520e" />


# Go back to the past with the dos-emulator

This application implements an emulation of the 8086-CPU family and provides a shell with MS-DOS emulation.

The whole application was written in Go to obtain experience with emulators. Although the code has been checked for errors several times, some errors will certainly pop up.
You may write .com or .exe runnable applications yourself using the NASM compiler (www.nasm.us), or provide existing executables. dos-emulator will load them into memory with the required layout (.exe or .com layout) and execute them. 

On errors, don't hesitate to contact me or - better - to report an issue.



## To build the dos-emulator use the Go compiler:

go build -o dos-emulator dos.go



## For cross-compilation use:

### Build for Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o dos-emulator.exe

### Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o dos-emulator-mac-intel

### Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dos-emulator-mac-arm

### Build for Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o dos-emulator-linux

**Enjoy the dos-emulator time machine!**
