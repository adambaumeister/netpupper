 param (
    [switch]$rpi = $false,
    [switch]$buildsrv = $false
 )

go build -o c:\Temp\netp.exe
if ($rpi) {
    write-output "Doing rpi"
    $Env:GOOS = "linux"
    $Env:GOARCH = "arm"
    $Env:GOARM = 7
    go build -o c:\Temp\netp_pi
    scp c:\Temp\netp_pi pi@rpi.home:code/

    $Env:GOOS = ""
    $Env:GOARCH = ""
}

if ($buildsrv) {
    write-output "Doing buildsrv"
    $Env:GOOS = "linux"
    $Env:GOARCH = "amd64"
    go build -o c:\Temp\netp_pi_amd
    scp c:\Temp\netp_pi_amd adam@buildserver.home:

    $Env:GOOS = ""
    $Env:GOARCH = ""
}

c:\Temp\netp.exe -daemon -config='C:\users\adam\go\src\github.com\adamb\netpupper\daemon_server.yml'