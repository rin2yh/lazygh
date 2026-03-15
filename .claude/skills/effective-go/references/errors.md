# Errors

## Return errors

Return errors as the last return value.

func Read() ([]byte, error)

## Do not panic for normal errors

Use panic only for programmer errors or unrecoverable states.

## Error messages

Error strings should:

- start lowercase
- avoid punctuation
- describe the failure

Bad

"File Not Found."

Good

"file not found"
