use std::io;
use std::io::BufRead;

fn main() {
    for line in io::stdin().lock().lines() {
        print!("{}", line.unwrap());
    }
}

//fn main() {
//    let mut reader = io::stdin();
//    let mut line;
//    loop {
//        line = reader.read_line();
//        print!("{}\n", line);
//    }
//}
