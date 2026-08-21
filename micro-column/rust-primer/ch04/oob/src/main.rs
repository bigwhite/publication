fn get_score(scores: [i32; 5], index: usize) -> i32 {
    scores[index]
}

fn main() {
    let scores = [90, 85, 77, 88, 95];
    let idx = scores.len(); // 数组长度是 5，合法下标是 0..4，5 已经越界了
    println!("{}", get_score(scores, idx));
}
