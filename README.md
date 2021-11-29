# pswdmngr

A lightweight Secure CLI Password Manager using Go

![image](https://user-images.githubusercontent.com/74897008/143892576-cdeaa8a1-c37f-4a7a-9b96-561d33da6334.png)
![image](https://user-images.githubusercontent.com/74897008/143892793-280abbd4-c783-4fc2-aa59-69669ceba057.png)

## Installation

NOTE: The binaries are only for x86, 64bit devices. For others, please compile from source.

The required binaries can be installed directly from the releases tab: `pswdmngr` for Linux, `pswdmngr.exe` for Windows.

Alternative, <br/>
Get the required files
```
git clone https://github.com/resyfer/pswdmngr
cd pswdmngr
```

### Linux
Install
```
./install.sh
```

Run with <br/>
NOTE: The `password.txt` file gets initialised in your current directory.
```
pswdmngr
```

### Windows
The bin folder has the executable `pswdmngr.exe`. Copy to preferred location and run from there.


## Compiing From Source
Pre-requirements: `Go 1.13 or higher`

`cd` into project folder
```
rm -rf bin
mkdir bin
go build -o ./bin
cd bin
```
and run the executable binary
