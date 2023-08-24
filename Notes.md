# Lecture One and Two -- MapReduce

> Distributed System tend to use GoLang!

### Lab One  
Lab 1 is to implement one MapReduce. In mrsequential.go, the sequential function has been tested ok.  

Then the CallExample could work. Next is to Create specific RPC connection.

wc.so has Map and Reduce function.  


I know, I know. How STUPID my design and my choice of data structure are!   
So hard for me to achieve this...  
```bash
(base) MacBook-Pro:main $ go build -buildmode=plugin ../mrapps/wc.go
(base) MacBook-Pro:main $ bash test-mr.sh
*** Starting wc test.
--- wc test: PASS
*** Starting indexer test.
--- indexer test: PASS
*** Starting map parallelism test.
--- map parallelism test: PASS
*** Starting reduce parallelism test.
--- reduce parallelism test: PASS
*** Starting job count test.
--- job count test: PASS
*** Starting early exit test.
--- early exit test: PASS
*** Starting crash test.
2023/08/23 23:32:23 dialing:dial unix /var/tmp/5840-mr-501: connect: connection refused
2023/08/23 23:32:23 dialing:dial unix /var/tmp/5840-mr-501: connect: connection refused
2023/08/23 23:32:23 dialing:dial unix /var/tmp/5840-mr-501: connect: connection refused
sort: No such file or directory
cmp: EOF on mr-crash-all
--- crash output is not the same as mr-correct-crash.txt
--- crash test: FAIL
*** FAILED SOME TESTS
```  
Notes for me:  
1. As soon as the worker finish one task, it would immediately go into another one! Thus, the Sleep should be placed outside the RPC handler!  
2. Better design includes TaskState, TaskQueue and so on.  
3. At first, I hope to simplify the problem which is when to end the Worker Map Tasks Allocation and when to the end the Coordinator. It turns out that I could not just end it when the worker who is the last received the task! There might be someone doing so slowly that the last one who receives the task is not the last one who completes the task!

I am too lazy to do the crash test and multi-processor test... :)