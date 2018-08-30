* A Go-specific Text-to-Speech converter

This is a proof-of-concept program to do text-to-speech using the Go AST parser. It
attempts to do an English-language rendering of what the code is doing, rather than
just reading individual words and characters as a regular TTS engine would.

It currently relies on the Mac OSX "say" program, so it won't run on another platform
without reworking the speak function.

The executable is in src/go-to-speech/cmd/speaker/main.go. You can give it a -q option to disable
the speaking if you are just debugging the language processing. Otherwise, just specify Go files
on the command-line and it will read out each one.