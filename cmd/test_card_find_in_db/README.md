## PACKAGE

This package is for testing plain unciphered string of card UID input and check if that card UID exists in database.

### EXAMPLE USAGE

```
./<binary name> -cardUID "<selected card uid>"
```

### CARD UID FORMAT

Input should be card UID.
If card UID is [0f 48 5e 69] then the string should be "0f 48 5e 69".
Any amount of spaces are valid.
Eg. "0f485e69" or "   0f         48   5e69   " is valid.
