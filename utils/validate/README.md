# validate

validate is a small tool manage your validators of string as a chain.

# Example
```Go
var wrongLength = errors.New("wrong length")
var wrongChars = errors.New("wrong chars")
var Email = validate.Use(
    validate.Length{
        Min: 3,
        Max: 30,
        Err: wrongLength,
    }.Validate,
    validate.Chars{
        Chars: "abcdefghijklmn@",
        Err:   wrongChars,
    }.Validate,
)
var Name = validate.Use(
    validate.Length{
        Min: 3,
        Max: 16,
        Err: wrongLength,
    }.Validate,
    validate.Chars{
        Chars: "abcdefg",
        Err:   wrongChars,
    }.Validate,
)
var Password = validate.Length{
    Min: 8,
    Max: 16,
    Err: wrongLength,
}.Validate

var User = validate.UseMul(Email, Name, Password)

Email("ab") == wrongLength
Email("abcde12345") == wrongChars

Name("12345") == wrongChars

Password("111") == wrongLength

//     Email    Name  Password Password
//      |         |      |       |
User("abcde", "12345", "111", "ddd") == ?
```
