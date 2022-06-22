package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Arguments map[string]string

const operation = "operation"
const filename = "fileName"
const item = "item"
const id = "id"

/* type items struct {
    Items []ItemO `json:"items"`
} */

type ItemO struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int64  `json:"age"`
}

func isValidOps(ops string) bool {
	switch ops {
	case
		"add",
		"list",
		"findById",
		"remove":
		return true
	}
	return false
}

func parseArgs() Arguments {
	ops := flag.String(operation, "", "specify operation")
	file := flag.String(filename, "", "specify filename")
	it := flag.String(item, "", "specify item")
	i := flag.String(id, "", "specify id")

	flag.Parse()
	args := Arguments{
		operation: *ops,
		item:      *it,
		filename:  *file,
		id:        *i,
	}

	return args
}

func Perform(args Arguments, writer io.Writer) error {
	var errors []string
	if args[operation] == ""{
		return fmt.Errorf("-%s flag has to be specified", operation)
	}else if !isValidOps(args[operation]){
		return fmt.Errorf("Operation %s not allowed!", args[operation])
	}

	 if args[filename] == "" {
		return fmt.Errorf("-%s flag has to be specified", filename)
	}
	/*
		required := []string{operation, filename}
		present := make(map[string]bool)

		flag.Visit(func(f *flag.Flag) { present[f.Name] = true })
		for _, req := range required {
			if !present[req] {
				errors = append(errors, fmt.Sprintf("-%s flag has to be specified", req))
				//os.Exit(2) // the same exit code flag.Parse uses
				if v := strings.Join(errors, ""); len(v) > 0 {
					return fmt.Errorf(v)
				}

			}
		} */

	switch args[operation] {

	case "add":
		if args[item] == "" {
			errors = append(errors, fmt.Sprintf("-%s flag has to be specified", item))
		} else {
		// create a new decoder
		var itemList []ItemO
		itemDecoder := json.NewDecoder(strings.NewReader(readFile(args[filename])))
		err := itemDecoder.Decode(&itemList)
		if err != nil {
			fmt.Println(err)
			//os.Exit(2)
		}
		//parse items
		var it ItemO
		err = json.Unmarshal([]byte(args[item]),&it)
		if err != nil {
			fmt.Println(err)
			//os.Exit(2)
		}
		
		for i := range itemList {
			
			if itemList[i].Id == string(it.Id) {
				writer.Write([]byte(fmt.Sprintf("Item with id %s already exists",it.Id)))
				return nil
			}
		}
			
			itemList = append(itemList, it)
			items, err := json.Marshal(itemList)
			if err != nil {
				fmt.Println(err)
				//os.Exit(2)
			}
			err = os.WriteFile(args[filename], []byte(items), os.ModeAppend)
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}
	case "list":
		//read file and print
		writer.Write([]byte(readFile(args[filename])))
	case "findById":
		//check for id passed
		var err error
		if args[id] == "" {
			errors = append(errors, fmt.Sprintf("-%s flag has to be specified", id))
		} else {
			// create a new decoder
			var itemList []ItemO

			itemDecoder := json.NewDecoder(strings.NewReader(readFile(args[filename])))
			err = itemDecoder.Decode(&itemList)
			if err != nil {
				fmt.Println(err)
			//	os.Exit(2)
			}
			for i := range itemList {

				if itemList[i].Id == string(args[id]) {
					item, err := json.Marshal(itemList[i])
					if err != nil {
						fmt.Println(err)
						os.Exit(2)
					}
					writer.Write(item)
				} else {
					writer.Write([]byte(""))
				}
			}
		}
	case "remove":
		//check for id passed
		var err error
		if args[id] == "" {
			errors = append(errors, fmt.Sprintf("-%s flag has to be specified", id))
		} else {
			// create a new decoder
			var itemList []ItemO
			deleteIndex := -1
			itemDecoder := json.NewDecoder(strings.NewReader(readFile(args[filename])))
			err = itemDecoder.Decode(&itemList)
			if err != nil {
				fmt.Println(err)
			//	os.Exit(2)
			}
			for i := range itemList {
				if itemList[i].Id == string(args[id]) {
					deleteIndex = i
					break
				}
			}
			if deleteIndex >= 0 {
				itemList = append(itemList[:deleteIndex], itemList[deleteIndex+1:]...)
				err := os.Remove(args[filename])
				if err != nil {
					fmt.Println(err)
					//os.Exit(2)
				}
				items, err := json.Marshal(itemList)
				if err != nil {
					fmt.Println(err)
				//	os.Exit(2)
				}
				err = os.WriteFile(args[filename], []byte(items), 0777)
				if err != nil {
					fmt.Println(err)
					os.Exit(2)
				}
			} else {
				writer.Write([]byte(fmt.Sprintf("Item with id %s not found", args[id])))
			}
		}
	default:
		errors = append(errors, fmt.Sprintf("Operation %s not allowed!", args[operation]))
	}
	if v := strings.Join(errors, ""); len(v) > 0 {
		return fmt.Errorf(v)
	}
	return nil
}

/* func isValidItem(item ItemO)bool{
   if item.Id !=nil && item.Age!=nil && item.Email !=nil {

   }
} */

func readFile(fileName string) string {

	f, err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE, 420)
	var content string
	if err != nil {
		fmt.Printf("unable to read file: %v", err)
		os.Exit(2)
	}
	defer f.Close()
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			//writer.Write(buf[:n])
			content += string(buf[:n])
		}
	}
	return content
}
/* func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
} */
