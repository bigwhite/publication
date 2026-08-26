fn main() {
    let mut stock = 3;
    while stock != 0 {
        println!("剩余库存：{}", stock);
        stock -= 1;
    }
    println!("售罄！");

    let books = ["Rust 第一课", "所有权", "生命周期"];
    let mut index = 0;
    while index < books.len() {
        println!("（不推荐写法）第 {} 本：{}", index, books[index]);
        index += 1;
    }

    for book in books.iter() {
        println!("（地道写法）书名：{}", book);
    }

    for n in 1..=5 {
        print!("{} ", n);
    }
    println!();

    for n in (0..10).step_by(2) {
        print!("{} ", n);
    }
    println!();
}
