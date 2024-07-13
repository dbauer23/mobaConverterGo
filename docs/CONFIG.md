# Config Files

The Config file is a json files, complied into moab-converter-go.

It serves as the main mapping table between `json` and `.mxtsessions` data. It may needed to be updated to work with future versions of MobaXterm.
You can also provide your own config file, using the global Flag `--configPath`

## Structure

The config file is three parts: 

* **_meta**

    *  **version**: (string) Specifies the version of the configuration file.
    * **changed_when**: (string) Indicates the date and time when the configuration file was last modified.


* **sessionTypes**: Defines options specific to different session types. Currently, only the "ssh" session type is supported. See [sessionTypes](#sessiontypes)
* **options**: Defines general configuration options that can be applied to any session type. See [options](#options).


### sessionTypes

* **ssh**
    * **tmplString**: (string) A template string used to construct the MXT session file content for SSH sessions. It includes placeholders that are replaced with corresponding values from the `options` section. 
    * **allowed_options**: (array) An empty array, likely intended for future use to restrict which options can be used for SSH sessions.

### options

Each option within this section has the following properties:

* **name**: (string) The name of the option.
* **default_value**: (string | number | boolean) The default value for the option.
* **section**: (string) The section of the MXT session file where this option is used.
* **help**: (string) A description of the option and its functionality.
* **options** (optional): (object) If present, defines a list of possible values for the option and how they are mapped to the MXT session file format. Each key-value pair within this object represents a user-friendly option name and the corresponding value used internally.

**Example Option with allowed values:**

```json
"X11Forwarding": {
  "default_value": "Enable",
  "section": "Advanced SSH settings",
  "help": "Enable X11 forwarding.",
  "options": {
    "Enable": "-1",
    "Disable": "0",
    "true": "-1",
    "false": "0"
  }
}
```

**Note:**

* Some options use special values like `"__PIPE__"` or `"__PERCENT__"`. These likely have specific purposes within the application and are not documented here.
* Not all options have a complete set of properties documented within the `help` section.

I hope this documentation provides a clearer understanding of the configuration options available in the `config.json` file.



# Cli Commands

There are several commands to view interact with config file. 
All of them can be reched using the subcommand `config`

* `info`: Print some basic meta information about the loaded config.
* `show`: Show all possible parameters for the input file available in the config file.