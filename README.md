# Morcego

**Blind SQLI Tool**

```go

       _,    _   _    ,_
  .o888P     Y8o8Y     Y888o.
 d88888      88888      88888b
d888888b_  _d88888b_  _d888888b
8888888888888888888888888888888
8888888888888888888888888888888
YJGS8P"Y888P"Y888P"Y888P"Y8888P
 Y888   '8'   Y8P   '8'   888Y
  '8o          V          o8'
    '                     '

```				                            

## What is it
Morcego is a Blind SQL Injection tool to brute force size and values.

Morcego is designed to localhost tests so it doesn't deal with rate limits and anything other than the SQLI.

**ATENTION**: Usage of this program for attacking targets without prior mutual consent is illegal.
It is the end user's responsibility to obey all applicable local, state and federal laws.
Developers assume no liability and are not responsible for any misuse or damage caused by this program

## Build
Clone the repository, install a GoLang compiler and run `go build`.

## Usage
![Screenshot](./screeshot.png)
### Form
Run `morcego`, then reply with:

(as example to a form with a vulnerable input string)

> URL: The form URL

> Method: POST

> Value Type: STRING

> Target Column: The column name

> Vulnerable Input: The name of the vulnerable input

> Extra Inputs: Use if the form requires any other input

> Extra Condition: If you wanna limit the query, then use it

> Error Message: The expect error message found when the input value is false

**Wait and done**

### GET request with parameters in the end of the path
Run `morcego`, then reply with:

(as example to a vulnerable "REST API" to get a entry by the id: `http://localhost/users/1`)

> URL: The URL (without the parameter)

> Method: GET

> Value Type: INT

> Target Column: The column name

> Vulnerable Input: Leave it empty

> Extra Condition: If you wanna limit the query, then use it

> Error Message: The expect error message found when the input value is false

**Wait and done**

### GET request with parameters in the query
Run `morcego`, then reply with:

(as example to a vulnerable "API" to get a entry by the id in "id" parameter: `http://localhost/user?id=1`)

> URL: The URL (without the parameter)

> Method: GET

> Value Type: INT

> Target Column: The column name

> Vulnerable Input: id (use your own here)

> Extra Condition: If you wanna limit the query, then use it

> Error Message: The expect error message found when the input value is false

**Wait and done**

## The name
Morcego was writen during the COVID-19 pandemic, so it's named after the goddam bat that locked us inside our houses (and that's not even the bad thing about it).

**Morcego** = Bat (in portuguese)  
Mor**cego** = Blind (also in portuguese)  
Morce**go** = Go (the programming language used in it)

## License
`GPL-2.0`, for more information, see the [LICENSE](./LICENSE) file.
