![Fidy](./assets/fidy.png)

# Fidy

Fidy is a command-line tool for organizing files in a directory based on their extensions. It helps you keep your workspace tidy by categorizing files into folders named after their extensions.

## Features

- **Organize Files by Extension**: Automatically move files into directories based on their file extensions.
- **Verbose Mode**: Get detailed output of the operations performed.
- **Dry Run**: Simulate the organization process without making any changes.
- **Custom Greetings**: Save a custom name to personalize your Fidy experience.
- **Custom Extensions**: Add or remove extensions to be considered using -include/-exclude flags.

## Installation

Make sure you have [Go](https://go.dev/doc/install) installed in your system.

You can install Fidy using the following command in your terminal:
```
go install github.com/ashmit0920/Fidy@latest
```

Run fidy
```
fidy
```

## Usage

Display info about Fidy
```
fidy -help
```

Fidy can remember your name and send personalized greetings! Set your name using the command:
```
fidy -name YOUR_NAME
```

Organize a directory
```
fidy -dir PATH_TO_YOUR_DIR
```

Include or exclude specific file extensions
```
fidy -dir PATH_TO_YOUR_DIR -include txt,png -exclude pdf
```

Organize files with verbose output
```
fidy -dir PATH_TO_YOUR_DIR -verbose
```

Perform a Dry Run, without doing any actual changes to your directory.
```
fidy -dir PATH_TO_YOUR_DIR -dryrun -include txt,png -exclude pdf
```

## Upcoming features

- Log file
- Delete empty directories automatically
- Custom directory naming
- Backup before running Fidy
- Recursive organization for sub-directories