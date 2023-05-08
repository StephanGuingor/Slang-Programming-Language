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

let i = 0
for (i = 0; i < len(y); i++) {
    print(y[i])

    let j = 0
}

if (i == 5) {
    let x = 10
    print("i is 5")
} else {
    print("i is not 5")
}

print("x: ", x)

print("I: ", i)


print(y)

//#region "Functions"
let reverse = magic(a, b) { quote(unquote(b) - unquote(a)) };
//#endregion

print(reverse(2 + 2, 10 - 5));


if (true) {
    let z = 20
    print(z)
}

print(z)