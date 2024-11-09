## Abstract

Persistence: writing data to durable storage.
For ensuring persistence in dicedb, we use Append only file, where we append every write operations performed on the file.

Let's say we have a key named `counter` and we increment 10000 times. There would be 10000 fields of `counter` key in the AOF file. This could lead to increase the time, it takes to reconstruct the original dataset, in case of crash or power failure. 

For that redis use `BGREWRITEAOF` command. It creates a new child process via `fork()`, and schedules creation of a new AOF in the background. So, at this point of time, we can have 2 versions of Append only file: base file and incremental file. Main file may contain duplicate instances of keys and is the current active AOF file to which all the write operations are being added. Till all the unique keys in the main file are not updated in the incremental file with their latest value, redis keep executing the recreation child process in the background. Once, the recreation is completed, it will assign the incremental file as new base file and all the new write operations will be appended to this file.

## Commit 9 (249cb3acc4f22888dfc81a9bd72fab3aca083302)
Current implemententation just dumps all the keys in the hash table to append only file via `BGREWRITEAOF` command.



## Reference
https://redis.io/docs/latest/operate/oss_and_stack/management/persistence