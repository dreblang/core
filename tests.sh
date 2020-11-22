#!/bin/bash
echo "Lexer tests..."
go test -cover -coverprofile=lexer.out ./lexer/... # && go tool cover -html=lexer.out

echo "Parser tests..."
go test -cover -coverprofile=parser.out ./parser/... # && go tool cover -html=parser.out

echo "AST tests..."
go test -cover -coverprofile=ast.out ./ast/... # && go tool cover -html=ast.out

echo "Object tests..."
go test -cover -coverprofile=object.out ./object/... # && go tool cover -html=object.out
