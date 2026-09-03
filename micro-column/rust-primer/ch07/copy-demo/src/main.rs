#[derive(Debug, Clone, Copy)]
struct Point {
    x: i32,
    y: i32,
}

#[derive(Debug, Clone)]
struct Player {
    name: String,
    score: i32,
}

fn main() {
    let p1 = Point { x: 1, y: 2 };
    let p2 = p1;
    println!("p1 = {:?}", p1);
    println!("p2 = {:?}", p2);
    println!("p1 依然可用，因为 Point 实现了 Copy trait");

    println!();

    let player1 = Player {
        name: String::from("Tony"),
        score: 100,
    };
    let player2 = player1.clone();
    println!("player1 = {:?}", player1);
    println!("player2 = {:?}", player2);
    println!("player1 依然可用，但这次是靠显式 .clone()，而不是 Copy");
}
