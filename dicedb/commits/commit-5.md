## Abstract 

### Commit(7ea411bac39830fb2a073d194f34fe8576434302): GET, SET, TTL redis like commands implementation 

#### evalSET(): package core
- sets key-value pair in hash table
- Supports the "EX" option to set expiration time in seconds
- returns "OK" on success
- error handling for wrong no. of arguments

```bash
SET <key> <value> EX <expiration-time-in-seconds>
```

#### evalGET(): package core
- retrieves value for given key
- Return null(serialized via RESP protocol) if:
    - key does not exist
    - key has expired
- return key's value if key exist and has not expired

#### evalTTL(): package core
- returns the remaining time to live of key in seconds
- returns specific values:
    - (-2): key does not exist
        - (key, value) pair not present in hash table
        - key has expired
    - (-1): key exists, but with no expiration
- returns remaining time(non -ve) in seconds for keys with expiration
