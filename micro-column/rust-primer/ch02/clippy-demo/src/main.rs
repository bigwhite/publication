fn is_even(n: i32) -> bool {
    if n % 2 == 0 {
        return true;
    } else {
        return false;
    }
}

fn main() {
    let numbers = vec![1, 2, 3, 4, 5, 6];
    let mut even_count = 0;

    for i in 0..numbers.len() {
        if is_even(numbers[i]) == true {
            even_count = even_count + 1;
        }
    }

    println!("偶数个数：{}", even_count);
}
