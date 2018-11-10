# eGo
Generates e with Go. Serial/parallel versions included  
https://github.com/NinjaJc01/eGo  

Instructions for use:  
  1. Install Go on the machine you want to compile FROM (doesn't have to be the machine you want to run it on).  
Use the website, not the package manager otherwise it's outdated
https://golang.org/  
  1. Install dependency for my code: ```go get -u github.com/ericlagergren/decimal```
  1. Build my project: ```go build ./parallel/egoParallel.go```
  To build for different architectures or operating systems:  
  ```GOARCH=arm64 GOOS=linux go build ./parallel/egoParallel.go```  
  Would build for arm64, linux (raspi etc, sometimes you may need just arm)  
  See entries in bold for supported OS and Architectures out of the box: https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63  
  1. Make sure you FILTER the output:  
  \*NIX ```| grep e```  
  Powershell ```| Select-String e```  
  CMD ```| find "e"```  
  1. New: try -hard as a flag, sets the iterations and precision higher to give your CPU a workout! Overrides any set precision or iterations!
