use greeter_core::build_greeting;
use rand::Rng;

fn main() {
    let names = ["Rustacean", "Ferris", "Tony"];
    let mut rng = rand::thread_rng();
    let idx = rng.gen_range(0..names.len());
    let picked = names[idx];

    let message = build_greeting(picked);
    println!("{}", message);
}
