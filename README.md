# moba-converter-go

moba-converter-go is a Go application to create MobaXterm Session files by using json data. 

Author: DominikBauer1

The session file format was used from: [.mxtsessions file format by Ruzgfpegk](https://gist.github.com/Ruzgfpegk/ab597838e4abbe8de30d7224afd062ea)

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
    Use the go documentation for more information on cross compilation.


## Usage

### Supported sessions

- ssh 


### Perquisites

1. **Configuration File:**
Ensure that `config.json` is located in the same directory as the executable or specify its location using the `--config-file` flag. The config.json serves as the main mapping table between `json` and `.mxtsessions` data. It may needed to be updated to work with future versions of MobaXterm.

```shell
./moba-converter-go.exe --config-file /path/to/config.json
```

2. **Input Data:**
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
      "SessionName": "first session",
      "sessionType": "ssh",
      "RemoteHost": "1.2.3.4",
      "Port": "2222",
      "template": "prodsrv"
    },
    {
      "SessionName": "another session",
      "sessionType": "rdp",
      "RemoteHost": "5.6.7.8",
      "template": "prodsrv"
    }
  ],
  "templates": {
    "prodsrv": {
      "CustomTabColor": "16711680",
      "Username": "my-prod-user"
    }
  }
}
```

To get information on all possible Options, use the `--value-info` flag:

```shell
./moba-converter-go.exe --value-info
```

### Special Keys
There are two keys which can be used in the input, but do not directly correspond to MobaXterm Setting: 

#### FolderPath (TODO)
With this key you can specify the folder path in which the session should reside.

The folders use `/` to separate folders and start with `/` as the root folder object.
The root folder is implicitly assumed and does not need to be specified.
Trailing slashes are ignored.

Example with root folder
```json
"FolderPath": "/"
```

Example with subfolder
```json
"FolderPath": "/my-ssh-sessions"
```

Folders will be Implicitly created as soon as they are used at least once in a session path. To customize the icon please use the optional key "folders" in the input data.


```json
{
  "_meta": {
    "description": "Example for Folder Section"
  },
  "sessions": [    
      {
        "SessionName": "another session",
        "sessionType": "ssh",
        "RemoteHost": "5.6.7.8",
        "template": "prodsrv",
        "folder": "/Test"
      }
    ],
  "templates": {},
  "folders":{
    "/folder1":{
      "Icon":"32"
    }
  }
}
```

#### Templates

The Template key allows to apply a set of options to multiple sessions.

*Note:* This is not a Mobaxterm feature and should not be seen as one.

Templating works by creating a set of options which then act as the default options for sessions which are connected to the template. 
This also means that template values only work if the value is NOT explicitly set in the session itself.

To create a template, add a section to the "templates" section in the input data.
Then add the "template" key to one or more sessions to apply the options.

```json
{
  "_meta": {
    "description": "Example JSON input file for sessions and templates"
  },
  "sessions": [
    {
      "SessionName": "first session",
      "sessionType": "ssh",
      "RemoteHost": "1.2.3.4",
      "template": "my-first-template"
    },
  ],
  "templates": {
    "my-first-template": {
      "CustomTabColor": "16711680",
      "someotherconfig": "othervalue"
    }
  }
}
```


*Possible Future Options*: 
- Allow for multiple templates to be applied to one session.
- Allow for templates to have the template key and allow for recursive templates

### Running the conversion

To run the conversion you need to provide the converter with the json data and it will print out a mobaxterm file to converted.mxtsessions.
To change the output file, use the --output <> flag.

*Note:* All log and  messages error messages which may be shown are printed to stderr.

moba-converter-go can accept session data in one of two ways: 

1. **From sdtin**
In its default mode moba-converter-go expects json data from stdin until EOF. The easiest way to achieve this is to pipe data from another tool to moba-converter-go.
This is very helpful when obtaining the session data from another tool via an api wrapper.

Example: 
```bash
# Pipe data from api script
my-api-wrapper | moba-converter-go.exe --output my-sessions.mxtsessions
```

2. **From File**
If you have an existing file with json data it can be used by using the --input flag.

Example: 
```bash
# Read json data from file
moba-converter-go.exe --input input.json --output your-new-mobafile.mxtsessions
```


# Other
## TODO
- specify template from separate file with flag
- make template key optional
- fix some stuff with calculated vars
- ENV var to set basic parameters like
  - don't allow risky options