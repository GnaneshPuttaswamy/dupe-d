package main

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	extensions []string
)

type HashedFileInfo struct {
	Name string
	Path string
	Size int64
	Hash string
}

var rootCmd = &cobra.Command{
	Use:   "dupe-d [directory]",
	Short: "dupe-d is a tool to identify file duplicates",
	Long: `dupe-d is a tool to identify file duplicates by generating sha-256 hash.
	To scan the current directory, use: dupe-d .`,
	Example: `  dupe-d 
  dupe-d /path/to/directory
  dupe-d --ext jpg --ext png /path/to/directory
  dupe-d --ext=jpg,png,pdf`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		folderPath, err := getFolderPath(args)
		if err != nil {
			return err
		}

		formattedExtensions := formatExtensions(extensions)

		hashedFilesInfo, err := processFiles(folderPath, formattedExtensions)
		if err != nil {
			return err
		}

		err = writeToCsv(hashedFilesInfo)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringSliceVarP(&extensions, "ext", "e", []string{}, "File extensions to process (can be specified multiple times or comma-separated)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		printToStdErr(err)
	}
}

func getFolderPath(args []string) (string, error) {
	if len(args) > 0 {
		return validateDirectory(args[0])
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return currentDir, nil
}

func validateDirectory(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("directory not accessible: %w", err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", path)
	}

	return path, nil
}

func formatExtensions(rawExts []string) []string {
	var formattedExts []string

	for _, ext := range rawExts {
		splitExts := strings.Split(ext, ",")

		for _, splitExt := range splitExts {
			splitExt = strings.TrimSpace(splitExt)

			if splitExt == "" {
				continue
			}

			if !strings.HasPrefix(splitExt, ".") {
				splitExt = "." + splitExt
			}

			formattedExts = append(formattedExts, splitExt)
		}
	}

	return formattedExts
}

func processFiles(folderPath string, exts []string) ([]HashedFileInfo, error) {
	printToStdOut(fmt.Sprintf("Scanning folder: %s\n", folderPath))
	if len(exts) > 0 {
		printToStdOut(fmt.Sprintf("Filtering by extensions: %s\n", strings.Join(exts, ", ")))
	} else {
		printToStdOut("Processing all file types\n")
	}

	var files []HashedFileInfo

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if matchesExtension(path, exts) {
			printToStdOut(fmt.Sprintf("Processing: %s\n", path))

			hash, err := hashFile(path)
			if err != nil {
				return fmt.Errorf("failed to hash file %s: %w", path, err)
			}

			info, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("failed to get file stats for %s: %w", path, err)
			}

			fileInfo := HashedFileInfo{
				Name: info.Name(),
				Size: info.Size(),
				Hash: hash,
				Path: path,
			}

			files = append(files, fileInfo)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := sha256.New()
	buf := make([]byte, 1024*1024)

	_, err = io.CopyBuffer(hash, file, buf)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func writeToCsv(hashedFilesInfo []HashedFileInfo) error {

	timestamp := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	outputFilename := fmt.Sprintf("hash_results_%s.csv", timestamp)

	file, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Name", "Path", "Size (MB)", "Hash"})
	if err != nil {
		return fmt.Errorf("failed to write header to CSV: %w", err)
	}

	for _, hashedFileInfo := range hashedFilesInfo {

		sizeInMB := float64(hashedFileInfo.Size) / 1048576.0

		err = writer.Write([]string{
			hashedFileInfo.Name,
			hashedFileInfo.Path,
			fmt.Sprintf("%.2f", sizeInMB),
			hashedFileInfo.Hash,
		})
		if err != nil {
			return fmt.Errorf("failed to write content to CSV: %w", err)
		}
	}

	absPath, err := filepath.Abs(outputFilename)
	if err != nil {
		absPath = outputFilename
	}

	printToStdOut(fmt.Sprintf("Output written to: %s\n", absPath))

	return nil
}

func matchesExtension(path string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}

	ext := filepath.Ext(path)
	for _, e := range exts {
		if ext == e {
			return true
		}
	}

	return false
}

func printToStdErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
}

func printToStdOut(s string) {
	fmt.Fprint(os.Stdout, s)
}
