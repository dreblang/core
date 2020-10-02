#!/bin/bash
echo "Lexer tests..."
go test -cover -coverprofile=lexer.out ./lexer/... # && go tool cover -html=lexer.out

echo "Parser tests..."
go test -cover -coverprofile=parser.out ./parser/... && go tool cover -html=parser.out
