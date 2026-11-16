use std::env;
use std::process;

use cligrep::Config;

fn main() {
    let args: Vec<String> = env::args().collect();

    let config = Config::build(&args).unwrap_or_else(|err| {
        eprintln!("参数错误：{}", err);
        process::exit(1);
    });

    if let Err(e) = cligrep::run(&config) {
        eprintln!("运行出错：{}", e);
        process::exit(1);
    }
}
