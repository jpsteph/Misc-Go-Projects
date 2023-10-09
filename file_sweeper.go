package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func main() {

	source_dir, _ := os.Getwd()

	type file_struct struct {
		directory string
		size      int
	}

	file_list_sorted := []file_struct{}

	var files []string
	dir_str := "**"
	next_dir_str := "/**"

	// Get all files and directories in the source directory and its subdirectories
	pattern := filepath.Join(source_dir, dir_str)
	files_temp, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("Error getting file paths in directory")
		os.Exit(1)
		//return err
	}
	files = append(files, files_temp...)
	file_num := len(files_temp)

	for file_num != 0 { //keeping getting all file paths until there are no files/folders to go into
		dir_str = dir_str + next_dir_str //going one directory down

		pattern := filepath.Join(source_dir, dir_str)
		files_temp, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Println("Error getting file paths in directory")
			os.Exit(1)
			//return err
		}

		files = append(files, files_temp...)
		file_num = len(files_temp)
	}

	for i := 0; i < len(files); i++ {
		file_info, err := os.Stat(files[i])
		if err != nil {
			fmt.Println("Error getting directory info")
		}

		if !file_info.IsDir() { //skip directory if it is just a folder
			file_data, err := os.ReadFile(files[i])
			if err != nil {
				fmt.Println("Error reading file:", err.Error())
				os.Exit(1)
			}

			file_size := len(file_data) / 1000 //str of file size in kB
			//fmt.Println("File: " + files[i] + " File Size: " + strconv.Itoa(file_size) + "kb") //debug

			file_list_sorted = append(file_list_sorted, file_struct{files[i], file_size})
		}
	}

	// Define a custom sorting function
	sort.Slice(file_list_sorted, func(i, j int) bool {
		return file_list_sorted[i].size < file_list_sorted[j].size // Sort in descending order
	})

	file, err := os.Create("File_Sweeper_Output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 0; i < len(file_list_sorted); i++ {
		fmt.Println("File Name: " + file_list_sorted[i].directory)
		fmt.Println("File Size: " + strconv.Itoa(file_list_sorted[i].size) + " kb\n")

		_, err = writer.WriteString("File Name: " + file_list_sorted[i].directory + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}
		_, err = writer.WriteString("File Size: " + strconv.Itoa(file_list_sorted[i].size) + " kb\n\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}
	}

	writer.Flush()
}
