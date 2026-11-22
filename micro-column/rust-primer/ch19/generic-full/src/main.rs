fn largest<T: PartialOrd + Copy>(list: &[T]) -> T {
    let mut largest = list[0];
    for &item in list {
        if item > largest {
            largest = item;
        }
    }
    largest
}

// 泛型结构体：两个字段共用同一个类型参数 T
#[derive(Debug)]
struct Point<T> {
    x: T,
    y: T,
}

// 泛型结构体：两个字段可以是不同类型 T 和 U
#[derive(Debug)]
struct Pair<T, U> {
    first: T,
    second: U,
}

// 泛型枚举：我们熟悉的 Option<T> 和 Result<T, E> 本身就是标准库定义的泛型枚举，
// 这里我们自己定义一个类似语义的简化版，感受一下语法
#[derive(Debug)]
enum MyOption<T> {
    MySome(T),
    MyNone,
}

fn main() {
    let numbers = vec![34, 50, 25, 100, 65];
    println!("最大的数字：{}", largest(&numbers));

    let chars = vec!['y', 'm', 'a', 'q'];
    println!("最大的字符：{}", largest(&chars));

    let floats = vec![1.5, 3.2, 0.8];
    println!("最大的浮点数：{}", largest(&floats));

    let int_point = Point { x: 5, y: 10 };
    let float_point = Point { x: 1.0, y: 4.0 };
    println!("{:?}", int_point);
    println!("{:?}", float_point);

    let pair = Pair {
        first: "分数",
        second: 95,
    };
    println!("{:?}", pair);

    let some_val: MyOption<i32> = MyOption::MySome(42);
    let none_val: MyOption<i32> = MyOption::MyNone;
    println!("{:?}", some_val);
    println!("{:?}", none_val);
}
