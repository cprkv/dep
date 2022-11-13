go build -ldflags '-s' -o installer/dep.exe
cd installer
rm -r Output -ErrorAction SilentlyContinue
."C:\Program Files (x86)\Inno Setup 6\Compil32.exe" /cc .\main.iss
cd ..
