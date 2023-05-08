// import "functions.sl", We can't import files yet

print("Hello World!")

let x = {
    "name": "John",
    "nested": {
        "name": "Doe"
    }
}

print(x["name"])
print(x["nested"]["name"])

let y = [1, 2, 3, 4, 5]

let i = 10;
// FIXME: add own environment for for loops
for (i = 0; i < len(y); i++) {
    print(y[i])
}

print("I: ", i)


print(y)

//#region "Functions"
let reverse = magic(a, b) { quote(unquote(b) - unquote(a)) };
//#endregion

print(reverse(2 + 2, 10 - 5));