// 具名字段结构体：字段各有名字，适合表示"有多个不同含义属性"的实体
#[derive(Debug)]
struct User {
    username: String,
    email: String,
    active: bool,
}

// 元组结构体：字段没有名字，只有位置和类型，适合"给一个已有类型包一层新含义"
#[derive(Debug)]
struct Meters(f64);

#[derive(Debug)]
struct Color(u8, u8, u8);

// 单元结构体：没有任何字段，适合"只需要一个类型标记，不需要携带数据"的场景
struct AlwaysEqual;

fn main() {
    let user = User {
        username: String::from("tony"),
        email: String::from("tony@example.com"),
        active: true,
    };
    println!("{:?}", user);

    let distance = Meters(42.0);
    println!("{:?}，取出内部值：{}", distance, distance.0);

    let bg = Color(255, 0, 0);
    println!("{:?}，红色分量：{}", bg, bg.0);

    let _marker = AlwaysEqual;
    println!("单元结构体没有字段，仅仅是一个类型标记");
}
