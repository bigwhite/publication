// 省略写法：编译器靠"规则二"自动帮我们补全
fn first_word(s: &str) -> &str {
    let bytes = s.as_bytes();
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }
    s
}

// 完全等价的显式写法
fn first_word_explicit<'a>(s: &'a str) -> &'a str {
    let bytes = s.as_bytes();
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }
    s
}

struct Excerpt<'a> {
    part: &'a str,
}

impl<'a> Excerpt<'a> {
    // 省略写法：靠"规则三"自动把 &self 的生命周期赋给返回值
    fn part(&self) -> &str {
        self.part
    }
}

fn main() {
    let sentence = String::from("Rust 第一课 事不过三");
    println!("first_word: {}", first_word(&sentence));
    println!("first_word_explicit: {}", first_word_explicit(&sentence));

    let excerpt = Excerpt {
        part: "摘录内容"
    };
    println!("excerpt.part(): {}", excerpt.part());
}
