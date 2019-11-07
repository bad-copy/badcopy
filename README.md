# badcopy
badcopy is a command line rescue tool. You can use it to copy corrupted file or directory. 

Mac(64), Linux(32/64) and Windows(32/64) are supported. 

Source file is written in GOLANG and executable files are packed by UPX.

Usage:

```
badcopy Options
```

Options:

```
  -b int
    	Block size: 0-512, 1-1K,  2-2K,  3-4K,   4-8K,
    	            5-16K, 6-32K, 7-64K, 8-128K, 9-256K
    	            (default 3)
  -c	Check timestamp. File with different timestamp will be treated 
    	as different file, which will be overwriten.
  -f	Force overwrite.
  -i string
    	*Input directory or file
  -o string
    	*Output directory
  -r int
    	Retry times: 0 - 9
  -s	Skip if failed to read file data. Only readed data will be stored.
    	Left data will be skipped.
```

Example:

```
./dist/badcopy-darwin64-2019a -i ~/go/src -o ~/go/test
======================================================
badcopy 2019a is created by luomao2000@tom.com
You can use it to copy corrupted file
======================================================
         input: /Users/luomc/go/src
        output: /Users/luomc/go/test
     blockSize: 4K
    retryTimes: 0
  skipIfFailed: false
checkTimestamp: false
forceOverwrite: false
------------------------------------------------------
14:12:36 started
14:12:38 784/2418/0
   784: directory
  2418: file
  2418: copied
     0: reused
     0: failed
Duration:  0:00:01
Done!
```

