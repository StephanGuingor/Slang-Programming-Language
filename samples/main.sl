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

for (let i = 0; i < len(y); i++) {
    print(x["nested"]["name"], y[i])
}