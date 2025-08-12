#!/bin/bash

# TimeSeriesDB Benchmark Runner Script
# This script helps run benchmarks locally and compare results

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BENCHMARK_DIR="./test"
RESULTS_DIR="benchmark-results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
CURRENT_RESULTS="$RESULTS_DIR/benchmark_$TIMESTAMP.txt"
BASELINE_RESULTS="$RESULTS_DIR/baseline.txt"

# Create results directory if it doesn't exist
mkdir -p "$RESULTS_DIR"

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

# Function to show help
show_help() {
    echo "TimeSeriesDB Benchmark Runner"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -a, --all               Run all benchmarks (default)"
    	echo "  -p, --ingestion         Run ingestion benchmarks only"
    echo "  -s, --storage           Run storage benchmarks only"
    echo "  -e, --http              Run HTTP endpoint benchmarks only"
    echo "  -t, --e2e               Run end-to-end benchmarks only"
    echo "  -m, --memory            Run memory usage benchmarks only"
    echo "  -c, --compare           Compare with baseline results"
    echo "  -b, --baseline          Set current results as baseline"
    echo "  -o, --output FILE       Output file (default: benchmark_YYYYMMDD_HHMMSS.txt)"
    echo "  -v, --verbose           Verbose output"
    echo "  -t, --timeout SECONDS   Timeout for benchmarks (default: 10m)"
    echo ""
    echo "Examples:"
    echo "  $0                      # Run all benchmarks"
    	echo "  $0 -p                   # Run ingestion benchmarks only"
    echo "  $0 -c                   # Compare with baseline"
    echo "  $0 -b                   # Set current results as baseline"
    echo "  $0 -o my_results.txt    # Save to specific file"
}

# Function to run benchmarks
run_benchmarks() {
    local benchmark_type="$1"
    local output_file="$2"
    local timeout="$3"
    
    print_status "Running benchmarks: $benchmark_type"
    print_status "Output file: $output_file"
    print_status "Timeout: $timeout"
    
    case "$benchmark_type" in
        "all")
            go test -bench=. -benchmem -timeout="$timeout" -v "$BENCHMARK_DIR" | tee "$output_file"
            ;;
        	"ingestion")
            go test -bench=BenchmarkParse -benchmem -timeout="$timeout" -v "$BENCHMARK_DIR" | tee "$output_file"
            ;;
        "storage")
            go test -bench=BenchmarkWrite -benchmem -timeout="$timeout" -v "$BENCHMARK_DIR" | tee "$output_file"
            ;;
        "http")
            go test -bench=BenchmarkHTTP -benchmem -timeout="$timeout" -v "$BENCHMARK_DIR" | tee "$output_file"
            ;;
        "e2e")
            go test -bench="BenchmarkEndToEnd|BenchmarkConcurrent" -benchmem -timeout="$timeout" -v "$BENCHMARK_DIR" | tee "$output_file"
            ;;
        "memory")
            go test -bench=BenchmarkMemory -benchmem -timeout="$timeout" -v "$BENCHMARK_DIR" | tee "$output_file"
            ;;
        *)
            print_error "Unknown benchmark type: $benchmark_type"
            exit 1
            ;;
    esac
    
    print_success "Benchmarks completed and saved to $output_file"
}

# Function to compare results
compare_results() {
    local current_file="$1"
    local baseline_file="$2"
    
    if [ ! -f "$baseline_file" ]; then
        print_warning "No baseline file found at $baseline_file"
        print_status "Run '$0 -b' to set a baseline first"
        return 1
    fi
    
    if [ ! -f "$current_file" ]; then
        print_error "Current results file not found: $current_file"
        return 1
    fi
    
    print_status "Comparing current results with baseline..."
    
    # Extract benchmark results and compare
    echo "=== Benchmark Comparison ===" > "$RESULTS_DIR/comparison_$TIMESTAMP.txt"
    echo "Current: $current_file" >> "$RESULTS_DIR/comparison_$TIMESTAMP.txt"
    echo "Baseline: $baseline_file" >> "$RESULTS_DIR/comparison_$TIMESTAMP.txt"
    echo "" >> "$RESULTS_DIR/comparison_$TIMESTAMP.txt"
    
    # Use a simpler approach to compare results
    while IFS= read -r line; do
        if [[ $line =~ ^Benchmark ]]; then
            # Extract benchmark name and ns/op value
            benchmark_name=$(echo "$line" | awk '{print $1}')
            current_ns=$(echo "$line" | awk '{print $3}' | sed 's/[^0-9.]//g')
            
            if [[ -n "$current_ns" ]]; then
                # Find corresponding baseline result
                baseline_line=$(grep "^$benchmark_name" "$baseline_file" | head -1)
                if [[ -n "$baseline_line" ]]; then
                    baseline_ns=$(echo "$baseline_line" | awk '{print $3}' | sed 's/[^0-9.]//g')
                    
                    if [[ -n "$baseline_ns" ]]; then
                        # Calculate difference
                        diff=$(echo "$current_ns - $baseline_ns" | bc -l 2>/dev/null || echo "0")
                        pct=$(echo "scale=2; ($diff / $baseline_ns) * 100" | bc -l 2>/dev/null || echo "0")
                        
                        if (( $(echo "$diff < 0" | bc -l) )); then
                            status="✅ FASTER"
                        elif (( $(echo "$diff > 0" | bc -l) )); then
                            status="❌ SLOWER"
                        else
                            status="➡️  SAME"
                        fi
                        
                        printf "%-25s | %15s | %15s | %+10.2f%% | %s\n" \
                               "$benchmark_name" "$current_ns" "$baseline_ns" "$pct" "$status" \
                               >> "$RESULTS_DIR/comparison_$TIMESTAMP.txt"
                    fi
                fi
            fi
        fi
    done < "$current_file"
    
    print_success "Comparison saved to $RESULTS_DIR/comparison_$TIMESTAMP.txt"
    cat "$RESULTS_DIR/comparison_$TIMESTAMP.txt"
}

# Function to set baseline
set_baseline() {
    local current_file="$1"
    
    if [ ! -f "$current_file" ]; then
        print_error "Current results file not found: $current_file"
        print_status "Run benchmarks first to generate results"
        return 1
    fi
    
    cp "$current_file" "$BASELINE_RESULTS"
    print_success "Baseline set to: $BASELINE_RESULTS"
}

# Parse command line arguments
BENCHMARK_TYPE="all"
OUTPUT_FILE="$CURRENT_RESULTS"
TIMEOUT="10m"
COMPARE_MODE=false
SET_BASELINE=false
VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -a|--all)
            BENCHMARK_TYPE="all"
            shift
            ;;
        	-p|--ingestion)
		BENCHMARK_TYPE="ingestion"
            shift
            ;;
        -s|--storage)
            BENCHMARK_TYPE="storage"
            shift
            ;;
        -e|--http)
            BENCHMARK_TYPE="http"
            shift
            ;;
        -t|--e2e)
            BENCHMARK_TYPE="e2e"
            shift
            ;;
        -m|--memory)
            BENCHMARK_TYPE="memory"
            shift
            ;;
        -c|--compare)
            COMPARE_MODE=true
            shift
            ;;
        -b|--baseline)
            SET_BASELINE=true
            shift
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        --timeout)
            TIMEOUT="$2"
            shift 2
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
    print_status "TimeSeriesDB Benchmark Runner"
    print_status "Working directory: $(pwd)"
    
    # Check if we're in the right directory
    if [ ! -d "$BENCHMARK_DIR" ]; then
        print_error "Benchmark directory not found: $BENCHMARK_DIR"
        print_status "Make sure you're running this script from the project root"
        exit 1
    fi
    
    # Check if Go is available
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Download dependencies
    print_status "Downloading Go dependencies..."
    go mod tidy
    
    if [ "$COMPARE_MODE" = true ]; then
        compare_results "$OUTPUT_FILE" "$BASELINE_RESULTS"
        exit 0
    fi
    
    if [ "$SET_BASELINE" = true ]; then
        # If output file was specified, use it; otherwise use the most recent results file
        if [ "$OUTPUT_FILE" != "$CURRENT_RESULTS" ]; then
            set_baseline "$OUTPUT_FILE"
        else
            # Find the most recent benchmark results file
            latest_file=$(ls -t "$RESULTS_DIR"/benchmark_*.txt 2>/dev/null | head -1)
            if [ -n "$latest_file" ]; then
                set_baseline "$latest_file"
            else
                print_error "No benchmark results found. Run benchmarks first."
                exit 1
            fi
        fi
        exit 0
    fi
    
    # Run benchmarks
    run_benchmarks "$BENCHMARK_TYPE" "$OUTPUT_FILE" "$TIMEOUT"
    
    print_success "Benchmark run completed!"
    print_status "Results saved to: $OUTPUT_FILE"
    print_status "To compare with baseline: $0 -c"
    print_status "To set as baseline: $0 -b"
}

# Run main function
main "$@"
