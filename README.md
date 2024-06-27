# moba-converter-go

moba-converter-go is a Go application for [brief description of your project].

## Build Instructions

To build and run the project, follow these steps:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/DominikBauer1/moba-converter-go.git
   cd moba-converter-go
   ```
2. **Build the project:**
    ```bash
    go build -o moba-converter-go.exe
    ```
3. **Run the executable:**
    ```bash
    ./moba-converter-go.exe
    ```
4. **(Optional) Cross-compilation:**
    If you need to build for a different operating system or architecture, use environment variables GOOS and GOARCH:
    ```bash
    # Example: Build for Windows 64-bit from a Unix-like system
    GOOS=windows GOARCH=amd64 go build -o moba-converter-go.exe
    ```

## Usage



### Perquisites

1. Configuration File:
Ensure that `config.json` is located in the same directory as the executable or specify its location using the `--config-file` flag:

```shell
./moba-converter-go.exe --config-file /path/to/config.json
```

2. Input Data:
Input data must be provided in one of the two methods specified in the Running the conversion section.


### Data format

The input data must be a valid JSON with the following format:

```json
{
  "_meta": {
    "description": "Example JSON input file for sessions and templates"
  },
  "sessions": [
    {
      "sessionName": "abc",
      "session_type": "ssh",
      "remote_host": "1.2.3.4",
      "someotherconfig": "othervalue",
      "template": "prodsrv"
    },
    {
      "sessionName": "def",
      "session_type": "rdp",
      "remote_host": "5.6.7.8",
    }
  ],
  "templates": {
    "prodsrv": {
      "tab_color": "16711680",
      "someotherconfig": "othervalue"
    }
  }
}
```

To get information on all possible values, use the `--value-info` flag:

```shell
./moba-converter-go.exe --value-info
```

### Running the conversion

To run the conversion you need to provide the converter with the data and it will print out a mobaxterm file to stdout.

*Note:* All log and  messages error messages which may be shown are printed to stderr.

moba-converter-go can accept session data in one of two ways: 

1. **From sdtin**
In its default mode moba-converter-go expects json data from stdin until EOF. The easiest way to achieve this is to pipe data from another tool to moba-converter-go.
This is very helpful when obtaining the session data from another tool via an api wrapper.

Example: 
```bash
# Pipe data from api script
my-api-wrapper | moba-converter-go.exe 1> your-new-mobafile.mxtsessions
```

2. **From File**
If you have an existing file with json data it can be used by using the --input flag.

Example: 
```bash
# Read json data from file
moba-converter-go.exe --input input.json 1> your-new-mobafile.mxtsessions
```

# upcoming features
- output file
- specify template with flag
- add support for bookmarks



...