// Claude

use std::fs::File;
use std::io::{self, BufRead, BufReader};

// Method 2: Using BufReader with read_line - More control, similar performance
fn read_lines_manual(path: &str) -> io::Result<()> {
    let file = File::open(path)?;
    let mut reader = BufReader::new(file);
    let mut line = String::new();
    
    while reader.read_line(&mut line)? > 0 {
        // Process line here
        line.clear(); // Reuse the string buffer
    }
    Ok(())
}

fn main() {
  let _ = read_lines_manual("/Users/kerl/tmp/big");
}
