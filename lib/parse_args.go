package parse_args
//package main

import (
	"fmt"
	"os"
)

type ArgStruct struct {
	Name string
	Format []string
	GetNext bool
}

type ArgOutput struct {
	ArgMap map[string]string
	Rest []string
}

func RemoveIndex(s []string, index int) []string {
    return append(s[:index], s[index+1:]...)
}

func ParseArgs(input_args []string, to_find []ArgStruct) (end_output ArgOutput) {
	var output []string
	GetNext := false
	found := ""
	to_output := true
	output_map := make(map[string]string)
	//
	for in_i, in_a := range input_args {
		to_output = true
		_ = in_i
		if (GetNext) {
			//fmt.Println("Foud:", found, in_a)
			output_map[found] = in_a
			GetNext = false
			to_output = false
		} else {
			for find_i, find_a := range to_find {
				_ = find_i
				for format_i, format_a := range find_a.Format {
					_ = format_i
					//fmt.Println(format_a, in_a)
					if (format_a == in_a) {
						found = find_a.Name
						if (find_a.GetNext) {
							GetNext = find_a.GetNext
						} else {
							//fmt.Println("Foud:", found, true)
							output_map[found] = "true"
							//fmt.Println(in_i)
						}
						to_output = false
						break
					}
				}
				if !(to_output) {
					break
				}
			}
			if (to_output) {
				if (in_a[:2] == "--") || (in_a[:1] == "-") {
					fmt.Println("Unrecognized parameter:", in_a)
					os.Exit(1)
				}
			}
		}
		if (to_output) {
			output = append(output, in_a)
		}
	}
	end_output = ArgOutput{ArgMap: output_map, Rest: output}
	//fmt.Println(output)
	//fmt.Println(output_map)
	//fmt.Println(input_args)
	//fmt.Println(to_find)
	return
}

func main() {
	var command ArgOutput
	format_namespace := []string {"--namespace", "-n"}
	fromat_verbose := []string {"--verbose", "-v"}
	var to_find = []ArgStruct {
		ArgStruct {
			Name: "namespace",
			Format: format_namespace,
			GetNext: true,
		},
		ArgStruct {
			Name: "verbose",
			Format: fromat_verbose,
			GetNext: false,
		},
	}

    argsWithoutProg := os.Args[1:]
    command = ParseArgs(argsWithoutProg, to_find)
    fmt.Println(command.Rest[0], command.ArgMap["namespace"], command.ArgMap["verbose"])
}
