package mr

import "fmt"
import "log"
import "net/rpc"
import "hash/fnv"

import "os"
import "io/ioutil"
import "encoding/json"
import "sort"



//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}


//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue, reducef func(string, []string) string) {

	// Your worker implementation here.
    ReplyFromCoor := AskForTask()
	fmt.Printf("Oops - ReplyFromCoor.Y %v [OK]\n", ReplyFromCoor.Y)

	if ReplyFromCoor.Map == true {

		fmt.Printf("Oops - Get Map Tasks\n")
		fmt.Println("Map Task with File <", ReplyFromCoor.File, ">")
		
		intermediate := []KeyValue{}
		eachFile, err := os.Open(ReplyFromCoor.File)
		if err != nil {
			log.Fatalf("cannot open %v", ReplyFromCoor.File)
		}
		content, err  := ioutil.ReadAll(eachFile)
		if err != nil {
			log.Fatalf("cannot read %v", ReplyFromCoor.File)
		}
		eachFile.Close()
		
		kva := mapf(ReplyFromCoor.File, string(content))
		intermediate = append(intermediate, kva...)
        sort.Sort(ByKey(intermediate))

		// allocate nReduce intermediate files
		//
		// ==== Create the JSON files ====
		for i := 0; i < ReplyFromCoor.N; i++ {
			filename := fmt.Sprintf("mr-%d-%d.json", ReplyFromCoor.MapTaskId, i)
		    file, err := os.Create(filename)	
		    if err != nil {
				panic(err)
			}
			defer file.Close()
		}
		// ========== Allocate ==========
		for _, kv := range intermediate {
			filename := fmt.Sprintf("mr-%d-%d.json", ReplyFromCoor.MapTaskId, ihash(kv.Key) % ReplyFromCoor.N)
			file, openErr := os.OpenFile(filename, os.O_RDWR | os.O_APPEND, 0644)
			if openErr != nil {
				panic(openErr)
			}
			enc := json.NewEncoder(file)
			err := enc.Encode(kv)
			if err != nil {
				panic(err)
			}
		}


	} else if ReplyFromCoor.Reduce == true {
		fmt.Printf("Oops - Get Reduce Tasks\n")

		kva := []KeyValue{}
		for i := 1; i <= ReplyFromCoor.MapTaskNum; i++ {
			filename := fmt.Sprintf("mr-%d-%d.json", i, ReplyFromCoor.ReduceTaskId)
			file, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			dec := json.NewDecoder(file)
			for {
				var kv KeyValue
				if err := dec.Decode(&kv); err != nil {
					break
				}
				kva = append(kva, kv)
			}
			defer file.Close()
		}
		// All in one Reduce Task has been collected.
		// Do Reduce
		outputFileName := fmt.Sprintf("mr-out-%d", ReplyFromCoor.ReduceTaskId)
		ofile, _ := os.Create(outputFileName)
		i := 0
		for i < len(kva) {
			j := i + 1
			for j < len(kva) && kva[j].Key == kva[i].Key {
				j++
			}
			values := []string{}
			for k := i; k < j; k++ {
				values = append(values, kva[k].Value)
			}
			output := reducef(kva[i].Key, values)
			fmt.Fprintf(ofile, "%v %v\n", kva[i].Key, output)
			
			i = j
		}
		ofile.Close()

	} else {
		fmt.Printf("Oops - No Tasks at all!\n")
		fmt.Printf("Oops - ReplyFromCoor.Y %v [END]\n", ReplyFromCoor.Y)
	}

	// --------------------------------------------------------
	// -- Begin with this step --
	// uncomment to send the Example RPC to the coordinator.
	////CallExample()
}













// ============================================================================================
// A Initialized Worker ask for a task (RPC)
func AskForTask() ReplyRPC {
	args  := ArgsRPC{}
	args.X          = 99
	args.WaitingMap = true
	reply := ReplyRPC{} 
	
	ok := request("Coordinator.AllocateTask", &args, &reply)
	if ok {
        return reply
	} else {
		fmt.Printf("call failed!\n")
		return ReplyRPC{}
	}
}
func request(rpcname string, args interface{}, reply interface{}) bool {
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}









// ============================================================================================
//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {
	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply. the "Coordinator.Example" tells the
	// receiving server that we'd like to call the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
		// It does.
	} else {
		fmt.Printf("call failed!\n")
	}
}
//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname  := coordinatorSock()

	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
