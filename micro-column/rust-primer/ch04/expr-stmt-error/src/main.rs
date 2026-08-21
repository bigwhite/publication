fn main() {
    let y = {
        let inner = 5 * 2;
        inner + 1;
    };
    println!("y = {:?}", y);
}
