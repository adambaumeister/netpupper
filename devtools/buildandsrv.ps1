 param (
    [switch]$rpi = $false
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

c:\Temp\netp.exe -server -address='127.0.0.1:5000'