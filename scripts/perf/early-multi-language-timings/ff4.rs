// Claude

use std::fs::File;
use std::io::{self, BufRead, BufReader};

// Method 4: Memory mapped files for very large files
fn read_lines_memmap(path: &str) -> io::Result<()> {
    use memmap2::Mmap;
    
    let file = File::open(path)?;
    let mmap = unsafe { Mmap::map(&file)? };
    
    // Convert memory map to string and iterate over lines
    if let Ok(contents) = std::str::from_utf8(&mmap) {
        for line in contents.lines() {
            // Process line here
        }
    }
    Ok(())
}

fn main() {
  let _ = read_lines_custom_buffer("/Users/kerl/tmp/big");
}
