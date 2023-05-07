let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated;
    } 

    push(accumulated, f(first(arr)))
    iter(rest(arr), accumulated);
  };

  iter(arr, []);
};

let reduce = fn(arr, initial, f) {
  let iter = fn(arr, result) {
    if (len(arr) == 0) {
      return result;
    }

    iter(rest(arr), f(result, first(arr)));
  };

  iter(arr, initial);
};

let sum = fn(arr) {
  reduce(arr, 0, fn(initial, el) { initial + el });
};

let product = fn(arr) {
  reduce(arr, 1, fn(initial, el) { initial * el });
};

let filter = fn(arr, f) {
  let iter = fn(arr, result) {
    if (len(arr) == 0) {
      return result;
    }

    let x = first(arr);
    if (f(x)) {
      push(result, x);
    }

    iter(rest(arr), result);
  };

  iter(arr, []);
};


// We have no inbuilt modulo operator, so we define one
let mod = fn(a, b) { a - (b * (a / b)); };

print("Filtering even numbers: ", filter([1, 2, 3, 4], fn(x) { mod(x, 2) == 0 }));

