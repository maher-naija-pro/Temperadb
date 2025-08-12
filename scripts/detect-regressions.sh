#!/bin/bash

# TimeSeriesDB Performance Regression Detection Script
# This script analyzes benchmark results and detects performance regressions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
ORANGE='\033[0;33m'
NC='\033[0m' # No Color

# Configuration
RESULTS_DIR="benchmark-results"
BASELINE_FILE="$RESULTS_DIR/baseline.txt"
REGRESSION_THRESHOLD=5.0  # 5% performance regression threshold
CRITICAL_THRESHOLD=15.0   # 15% critical regression threshold
OUTPUT_FILE="$RESULTS_DIR/regression_report_$(date +"%Y%m%d_%H%M%S").txt"
HTML_REPORT="$RESULTS_DIR/regression_report_$(date +"%Y%m%d_%H%M%S").html"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_critical() {
    echo -e "${ORANGE}[CRITICAL]${NC} $1"
}

# Function to show help
show_help() {
    echo "TimeSeriesDB Performance Regression Detection"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -b, --baseline FILE     Baseline file to compare against (default: $BASELINE_FILE)"
    echo "  -c, --current FILE      Current results file to analyze"
    echo "  -t, --threshold PERCENT Performance regression threshold (default: $REGRESSION_THRESHOLD%)"
    echo "  -o, --output FILE       Output file for regression report (default: auto-generated)"
    echo "  -v, --verbose           Verbose output"
    echo "  -j, --json              Output results in JSON format"
    echo "  -H, --html              Generate HTML report"
    echo "  --critical-threshold    Critical regression threshold (default: $CRITICAL_THRESHOLD%)"
    echo ""
    echo "Examples:"
    echo "  $0 -c latest_results.txt                    # Analyze specific file against baseline"
    echo "  $0 -c latest_results.txt -t 10              # Use 10% threshold"
    echo "  $0 -c latest_results.txt -j                 # JSON output"
    echo "  $0 -c latest_results.txt -H                 # Generate HTML report"
}

# Function to find the most recent benchmark results file
find_latest_results() {
    local latest_file=$(ls -t "$RESULTS_DIR"/benchmark_*.txt 2>/dev/null | head -1)
    if [[ -z "$latest_file" ]]; then
        print_error "No benchmark results found in $RESULTS_DIR"
        print_status "Run benchmarks first using: ./scripts/run-benchmarks.sh"
        exit 1
    fi
    echo "$latest_file"
}

# Function to parse benchmark results
parse_benchmark_results() {
    local file="$1"
    local -A results
    
    while IFS= read -r line; do
        if [[ $line =~ ^Benchmark ]]; then
            # Extract benchmark name, ns/op, and B/op
            # Format: BenchmarkName-16      	 1343606	       895.9 ns/op	     880 B/op	      12 allocs/op
            if [[ $line =~ ^([^[:space:]]+)[[:space:]]+([0-9]+)[[:space:]]+([0-9.]+)[[:space:]]ns/op[[:space:]]+([0-9]+)[[:space:]]B/op[[:space:]]+([0-9]+)[[:space:]]allocs/op ]]; then
                benchmark_name="${BASH_REMATCH[1]}"
                ns_per_op="${BASH_REMATCH[3]}"
                bytes_per_op="${BASH_REMATCH[4]}"
                allocs_per_op="${BASH_REMATCH[5]}"
                
                if [[ -n "$ns_per_op" && "$ns_per_op" != "0" ]]; then
                    results["$benchmark_name"]="$ns_per_op|$bytes_per_op|$allocs_per_op"
                fi
            fi
        fi
    done < "$file"
    
    echo "${results[@]}"
}

# Function to detect regressions
detect_regressions() {
    local baseline_file="$1"
    local current_file="$2"
    local threshold="$3"
    local critical_threshold="$4"
    
    print_status "Analyzing performance regressions..."
    print_status "Baseline: $baseline_file"
    print_status "Current: $current_file"
    print_status "Threshold: ${threshold}%"
    print_status "Critical threshold: ${critical_threshold}%"
    
    # Parse both files
    local -A baseline_results
    local -A current_results
    
    # Parse baseline results
    while IFS= read -r line; do
        if [[ $line =~ ^Benchmark ]]; then
            # Format: BenchmarkName-16      	 1343606	       895.9 ns/op	     880 B/op	      12 allocs/op
            if [[ $line =~ ^([^[:space:]]+)[[:space:]]+([0-9]+)[[:space:]]+([0-9.]+)[[:space:]]ns/op[[:space:]]+([0-9]+)[[:space:]]B/op[[:space:]]+([0-9]+)[[:space:]]allocs/op ]]; then
                benchmark_name="${BASH_REMATCH[1]}"
                ns_per_op="${BASH_REMATCH[3]}"
                bytes_per_op="${BASH_REMATCH[4]}"
                allocs_per_op="${BASH_REMATCH[5]}"
                
                if [[ -n "$ns_per_op" && "$ns_per_op" != "0" ]]; then
                    baseline_results["$benchmark_name"]="$ns_per_op|$bytes_per_op|$allocs_per_op"
                fi
            fi
        fi
    done < "$baseline_file"
    
    # Parse current results
    while IFS= read -r line; do
        if [[ $line =~ ^Benchmark ]]; then
            # Format: BenchmarkName-16      	 1343606	       895.9 ns/op	     880 B/op	      12 allocs/op
            if [[ $line =~ ^([^[:space:]]+)[[:space:]]+([0-9]+)[[:space:]]+([0-9.]+)[[:space:]]ns/op[[:space:]]+([0-9]+)[[:space:]]B/op[[:space:]]+([0-9]+)[[:space:]]allocs/op ]]; then
                benchmark_name="${BASH_REMATCH[1]}"
                ns_per_op="${BASH_REMATCH[3]}"
                bytes_per_op="${BASH_REMATCH[4]}"
                allocs_per_op="${BASH_REMATCH[5]}"
                
                if [[ -n "$ns_per_op" && "$ns_per_op" != "0" ]]; then
                    current_results["$benchmark_name"]="$ns_per_op|$bytes_per_op|$allocs_per_op"
                fi
            fi
        fi
    done < "$current_file"
    
    # Generate regression report
    local regression_count=0
    local critical_count=0
    local improvement_count=0
    
    echo "=== Performance Regression Analysis ===" > "$OUTPUT_FILE"
    echo "Generated: $(date)" >> "$OUTPUT_FILE"
    echo "Baseline: $baseline_file" >> "$OUTPUT_FILE"
    echo "Current: $current_file" >> "$OUTPUT_FILE"
    echo "Threshold: ${threshold}%" >> "$OUTPUT_FILE"
    echo "Critical threshold: ${critical_threshold}%" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    
    echo "Benchmark Name                    | Baseline (ns/op) | Current (ns/op) | Change (%) | Status" >> "$OUTPUT_FILE"
    echo "----------------------------------|------------------|-----------------|------------|---------" >> "$OUTPUT_FILE"
    
    # Check each benchmark
    for benchmark in "${!current_results[@]}"; do
        if [[ -n "${baseline_results[$benchmark]}" ]]; then
            local baseline_ns=$(echo "${baseline_results[$benchmark]}" | cut -d'|' -f1)
            local current_ns=$(echo "${current_results[$benchmark]}" | cut -d'|' -f1)
            
            if [[ -n "$baseline_ns" && -n "$current_ns" ]]; then
                local change_percent=$(echo "scale=2; ($current_ns - $baseline_ns) / $baseline_ns * 100" | bc -l 2>/dev/null || echo "0")
                
                local status="✅"
                if (( $(echo "$change_percent > $threshold" | bc -l) )); then
                    if (( $(echo "$change_percent > $critical_threshold" | bc -l) )); then
                        status="🚨 CRITICAL"
                        ((critical_count++))
                    else
                        status="⚠️  REGRESSION"
                        ((regression_count++))
                    fi
                elif (( $(echo "$change_percent < -$threshold" | bc -l) )); then
                    status="🚀 IMPROVEMENT"
                    ((improvement_count++))
                fi
                
                printf "%-35s | %15.0f | %15.0f | %+10.2f%% | %s\n" \
                       "$benchmark" "$baseline_ns" "$current_ns" "$change_percent" "$status" \
                       >> "$OUTPUT_FILE"
            fi
        fi
    done
    
    echo "" >> "$OUTPUT_FILE"
    echo "=== Summary ===" >> "$OUTPUT_FILE"
    echo "Total benchmarks analyzed: ${#current_results[@]}" >> "$OUTPUT_FILE"
    echo "Performance regressions: $regression_count" >> "$OUTPUT_FILE"
    echo "Critical regressions: $critical_count" >> "$OUTPUT_FILE"
    echo "Performance improvements: $improvement_count" >> "$OUTPUT_FILE"
    
    # Print summary to console
    echo ""
    print_status "=== Regression Analysis Summary ==="
    print_status "Total benchmarks: ${#current_results[@]}"
    
    if [[ $regression_count -gt 0 ]]; then
        print_warning "Performance regressions: $regression_count"
    else
        print_success "No performance regressions detected"
    fi
    
    if [[ $critical_count -gt 0 ]]; then
        print_critical "Critical regressions: $critical_count"
    fi
    
    if [[ $improvement_count -gt 0 ]]; then
        print_success "Performance improvements: $improvement_count"
    fi
    
    # Set exit code based on regressions
    if [[ $critical_count -gt 0 ]]; then
        exit 2  # Critical regression
    elif [[ $regression_count -gt 0 ]]; then
        exit 1  # Performance regression
    else
        exit 0  # No regressions
    fi
}

# Function to generate HTML report
generate_html_report() {
    local output_file="$1"
    local html_file="$2"
    
    print_status "Generating HTML report: $html_file"
    
    cat > "$html_file" << EOF
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Performance Regression Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007acc; padding-bottom: 10px; }
        .summary { background: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .summary h3 { margin-top: 0; color: #495057; }
        .regression { color: #dc3545; font-weight: bold; }
        .improvement { color: #28a745; font-weight: bold; }
        .critical { color: #fd7e14; font-weight: bold; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #007acc; color: white; }
        tr:nth-child(even) { background-color: #f2f2f2; }
        .status { font-weight: bold; }
        .timestamp { color: #6c757d; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 TimeSeriesDB Performance Regression Report</h1>
        <div class="timestamp">Generated: $(date)</div>
        
        <div class="summary">
            <h3>📊 Analysis Summary</h3>
            <p>This report analyzes performance changes between benchmark runs to detect regressions.</p>
        </div>
        
        <h2>📋 Detailed Results</h2>
        <pre>$(cat "$output_file")</pre>
    </div>
</body>
</html>
EOF

    print_success "HTML report generated: $html_file"
}

# Function to output JSON results
output_json() {
    local output_file="$1"
    
    # This is a simplified JSON output - you could enhance this to parse the actual results
    cat > "${output_file%.txt}.json" << EOF
{
  "report": {
    "generated": "$(date -Iseconds)",
    "baseline_file": "$BASELINE_FILE",
    "current_file": "$CURRENT_FILE",
    "threshold": $REGRESSION_THRESHOLD,
    "critical_threshold": $CRITICAL_THRESHOLD,
    "output_file": "$OUTPUT_FILE"
  },
  "summary": {
    "message": "Performance regression analysis completed. Check the output file for detailed results."
  }
}
EOF

    print_success "JSON output generated: ${output_file%.txt}.json"
}

# Parse command line arguments
CURRENT_FILE=""
VERBOSE=false
JSON_OUTPUT=false
HTML_OUTPUT=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -b|--baseline)
            BASELINE_FILE="$2"
            shift 2
            ;;
        -c|--current)
            CURRENT_FILE="$2"
            shift 2
            ;;
        -t|--threshold)
            REGRESSION_THRESHOLD="$2"
            shift 2
            ;;
        --critical-threshold)
            CRITICAL_THRESHOLD="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -j|--json)
            JSON_OUTPUT=true
            shift
            ;;
        -H|--html)
            HTML_OUTPUT=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
main() {
    print_status "TimeSeriesDB Performance Regression Detection"
    
    # Check if we're in the right directory
    if [ ! -d "$RESULTS_DIR" ]; then
        print_error "Results directory not found: $RESULTS_DIR"
        print_status "Make sure you're running this script from the project root"
        exit 1
    fi
    
    # Check if baseline exists
    if [ ! -f "$BASELINE_FILE" ]; then
        print_error "Baseline file not found: $BASELINE_FILE"
        print_status "Set a baseline first using: ./scripts/run-benchmarks.sh -b"
        exit 1
    fi
    
    # Find current results file if not specified
    if [ -z "$CURRENT_FILE" ]; then
        CURRENT_FILE=$(find_latest_results)
        print_status "Using latest results file: $CURRENT_FILE"
    fi
    
    # Check if current results file exists
    if [ ! -f "$CURRENT_FILE" ]; then
        print_error "Current results file not found: $CURRENT_FILE"
        exit 1
    fi
    
    # Run regression detection
    detect_regressions "$BASELINE_FILE" "$CURRENT_FILE" "$REGRESSION_THRESHOLD" "$CRITICAL_THRESHOLD"
    
    # Generate additional outputs if requested
    if [ "$HTML_OUTPUT" = true ]; then
        generate_html_report "$OUTPUT_FILE" "$HTML_REPORT"
    fi
    
    if [ "$JSON_OUTPUT" = true ]; then
        output_json "$OUTPUT_FILE"
    fi
    
    print_success "Regression detection completed!"
    print_status "Report saved to: $OUTPUT_FILE"
    
    if [ "$HTML_OUTPUT" = true ]; then
        print_status "HTML report: $HTML_REPORT"
    fi
    
    if [ "$JSON_OUTPUT" = true ]; then
        print_status "JSON output: ${OUTPUT_FILE%.txt}.json"
    fi
}

# Run main function
main "$@"
