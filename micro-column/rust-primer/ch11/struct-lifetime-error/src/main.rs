struct Excerpt {
    part: &str,
}

fn main() {
    let novel = String::from("事不过三。这一次，我们扬帆起航，绝不返航。");
    let first_sentence = novel.split('。').next().expect("找不到句号");

    let excerpt = Excerpt {
        part: first_sentence,
    };

    println!("摘录：{}", excerpt.part);
}
