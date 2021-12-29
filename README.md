# Canary
CLI tool written in Go to generate Canary Tokens from https://canarytokens.org

## Usage
```
$ canary -h 
Usage of ./canary:
  -email string
    	Email to configure token on
  -memo string
    	Description of token to remember it by
  -token string
    	DNS token to prepend generic data to. Pipe data from stdin
  -type string
    	Type of token to create (dns or web) (default "dns")
```

#### Examples
```
$ canary -email test@gmail.com
1wbzrTESTm74hu2.canarytokens.com
```
```
$ canary -email test@gmail.com -memo "Someone visited the site" -type web
http://canarytokens.com/images/tags/feedback/d3c5TESTvdnejy0b/contact.php
```
```
$ echo 'test' | canary -token 0mcsepoTESTsr9ilsjo.canarytokens.com
ORSXG5A.G23.0mcsepoTESTsr9ilsjo.canarytokens.com
```
```
$ cat data.txt | canary -token 0mcsepoTESTsr9ilsjo.canarytokens.com
KRSXG5BR.G78.0mcsepoTESTsr9ilsjo.canarytokens.com
KRSXG5BS.G20.0mcsepoTESTsr9ilsjo.canarytokens.com
...
```
