package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"
//// import "fmt"


type Coordinator struct {
	// Your definitions here.
    filesNotMappedYet []string
	nReduce           int
    reduceUnfinished  int
	mapTasksTotal     int
	ending            bool
}

// Your code here -- RPC handlers for the worker to call.
func (c *Coordinator) AllocateTask(args *ArgsRPC, reply *ReplyRPC) error {
	if args.WaitingMap == true {
        if len(c.filesNotMappedYet) != 0 {
			reply.Y = args.X + 1
			reply.Map = true
			reply.Reduce = false
			reply.N = c.nReduce

			//// fmt.Println(c.filesNotMappedYet)
			lastFile := c.filesNotMappedYet[len(c.filesNotMappedYet) - 1]
			lastIdx  := len(c.filesNotMappedYet) - 1
			reply.File = lastFile
			reply.MapTaskId = lastIdx + 1
			c.filesNotMappedYet = append(c.filesNotMappedYet[:lastIdx], c.filesNotMappedYet[lastIdx + 1:]...)
			//// fmt.Println(c.filesNotMappedYet)
		} else {
            // Map and Reduce Done
			if c.reduceUnfinished == -1 {
				c.ending = true
				reply.Map = false
				reply.Reduce = false
				reply.Y = args.X + 10
                return nil
			}

			reply.Y = args.X + 2
			reply.Map = false
			reply.Reduce = true
			reply.N = c.nReduce
			reply.ReduceTaskId = c.reduceUnfinished
			c.reduceUnfinished = c.reduceUnfinished - 1
			reply.MapTaskNum   = c.mapTasksTotal
		}
	
	}
	return nil
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}


//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.
    if c.ending == true {
		ret = true
	}

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	// Notes
	// nReduce is the number of reduce tasks: 10
	// nReduce intermediate files will be created.
	c := Coordinator{}

	// Your code here.
    c.filesNotMappedYet = files
	c.nReduce = nReduce
	c.reduceUnfinished  = nReduce - 1
    c.mapTasksTotal     = len(files)
    c.ending  = false

	c.server()
	return &c
}
