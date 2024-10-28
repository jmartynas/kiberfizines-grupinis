## AUTHENTICATION SERVER

This is authentication server for office access with RFID card.

### LAUNCHING PROGRAM

Launching binaries should be enough.

```
./<binary name>
```

### Request example

Without response header

curl localhost:8080/ -X POST -d '{"UUID":"0192c9f5-02fc-7eb1-9e72-fdf12acf481e","content":"9f704f5456828ad6f04756ae7a85c5ab"}'

```
$ curl localhost:8080/ -X POST -d '{"UUID":"0192c9f5-02fc-7eb1-9e72-fdf12acf481e","content":"9f704f5456828ad6f04756ae7a85c5ab"}'
vardenis pavardenis entered
vardenis pavardenis
```

With response header

```
$ curl localhost:8080/ -i -d '{"UUID":"0192c9f5-02fc-7eb1-9e72-fdf12acf481e","content":"9f704f5456828ad6f04756ae7a85c5ab"}'
vardenis pavardenis entered
HTTP/1.1 200 OK
Date: Mon, 28 Oct 2024 17:35:10 GMT
Content-Length: 19
Content-Type: text/plain; charset=utf-8

vardenis pavardenis
```

### CARD UID FORMAT BEFORE ENCRYPTING

When hashing a string it should be of a string type card UID.
If card UID is [0f 48 5e 69] then the string should be "0f 48 5e 69".
Any amount of spaces are valid.
Eg. "0f485e69" or " 0f 48 5e69 " is valid.

### HARD CODED VALUES

`Device UUID`: `0192c9f5-02fc-7eb1-9e72-fdf12acf481e`

`key`: `{0x2b,0x7e,0x15,0x16,0x28,0xae,0xd2,0xa6,0xab,0xf7,0x97,0x99,0x89,0xcf,0xab,0x12}`

`iv`: `{0x2b,0x7e,0x15,0x16,0x28,0xae,0xd2,0xa6,0xab,0xf7,0x97,0x99,0x89,0xcf,0xab,0x12}`
