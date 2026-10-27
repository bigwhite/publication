fn main() {
    let mut stack = vec![1, 2, 3, 4, 5];

    // while let：只要 pop() 还能拿到 Some，就继续循环
    while let Some(top) = stack.pop() {
        println!("弹出：{}", top);
    }
    println!("栈已清空：{:?}", stack);

    // 对比：如果用 match + loop 手写，会啰嗦很多
    let mut stack2 = vec![10, 20, 30];
    loop {
        match stack2.pop() {
            Some(top) => println!("（手写版）弹出：{}", top),
            None => break,
        }
    }
}
