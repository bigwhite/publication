use strsplit_modular::StrSplit;

fn main() {
    let parts: Vec<&str> = StrSplit::new("Rust,第一课,事不过三", ",").collect();
    println!("{:?}", parts);
}
