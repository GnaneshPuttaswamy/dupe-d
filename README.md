# dupe-d

A lightweight command-line tool to identify duplicate files using SHA-256 hash comparison.

## Features

- Scan directories recursively to find duplicate files
- Filter by file extensions
- Generate detailed CSV report with file information
- Fast performance with efficient hashing algorithm

## Installation

### Prerequisites

- Go 1.23 or higher installed on your system

### From Source

```bash
# Clone the repository
git clone https://github.com/GnaneshPuttaswamy/dupe-d.git
cd dupe-d

# Build the binary
go build -o dupe-d

# Optional: Move to a location in your PATH
mv dupe-d /usr/local/bin/
```

## Usage

```bash
# Scan current directory
dupe-d

# Scan a specific directory
dupe-d /path/to/directory

# Scan with file extension filtering
dupe-d --ext jpg --ext png /path/to/directory

# Using comma-separated extensions
dupe-d --ext=jpg,png,pdf /path/to/directory
```

## Options

| Flag    | Short | Description                                                    |
| ------- | ----- | -------------------------------------------------------------- |
| `--ext` | `-e`  | File extensions to process (comma-separated or multiple flags) |

## Output

The tool generates a timestamped CSV file (`hash_results_YYYYMMDD_HHMMSS.csv`) containing:

- File name
- Full path
- File size (in MB)
- SHA-256 hash

Duplicate files will have identical hash values, making them easy to identify.

## Example Output

```bash
Scanning folder: /path/to/directory
Filtering by extensions: .jpg, .png
Processing: /path/to/directory/image1.jpg
Processing: /path/to/directory/image2.jpg
Processing: /path/to/directory/image3.png
Output written to: /path/to/directory/hash_results_20250101_120000.csv
```

## How to Find Duplicates

After running the tool, open the generated CSV file in any spreadsheet software and:

1. Sort by the "Hash" column
2. Files with identical hash values are duplicates

## License

This project is licensed under the MIT License - see the LICENSE file for details.
