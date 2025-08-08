use std::fs::File;
use std::io::Write;
use std::path::Path;

mod elf_extraction;
use elf_extraction::{collect_elf_payload, ElfPayload, DataSection};

fn format_data_section(section: &DataSection) -> String {
    let mut output = String::new();
    output.push_str(&format!("    Address: 0x{:08x}\n", section.addr));
    output.push_str(&format!("    Size: {} bytes (0x{:x})\n", section.data.len(), section.data.len()));
    
    // Show first and last few bytes for reference
    if section.data.len() > 0 {
        let preview_len = 32.min(section.data.len());
        let preview: Vec<String> = section.data[..preview_len]
            .iter()
            .map(|b| format!("{:02x}", b))
            .collect();
        output.push_str(&format!("    First {} bytes: {}\n", preview_len, preview.join(" ")));
        
        if section.data.len() > 32 {
            let tail_start = section.data.len().saturating_sub(16);
            let tail: Vec<String> = section.data[tail_start..]
                .iter()
                .map(|b| format!("{:02x}", b))
                .collect();
            output.push_str(&format!("    Last {} bytes: {}\n", section.data.len() - tail_start, tail.join(" ")));
        }
    }
    output
}

fn main() {
    println!("Reading go-program.elf...");
    
    let elf_path = Path::new("/home/kev/work/zisk/go-program.elf");
    
    match collect_elf_payload(elf_path) {
        Ok(payload) => {
            let mut output = String::new();
            
            output.push_str("=== ELF PAYLOAD ANALYSIS ===\n\n");
            output.push_str(&format!("Entry Point: 0x{:08x}\n\n", payload.entry_point));
            
            // Executable sections
            output.push_str(&format!("EXECUTABLE SECTIONS (count: {}):\n", payload.exec.len()));
            output.push_str("----------------------------------------\n");
            for (i, section) in payload.exec.iter().enumerate() {
                output.push_str(&format!("Exec Section #{}:\n", i + 1));
                output.push_str(&format_data_section(section));
                output.push_str("\n");
            }
            
            // Read-write sections
            output.push_str(&format!("READ-WRITE SECTIONS IN RAM (count: {}):\n", payload.rw.len()));
            output.push_str("----------------------------------------\n");
            for (i, section) in payload.rw.iter().enumerate() {
                output.push_str(&format!("RW Section #{}:\n", i + 1));
                output.push_str(&format_data_section(section));
                output.push_str("\n");
            }
            
            // Read-only sections
            output.push_str(&format!("READ-ONLY SECTIONS (count: {}):\n", payload.ro.len()));
            output.push_str("----------------------------------------\n");
            for (i, section) in payload.ro.iter().enumerate() {
                output.push_str(&format!("RO Section #{}:\n", i + 1));
                output.push_str(&format_data_section(section));
                output.push_str("\n");
            }
            
            // Summary statistics
            output.push_str("=== SUMMARY ===\n");
            let total_exec_size: usize = payload.exec.iter().map(|s| s.data.len()).sum();
            let total_rw_size: usize = payload.rw.iter().map(|s| s.data.len()).sum();
            let total_ro_size: usize = payload.ro.iter().map(|s| s.data.len()).sum();
            
            output.push_str(&format!("Total executable size: {} bytes\n", total_exec_size));
            output.push_str(&format!("Total read-write size: {} bytes\n", total_rw_size));
            output.push_str(&format!("Total read-only size: {} bytes\n", total_ro_size));
            output.push_str(&format!("Total size: {} bytes\n", total_exec_size + total_rw_size + total_ro_size));
            
            // Write to file
            let output_path = "elf_payload_analysis.txt";
            let mut file = File::create(output_path).expect("Failed to create output file");
            file.write_all(output.as_bytes()).expect("Failed to write to file");
            
            println!("Analysis complete! Results written to: {}", output_path);
            println!("\n{}", output);
        }
        Err(e) => {
            eprintln!("Error reading ELF file: {}", e);
        }
    }
}