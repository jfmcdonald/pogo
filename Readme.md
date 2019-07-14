Pogo the IO Clown
===================

Pogo the IO Clown is inteded as a network disk io 
test for HPCC environments. It is a simulation of a
extreamly porely written user job. Pogo should be run 
as an array job spawing many threads accross as many 
hosts as possible. The idea is to spread the IO as 
wide as possible to create a "worst case sanareo" for
any cluster connected storage.

Whats going on here?
--------------------
Ok so what are we actually doing here?
Each pogo thread makes a choice to create, read,
update or delete a file. It then randomly chooses a 
file to act on or generates a new random file name to 
create. File creation creates an empty file and an
index Key:Value, a read pulls in all lines and does a 
char count, an update adds a random 8 bit string and 
CR/LF, and delete removes the file and associated index.

Pogo relys on a redis key value store, created
files are stored with a name:path in redis so we never
have to index the working directory. This is to reduce
the advantage of storage systems with extreamly fast 
metadata servies. Each key pair is set with a TTL
(default is 1 hour) so anything that is not deleted 
gets auto cleaned up.

Usage?
------

Usage of ./pogo:
  -./ string
        Path wher run time files will be generated (default "/tmp")
  -K/V ttl int
        Index Key/Value store default key TTL (default 60)
  -count int
        Total number of files to generate (default 100)


Credit
-------
Created by Jason McDonald
For Harvard Medical School
Created 12 July 2019
Updated 12 July 2019


