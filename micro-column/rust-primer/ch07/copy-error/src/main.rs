#[derive(Debug, Clone, Copy)]
struct Player {
    name: String,
    score: i32,
}

fn main() {
    let p = Player {
        name: String::from("Tony"),
        score: 100,
    };
    println!("{:?}", p);
}
