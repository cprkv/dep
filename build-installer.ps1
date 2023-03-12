Remove-Item release/installer.exe -ErrorAction SilentlyContinue
go build -ldflags '-s' -o installer/dep.exe
Set-Location installer
Remove-Item -r Output -ErrorAction SilentlyContinue
."C:\Program Files (x86)\Inno Setup 6\iscc.exe" .\main.iss
Set-Location ..
Move-Item installer/Output/dep-installer.exe release/installer.exe