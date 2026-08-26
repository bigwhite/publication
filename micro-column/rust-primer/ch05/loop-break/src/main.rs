fn main() {
    let mut counter = 0;

    let result = loop {
        counter += 1;

        if counter == 10 {
            break counter * 2;
        }
    };

    println!("result = {}", result);

    // 带标签的嵌套循环
    let mut count = 0;
    'outer: loop {
        let mut remaining = 10;

        loop {
            println!("remaining = {}", remaining);
            if remaining == 9 {
                break;
            }
            if count == 2 {
                break 'outer;
            }
            remaining -= 1;
        }

        count += 1;
    }
    println!("count = {}", count);
}
