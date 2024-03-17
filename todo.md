profiler, results viewer
garbage collector (timing of when it runs)
errors?
structs?

- Mathematical operations (*, /, +, -) with bidmas (DONE)
- Comparisons (>, <, >=, <=, !=) (DONE)
- Not operator (!) (DONE)
- Arrays:
```
var arr = ["a", "b", "c"]
var arr [5]string = [] // 5 empty strings
var arr []string = ["abc"] // Unkown size array
arr[0] // Access
arr[0] = "d" // Assignment
```
(DONE)
- Maps
- Structs
```
struct Pet {
    name: string
    age: uint8

    fn ageInMonths(): uint8 {
        return self.age * 12
    }
}
```
- Loops:  
```
for i = range 5 {}
for i = range 2,8 {}

var arr = ["a", "b", "c"]
for i, el = range arr {}

var ages map[string]uint8 = {
    "joe": 28,
    "john": 32
}
for name, arge = range ages {}
```
(DONE)
- Error handling:
```
// ! after function specifies that it can return an error
fn myFunction(): string! { return "hi" }
fn myFunction()! {}

// If a function is marked as being able to return an error
fn myFunction(input: string): string! {
    if input == "yes" {
        return "nice"
    } else {
        return NewError("Expected yes.")
    }
}

// A caller of this function must then handle the potential error case
myFunction("no") catch(err) {

}

// Or if inside another function that returns an error it can use "try" to pass the error along to it's previous function
fn main()! {
    var result = try myFunction("no") // If the function returns an error, the error will be returned by main, if not the code will continue
}
main() catch(err) {
    print("Caught error", err)
}

// enums can be used for specific errors
// enums take the type uint16 by default but a type can also be passed. Errors expect uint16
enum MyFunctionError: uint16 {
    InputTooShort
    InputTooLong
}
fn doSomething(input: string)! {
    if len(input) < 3 {
        // NewError() creates a MutableError
        return NewError().code(MyFunctionError.InputTooShort)
    } else if len(input) > 5 {
        return NewError().code(MyFunctionError.InputTooLong)
    }
}
fn main(): string {
    doSomething(input()) catch(err) {
        // When the error is passed to a catch it is converted to a StaticError
        return match err.code {
            MyFunctionError.InputTooShort => "Too short!",
            MyFunctionError.InputTooLong => {
                print("Too long")
                return "Too long!"
            }
        }
    }
    return "No error."
}
```
- Imports, exports:
```
// other.lang
export fn hello() {
    print("Hello world!")
}
export msg = "hello" // Static values can be exported
// or
var msg = "hello"
export msg 
msg = "hi" // This has no impact on the export
```
```
// main.lang
import hello, msg from "./other.lang"

hello()
print(msg)
```
- Garbage collection:
Mark and sweep garbage collection, each node is going to need a method to return every identifier that it references
- key-value library
- profiler