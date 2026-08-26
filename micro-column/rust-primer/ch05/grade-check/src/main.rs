fn classify(score: u32) -> &'static str {
    if score > 100 {
        "无效分数"
    } else if score >= 90 {
        "优秀"
    } else if score >= 75 {
        "良好"
    } else if score >= 60 {
        "及格"
    } else {
        "不及格"
    }
}

fn main() {
    let scores = [95, 88, 76, 59, 60, 101];

    let mut excellent_count = 0;
    let mut index = 0;

    for &score in scores.iter() {
        let grade = classify(score);
        println!("第 {} 位同学，分数 {}，等级：{}", index + 1, score, grade);

        if grade == "优秀" {
            excellent_count += 1;
        }

        index += 1;
    }

    println!();
    println!("本次一共有 {} 位同学获得\"优秀\"", excellent_count);

    let summary = if excellent_count >= 2 {
        "本次测验整体表现不错"
    } else {
        "本次测验还有提升空间"
    };
    println!("总评：{}", summary);
}
