Pogo the IO Clown
===================

Pogo the IO Clown is intended as a network disk io  
test for HPCC environments. It is a simulation of a  
extremely poorly written user job. Pogo should be run  
as an array job spawning many threads across as many  
hosts as possible. The idea is to spread the IO as  
wide as possible to create a "worst case scenario" for  
any cluster connected storage.  

Warnings!
=========

Don't run pogo against production file systems!
No Really, Don't!
Pogo the IO Clown generates rather small files and  
expands them slowly if at all. That said, Pogo is able  
to create a file-system crushing number of files depending  
on configuration. The point of Pogo is to break stuff.  
Don't break stuff your boss is going to get mad about.  
You have been warned. If Mis-use of Pogo gets you fired or  
arrested it's so not my fault.

Whats going on here?
--------------------
Ok so what are we actually doing here?  
Each pogo thread makes a choice to create, read,  
update or delete a file. It then randomly chooses a  
file to act on or generates a new random file name to  
create. File creation creates an empty file and an  
index Key:Value, a read pulls in all lines and does a  
char count, an update adds a random 8 bit string and  
and delete removes the file and associated index.  

Pogo relies on a Redis key value store, created  
files are stored with a name:path in Redis so we never  
have to index the working directory. This is to reduce  
the advantage of storage systems with extremely fast  
meta-data services. Each key pair is set with a TTL  
(default is 1 hour) so anything that is not deleted  
gets auto cleaned up.  

Usage?
------

Usage of ./pogo:  
  -./ string  
        Path where run time files will be generated (default "/tmp")  
  -k/v ttl int    
        Index Key/Value store default key TTL (default 60)  
  -count int  
        Total number of files to generate (default 100)  


Credit
-------
Created by Jason McDonald  
For Harvard Medical School  
Created 12 July 2019  
