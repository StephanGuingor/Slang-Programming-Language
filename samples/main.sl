let createAdder = fn (x) {
    fn (y) {
        return x + y;
    }
}

let add5 = createAdder(5);
let add10 = createAdder(10);

let apply = fn (funcA, funcB) {
    return funcA(funcB);
}
 