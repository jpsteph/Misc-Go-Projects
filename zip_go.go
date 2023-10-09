package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func compressFiles(sourceDir string, destZip string) error {

	//will use later to skip zipping the exe
	exeName := os.Args[0]

	//if zip file exists remove it
	_, err_file := os.Stat(destZip)

	if err_file == nil {
		os.Remove(destZip)
	}

	var files []string
	dir_str := "**"
	next_dir_str := "/**"

	// Get all files and directories in the source directory and its subdirectories
	pattern := filepath.Join(sourceDir, dir_str)
	files_temp, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	files = append(files, files_temp...)
	file_num := len(files_temp)

	for file_num != 0 {
		dir_str = dir_str + next_dir_str

		// Get all files and directories in the source directory and its subdirectories
		pattern := filepath.Join(sourceDir, dir_str)
		files_temp, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		files = append(files, files_temp...)
		file_num = len(files_temp)
	}

	// Create the destination zip file
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a new zip archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Loop through all the files and directories
	for i, path := range files {
		fmt.Println("File " + strconv.Itoa(i) + " of " + strconv.Itoa(len(files)) + "...")

		//if this is the zip compress exe skip this index in the for loop
		if strings.Contains(path, exeName) {
			continue
		}

		// Get information about the file or directory
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		// If the file is a directory, skip it
		if info.IsDir() {
			continue
		}

		// Open the file for reading
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Get the relative path of the file or directory inside the source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		fmt.Println("Processing " + relPath + "...")

		// Create a new file header for the file or directory
		header := &zip.FileHeader{
			Name:   relPath,
			Method: zip.Deflate,
		}

		// Set the file or directory's modified time to the same as the original file or directory's modified time
		header.SetModTime(info.ModTime())

		// Add the file or directory header to the zip archive
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// If the file or directory is a symlink, write the symlink target to the zip archive
		if (info.Mode() & os.ModeSymlink) != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			_, err = writer.Write([]byte(link))
			if err != nil {
				return err
			}
			continue
		}

		// If the file or directory is not a symlink, copy its contents to the zip archive
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}

	fmt.Println("Compression complete")
	return nil
}

func main() {

	cwd, _ := os.Getwd()

	if (len(os.Args) > 1) && (os.Args[1] == "drive") {
		_, err := os.Stat("D:/")
		if err == nil {
			compressFiles(cwd, "D:/output.zip")
		} else {
			fmt.Println("Error, External Flash Memory Device not Available")
			os.Exit(1)
		}
	} else if len(os.Args) > 1 {
		_, err := os.Stat(os.Args[1])
		if err == nil {
			compressFiles(cwd, os.Args[1]+"/output.zip")
		} else {
			fmt.Println("Error, Directory not Found")
			os.Exit(1)
		}
	} else {
		compressFiles(cwd, cwd+"/output.zip")
	}
}
