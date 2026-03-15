# Naming

## General rules

- Use short, clear names
- Avoid redundant names (stutter)
- Prefer common Go conventions

Bad

httpserver.Server

Good

http.Server

## Package names

- lowercase
- single word if possible
- no underscores or mixed caps

Bad

string_utils

Good

strings

## Variable names

Short names are acceptable in small scopes.

Common conventions:

i, j → loop index
err → error value
ctx → context.Context

Avoid overly descriptive names.

Bad

userInformationMap

Good

users
