## Abstract 

### Commit 6(1ebca1b96352cd4b431c4acb19268d6b219de408): DEL, EXPIRE command implementation

#### evalDEL(): package core

- DEL command to delete the keys from the hash table
- it deletes the key that exists, and if there is a key that does not exist, it ignores them.
- it counts the no. of keys that are deleted from the database, and returns integer reply.

#### Del(): store.go, package core
- checks whether key exists or not
- if yes, deletes the key using `delete()` method and return true
- else return false

#### Get() store.go, package core
- checks if key has expiry date set or not.
- if yes, check if key is expired or not.
- if yes, delete the key and return nil.
- otherwise return object

#### evalEXPIRE(): eval.go, package core

- EXPIRE command with 2 arguments: [key] [seconds]
- if no. of arguments is less than 2, return error
- get key from the 0th index
- parse seconds to 64 bit base 10 integer, and if error occurred due to integer not being found or out of range, return error.
- get key from the hash table
- if key does not exist or has expired, that means timeout will not be set, thus return 0 integer in RESP encoded format(`:0\r\n`)
- else add new timeout to obj's `ExpiredAt` field and return RESP encoded integer 1 to the client.

#### expire.go: package core
- It mentions Sampling Approach which Redis actually uses for efficiently expiring of keys from the hash table.
- It initially assumes the sample size to be 20 and iterations in hash table to be randomized.
- Sample signifies total no. of keys with expiry duration set, i.e. ExpiredAt != -1.
- If fraction of keys expired out of the total keys with expiry time set is greater than 25%(i.e. very large no. of keys have expired, and keep iterating till fraction reduces further), then keep iterating till fraction becomes < 25%(i.e. most of the expired keys have been deleted).

#### async_tcp.go: package server
- Setup a cron job to automatically delete expired keys every second, via sampling approach.
- after 1 second(which is our cron frequency), delete the expired keys.
- update the last cron execution time to current time using `time.Now()`.