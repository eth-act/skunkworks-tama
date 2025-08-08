use std::path::Path;
use zisk_core::elf2rom;

fn main() {
    println!("Testing ELF processing with debug output...\n");
    
    let elf_path = Path::new("elf-regressions/elf-output/empty.elf");
    
    println!("Processing: {}\n", elf_path.display());
    
    match elf2rom(elf_path) {
        Ok(rom) => {
            println!("\n=== ROM Processing Complete ===");
            println!("Total instructions in ROM: {}", rom.insts.len());
        }
        Err(e) => {
            eprintln!("Error processing ELF: {}", e);
        }
    }
}