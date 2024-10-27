## AUTHENTICATION SERVER

This is authentication server for office access with RFID card.

### LAUNCHING PROGRAM

Launching binaries should be enough.

```
./<binary name>
```

### Request example

Without response header

```
curl localhost:8080/ -X POST -d '{"UUID":"0192c9f5-02fc-7eb1-9e72-fdf12acf481e","content":"37541266f93eadb996a826dd57cbcc540acf44fe19e4a2d53f0a233bd30c4a67723f7eb4402662be0c7312da2f41dea4f4e4c4a67dad9c08005e2db8f8c3045232cef57cae172e45b52ba9a5235bec53fa9212c7ceffac5fc96f0abb9e0ccbe5a7741c8a29829cdfb7c98a0bbe2a623b78f8b3abe875b1956745684f80531515b9e092115a05a06fb430250fe486bef2f56e77745aeb22cd9abd9440b764cd1fc5e9866da9109fdaaca64f0d5ea4150c16f2024b70c0bf919441cd2ba29bd9880ce5508ff5ccce6af1397c2b595d6bf076ff91f97154666e840a01cab613f23e7a73f049e721a7c9e4168ac9845b88bb0f77733b74947112d90c6f0c72e28280"}'
```

With response header

```
curl localhost:8080/ -i -d '{"UUID":"0192c9f5-02fc-7eb1-9e72-fdf12acf481e","content":"37541266f93eadb996a826dd57cbcc540acf44fe19e4a2d53f0a233bd30c4a67723f7eb4402662be0c7312da2f41dea4f4e4c4a67dad9c08005e2db8f8c3045232cef57cae172e45b52ba9a5235bec53fa9212c7ceffac5fc96f0abb9e0ccbe5a7741c8a29829cdfb7c98a0bbe2a623b78f8b3abe875b1956745684f80531515b9e092115a05a06fb430250fe486bef2f56e77745aeb22cd9abd9440b764cd1fc5e9866da9109fdaaca64f0d5ea4150c16f2024b70c0bf919441cd2ba29bd9880ce5508ff5ccce6af1397c2b595d6bf076ff91f97154666e840a01cab613f23e7a73f049e721a7c9e4168ac9845b88bb0f77733b74947112d90c6f0c72e28280"}'
```

### CARD UID FORMAT BEFORE ENCRYPTING

When hashing a string it should be of a string type card UID.
If card UID is [0f 48 5e 69] then the string should be "0f 48 5e 69".
Any amount of spaces are valid.
Eg. "0f485e69" or "   0f         48   5e69   " is valid.
