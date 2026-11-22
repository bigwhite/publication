fn largest<T: PartialOrd + Copy>(list: &[T]) -> T {
    let mut largest = list[0];
    for &item in list {
        if item > largest {
            largest = item;
        }
    }
    largest
}

fn main() {
    let numbers = vec![34, 50, 25];
    let chars = vec!['y', 'm', 'a'];

    // 编译器会分别为 T=i32 和 T=char 各生成一份专门的机器码，
    // 就好像我们手写了 largest_i32 和 largest_char 两个函数一样，
    // 这个过程叫"单态化"（monomorphization），发生在编译期，运行期没有任何额外开销。
    println!("{}", largest(&numbers));
    println!("{}", largest(&chars));
}
