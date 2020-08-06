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
