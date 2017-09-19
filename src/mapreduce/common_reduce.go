package mapreduce

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// doReduce manages one reduce task: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	kvmap := make(map[string][]string)
	decs := []*json.Decoder{}

	for i := 0; i < nMap; i++ {
		imdfilename := reduceName(jobName, i, reduceTaskNumber)
		imdfile, _ := os.Open(imdfilename)
		defer imdfile.Close()
		dec := json.NewDecoder(imdfile)
		decs = append(decs, dec)
	}

	for i := 0; i < nMap; i++ {
		for {
			var kv KeyValue
			if err := decs[i].Decode(&kv); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(reduceName(jobName, i, reduceTaskNumber))
				panic(err)
			}
			if _, ok := kvmap[kv.Key]; ok {
				kvmap[kv.Key] = append(kvmap[kv.Key], kv.Value)
			} else {
				kvmap[kv.Key] = []string{kv.Value}
			}
		}
	}

	file, _ := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	enc := json.NewEncoder(file)

	for key, values := range kvmap {
		enc.Encode(KeyValue{key, reduceF(key, values)})
	}
	//
	// You will need to write this function.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTaskNumber) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
}
