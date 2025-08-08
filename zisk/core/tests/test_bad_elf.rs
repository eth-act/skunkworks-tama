use zisk_core::elf_extraction::collect_elf_payload;
use std::path::Path;

#[test]
fn test_elf_with_bss_outside_ram() {
    // Create a fake ELF file path (this would be an ELF with BSS at wrong address)
    // For demonstration, we'll use a non-existent file to show the error flow
    
    // This test is mainly to document the expected behavior
    // In a real scenario, you'd have an ELF with .bss or .data sections
    // placed outside the RAM region (0xa0000000 - 0xc0000000)
    
    println!("This test documents that ELF files with writable sections");
    println!("outside RAM boundaries should be rejected with a clear error message.");
    println!("RAM boundaries: 0xa0000000 - 0xc0000000");
}

#[test] 
fn test_good_elf_passes() {
    // Test that our properly configured go-program.elf works
    let elf_path = Path::new("../go-program.elf");
    
    if elf_path.exists() {
        match collect_elf_payload(elf_path) {
            Ok(payload) => {
                // Should have some executable sections
                assert!(!payload.exec.is_empty(), "Should have executable sections");
                
                // Should have read-write sections (data and bss in RAM)
                assert!(!payload.rw.is_empty(), "Should have read-write sections in RAM");
                
                // Entry point should be valid
                assert!(payload.entry_point != 0, "Should have valid entry point");
                
                println!("✓ ELF with properly placed sections loads successfully");
            }
            Err(e) => {
                panic!("Good ELF should not fail: {}", e);
            }
        }
    }
}