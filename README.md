# goMidRain
golang for envoy

## Clone & Build
1. Clone

        git clone https://github.com/Gemrails/goMidRain.git`
  
2. Build use docker container `golang:1.7.3`

        cd goMidRain
        
        docker run -it --rm -v $PWD:/mnt -w /mnt golang:1.7-alpine bash
        
        //in container
        
        cp vendor/* /go/src/
        
        cd /mnt
        
        go build -o mid_rain
  
## Usage
1. Start model

        ./mid_rain -s
        
2. Running model
    
        ./mid_rain -r
