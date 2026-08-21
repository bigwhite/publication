fn main() {
    let a: u8 = 255;
    let b = a.wrapping_add(1);
    println!("wrapping_add: {}", b);

    let c: u8 = 200;
    let d = c.checked_add(100);
    println!("checked_add: {:?}", d);

    let e: u8 = 250;
    let f = e.saturating_add(100);
    println!("saturating_add: {}", f);
}
