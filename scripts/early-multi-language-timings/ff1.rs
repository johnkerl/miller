// Claude

use std::fs::File;
use std::io::{self, BufRead, BufReader};

// Method 1: Using BufReader with lines() - Good balance of speed and memory
fn read_lines_buf_reader(path: &str) -> io::Result<()> {
    let file = File::open(path)?;
    let reader = BufReader::new(file);

    let mut n = 0;
    for line in reader.lines() {
        let _line = line?;
        // Process line here
        n += 1;
    }
    println!("line count {}", n);

    Ok(())
}

fn main() {
  let _ = read_lines_buf_reader("/Users/kerl/tmp/big");
}
