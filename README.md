# dconf
Dynanic config file parser for Golang, support runtime modification.

# Features
- Ini file parser
- Runtime configuration modify
- Thread safe

# Example
```
    import (
        "github.com/iswarezwp/dconf"
    )

    // Open and parse the config file, if set the reload param to true,
    // will automatically reload the configurations when you change
    // something in the config file.
    conf, err := NewDConf("/config/file/name", true)
    if err != nil {
        t.Fatal(err)
    }

    // Get values from [Default] section
    conf.Get("someKey", "defaultValueIfNotExists")
    
    // Equals ...
    conf.GetValue("Default", "someKey", "defaultValueIfNotExists")

    // Get values from other sections
    conf.GetValue("secName", "keyName", "defaultValueIfNotExists")

    // Close the watcher backend if reload is set to true
    conf.Close()
```