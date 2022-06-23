### About The Project

Meant to be an easy way to launch programs or run commands 
by saving them in a configuration file.

[![Product Name Screen Shot][product-screenshot]](https://example.com)

Inspired by Jaime

https://github.com/juanibiapina/jaime

### Built With

* [Go](https://go.dev/)
* [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## Download

Go to the releases page

https://github.com/Enchufadoo/launshr/releases

## Getting Started

You will need a simple json config file with the command to run, name and wd (working directory)

`name` is optional, a description of the command

`command` program to run

`wd` Working directory in which to run the command, by the default the working directory
is the folder in which launshr is run.

#### Submenus
If you create an object without a command, that object becomes a submenu, 

#### $config

`$config` is an object that is inherited by the commands of the object and submenus
for now it only accepts the wd key. 


``

```json
{
  "first": {
    "command": "ls -lah /",
    "name": "List the root directory",
    "wd": "."
  },
  "second": {
    "command": "wget google.com"
  },
  "emptySubmenu": {
    
  },
  "submenuWithConfig": {
    "third": {
      "command": "wget yahoo.com"
    },
    "$config": {
      "wd": "/root"
    }
  }
}
```

Save the file and pass it as the first argument of launshr

```bash
./launshr file.json
```

## Roadmap

- [x] Commit something that works
- [x] Add README.md
- [ ] Make a better screenshot
- [ ] Make it good

[product-screenshot]: images/screenshot.png
