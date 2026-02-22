# Implementation Plan: Vietnamese Image-Based Printing

## Overview

This implementation plan converts the ESC/POS printer from character encoding (TCVN3) to image-based printing for Vietnamese text. The approach maintains the existing Printer interface while replacing the internal text-to-ESC/POS conversion with a text-to-image-to-ESC/POS pipeline.

## Tasks

- [ ] 1. Set up image rendering infrastructure
  - [x] 1.1 Create FontManager component for loading Vietnamese fonts
    - Implement font loading from file paths
    - Implement system font discovery (Arial Unicode MS, Roboto, DejaVu Sans)
    - Add error handling for missing fonts
    - _Requirements: 4.1, 4.2, 4.5_
  
  - [ ]* 1.2 Write unit tests for FontManager
    - Test loading valid TTF files
    - Test error handling for missing fonts
    - Test system font discovery
    - _Requirements: 4.1, 4.2, 4.5_
  
  - [x] 1.3 Create configuration structures (RendererConfig, FontConfig)
    - Define RendererConfig with pixel width, font path, font size, line spacing, margin
    - Define FontConfig with font paths, normal size, bold size
    - _Requirements: 10.1, 10.2, 10.3, 10.4_
  
  - [ ]* 1.4 Write property test for configuration application
    - **Property 24: Configuration Application**
    - **Validates: Requirements 10.1, 10.2, 10.3, 10.4, 10.5**

- [ ] 2. Implement FormatParser component
  - [x] 2.1 Create FormatParser with line format detection
    - Implement LineFormat struct (text, bold, alignment, isSeparator)
    - Implement Parse method to process text content
    - Implement detectAlignment, detectBold, isSeparator methods
    - _Requirements: 5.1, 5.2, 5.3, 5.4_
  
  - [ ]* 2.2 Write property tests for format detection
    - **Property 13: Receipt Header Formatting**
    - **Property 14: Receipt Item Formatting**
    - **Property 15: Receipt Total Formatting**
    - **Property 16: Receipt Footer Formatting**
    - **Validates: Requirements 5.1, 5.2, 5.3, 5.4**
  
  - [ ]* 2.3 Write unit tests for separator and empty line handling
    - Test separator line detection (=== and ---)
    - Test empty line handling
    - _Requirements: 2.5, 2.6_

- [ ] 3. Implement TextRenderer component
  - [x] 3.1 Create TextRenderer with basic rendering capability
    - Implement NewTextRenderer with font loading
    - Implement Render method to create grayscale images
    - Implement measureText for width calculation
    - Implement calculateImageHeight for height calculation
    - _Requirements: 1.1, 1.4_
  
  - [ ]* 3.2 Write property test for Vietnamese text rendering
    - **Property 1: Vietnamese Text Rendering**
    - **Property 2: Vietnamese Character Support**
    - **Validates: Requirements 1.1, 1.2**
  
  - [x] 3.3 Implement text alignment in TextRenderer
    - Implement drawLine method with alignment support (left, center, right)
    - Calculate X position based on alignment and text width
    - _Requirements: 2.2, 2.3, 2.4_
  
  - [ ]* 3.4 Write property test for text alignment
    - **Property 7: Text Alignment**
    - **Validates: Requirements 2.2, 2.3, 2.4**
  
  - [x] 3.5 Implement bold text rendering
    - Use bold font face for bold lines
    - Ensure bold text is visibly different from normal text
    - _Requirements: 2.1, 4.4_
  
  - [ ]* 3.6 Write property test for bold rendering
    - **Property 6: Bold Text Rendering**
    - **Validates: Requirements 2.1, 4.4**
  
  - [x] 3.7 Implement text wrapping at word boundaries
    - Add word wrapping logic in drawLine method
    - Ensure wrapping occurs at spaces, not mid-word
    - Handle lines that exceed pixel width
    - _Requirements: 1.5, 7.4_
  
  - [ ]* 3.8 Write property test for text wrapping
    - **Property 5: Text Wrapping at Word Boundaries**
    - **Validates: Requirements 1.5, 7.4**
  
  - [ ]* 3.9 Write property tests for line break and spacing preservation
    - **Property 4: Line Break Preservation**
    - **Property 17: Vertical Spacing Preservation**
    - **Validates: Requirements 1.4, 5.5**

- [ ] 4. Checkpoint - Ensure text rendering tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 5. Implement ImageConverter component
  - [x] 5.1 Create ImageConverter with ESC/POS conversion
    - Implement NewImageConverter
    - Implement ConvertToESCPOS method using GS v 0 format
    - Implement imageToRaster for pixel-to-byte conversion
    - Implement applyThreshold for monochrome conversion
    - _Requirements: 3.1, 3.2, 3.3, 9.2_
  
  - [ ]* 5.2 Write property tests for image conversion
    - **Property 8: Image to ESC/POS Conversion**
    - **Property 9: ESC/POS Command Format**
    - **Property 10: Byte Alignment**
    - **Property 11: Image Width Handling**
    - **Validates: Requirements 3.1, 3.2, 3.3, 3.4**
  
  - [ ]* 5.3 Write property tests for optimization and quality
    - **Property 12: Byte Stream Optimization**
    - **Property 21: Monochrome Threshold**
    - **Property 23: Memory Constraint Compliance**
    - **Validates: Requirements 3.5, 9.2, 9.4**
  
  - [ ]* 5.4 Write unit test for DPI setting
    - Test 203 DPI resolution is used
    - _Requirements: 9.5_

- [ ] 6. Implement pixel width calculation
  - [x] 6.1 Add pixel width calculation logic
    - Implement calculation: (paper_width_mm / 25.4) * 203 DPI
    - Support 58mm → 384 pixels and 80mm → 576 pixels
    - _Requirements: 1.3, 7.1, 7.2_
  
  - [ ]* 6.2 Write property test for pixel width calculation
    - **Property 3: Pixel Width Calculation**
    - **Validates: Requirements 1.3**
  
  - [ ]* 6.3 Write unit tests for specific paper widths
    - Test 58mm → 384 pixels
    - Test 80mm → 576 pixels
    - _Requirements: 7.1, 7.2_

- [ ] 7. Update ESCPOSPrinter to use image-based rendering
  - [x] 7.1 Modify NewESCPOSPrinter constructor
    - Initialize FormatParser, TextRenderer, ImageConverter
    - Load fonts during initialization
    - Calculate pixel width from paper width
    - Handle initialization errors
    - _Requirements: 6.1, 7.5_
  
  - [x] 7.2 Replace convertToESCPOS method implementation
    - Remove TCVN3 encoding logic
    - Implement new flow: Parse → Render → Convert
    - Add ESC_INIT at start and GS_CUT at end
    - _Requirements: 6.2_
  
  - [x] 7.3 Update Print method to use new convertToESCPOS
    - Ensure Print method calls updated convertToESCPOS
    - Maintain existing error handling
    - _Requirements: 6.2_
  
  - [ ]* 7.4 Write property tests for Print method
    - **Property 18: Print Method Execution**
    - **Property 19: Configuration Reading**
    - **Validates: Requirements 6.2, 7.5**
  
  - [ ]* 7.5 Write unit tests for Connect, Disconnect, GetStatus
    - Test connection establishment
    - Test disconnection
    - Test status retrieval
    - _Requirements: 6.3, 6.4, 6.5_

- [ ] 8. Implement error handling
  - [x] 8.1 Add comprehensive error handling throughout
    - Wrap errors with context using fmt.Errorf
    - Add descriptive error messages for all error types
    - Ensure errors propagate with context
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_
  
  - [ ]* 8.2 Write unit tests for error cases
    - Test font loading failure
    - Test image rendering failure
    - Test image conversion failure
    - Test printer communication failure
    - Test empty content validation
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

- [ ] 9. Implement image quality optimizations
  - [x] 9.1 Add anti-aliasing to text rendering
    - Configure font rendering with anti-aliasing enabled
    - _Requirements: 9.1_
  
  - [ ]* 9.2 Write property test for anti-aliasing
    - **Property 20: Anti-Aliasing Application**
    - **Validates: Requirements 9.1**
  
  - [x] 9.3 Optimize image height calculation
    - Calculate tight bounds around content
    - Add minimal padding (e.g., 10 pixels top/bottom)
    - _Requirements: 9.3_
  
  - [ ]* 9.4 Write property test for image height optimization
    - **Property 22: Image Height Optimization**
    - **Validates: Requirements 9.3**

- [ ] 10. Add ESC/POS command constants
  - [x] 10.1 Define GS v 0 command constants
    - Add GS_V_0 constant (0x1D, 0x76, 0x30)
    - Add IMAGE_MODE_NORMAL constant (0x00)
    - Add pixel width constants (PIXEL_WIDTH_58MM, PIXEL_WIDTH_80MM)
    - Add DPI constant (203)
    - _Requirements: 3.2, 9.5_

- [ ] 11. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 12. Integration and final wiring
  - [x] 12.1 Ensure all components are properly integrated
    - Verify FormatParser → TextRenderer → ImageConverter flow
    - Verify ESCPOSPrinter uses all components correctly
    - Test with sample receipt content
    - _Requirements: 6.2_
  
  - [ ]* 12.2 Write integration tests
    - Test full workflow with mock printer
    - Test with various receipt templates
    - Verify ESC/POS byte stream format
    - _Requirements: 6.2_
  
  - [x] 12.3 Add font file deployment support
    - Document required font files
    - Add font file paths to configuration
    - Provide fallback to system fonts
    - _Requirements: 4.1, 4.5_

- [ ] 13. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties (minimum 100 iterations each)
- Unit tests validate specific examples and edge cases
- Use `gopter` or `rapid` library for property-based testing in Go
- Each property test must include a comment tag: `// Feature: vietnamese-image-printing, Property {number}: {property_text}`
- Font files (TTF) must be available either in the deployment package or as system fonts
- Test with actual ZyWell ZY303 printer hardware before production deployment
