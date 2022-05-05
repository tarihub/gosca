## GoSCA [![Go 1.17](https://img.shields.io/badge/Go-1.17-blue.svg)](https://go.dev/)

GoSCA scans a Go project for vulnerable dependencies.

### Usage
```bash
./gosca -w /path/to/workdir
```

![dash-h.png](doc/dash-h.png)

Running

![run.png](doc/run.png)

### Download
Select your platform for download

https://github.com/TARI0510/gosca/releases

### Build
```bash
go build -o gosca cmd/gosca/main.go
```

or for cross platform builds
```bash
sh package.sh
```

### References
1. https://github.com/securego/gosec
2. https://github.com/fatedier/frp
