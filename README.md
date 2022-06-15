<div id="top"></div>

### About The Project

Meant to be an easy way to launch programs or run commands by saving them in a config file.

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

You will need a simple json config file with the command to run, name (optional) and working directory (optional)

```json
{
  "first": {
    "command": "ls -lah /",
    "name": "List the root directory",
    "wd": "."
  },
  "second": {
    "command": "wget google.com"
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
