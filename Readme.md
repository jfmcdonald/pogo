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
gets auto cleaned up. You will still need to manually  
clean up any temp files not deleted. This is also the  
case if you kill your job before it finishes.

Usage?
------

Usage of ./pogo:<br>
&nbsp;&nbsp;  -count uint<br>
&nbsp;&nbsp;&nbsp;&nbsp;        Total number of files to generate (default 10)<br>
&nbsp;&nbsp;  -db int<br>
&nbsp;&nbsp;&nbsp;&nbsp;        redis db id you want to store keys in<br>
&nbsp;&nbsp;  -dbhost string<br>
&nbsp;&nbsp;&nbsp;&nbsp;       Hostname of the network redis server (default "localhost")<br>
&nbsp;&nbsp;  -logfile string<br>
&nbsp;&nbsp;&nbsp;&nbsp;        location where you want to log message (default "/dev/null")<br>
&nbsp;&nbsp; -pass string<br>
&nbsp;&nbsp;&nbsp;&nbsp;        Password for redis db<br>
&nbsp;&nbsp;  -path string<br>
&nbsp;&nbsp;&nbsp;&nbsp;        Path wher run time files will be generated (default "/tmp/pogo")<br>
&nbsp;&nbsp;  -ttl int<br>
&nbsp;&nbsp;&nbsp;&nbsp;        Index Key/Value store default key TTL (default 60) <br>
&nbsp;&nbsp;&nbsp;&nbsp;        Index Key/Value store default key TTL (default 60)<br>
&nbsp;&nbsp;&nbsp;&nbsp;	Make sure you update this for how ever long you want your job to run<br>
&nbsp;&nbsp;&nbsp;&nbsp;	If indexes get deleted before the files are removed they won't be<br>
&nbsp;&nbsp;&nbsp;&nbsp;	cleaned up after the run.<br>

Credit
-------
Created by Jason McDonald  
For Harvard Medical School  
Created 12 July 2019  
