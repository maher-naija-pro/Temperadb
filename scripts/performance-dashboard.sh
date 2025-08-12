#!/bin/bash

# TimeSeriesDB Performance Monitoring Dashboard
# This script generates performance reports and visualizations

set -e

# Configuration
RESULTS_DIR="benchmark-results"
DASHBOARD_DIR="performance-dashboard"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Function to print colored output
print_status() {
    echo -e "\033[0;34m[INFO]\033[0m $1"
}

print_success() {
    echo -e "\033[0;32m[SUCCESS]\033[0m $1"
}

print_warning() {
    echo -e "\033[1;33m[WARNING]\033[0m $1"
}

print_error() {
    echo -e "\033[0;31m[ERROR]\033[0m $1"
}

# Function to show help
show_help() {
    echo "TimeSeriesDB Performance Monitoring Dashboard"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -g, --generate          Generate performance dashboard"
    echo "  -t, --trends            Generate performance trends analysis"
    echo "  -s, --summary           Generate performance summary report"
    echo "  -o, --output DIR        Output directory (default: $DASHBOARD_DIR)"
    echo "  --days N                Analyze last N days of data (default: 30)"
    echo ""
    echo "Examples:"
    echo "  $0 -g                   # Generate full dashboard"
    echo "  $0 -t --days 7          # Generate trends for last 7 days"
    echo "  $0 -s                   # Generate summary report only"
}

# Function to generate performance summary
generate_summary() {
    local output_dir="$1"
    local days="$2"
    
    print_status "Generating performance summary..."
    
    local summary_file="$output_dir/performance_summary.txt"
    
    echo "=== TimeSeriesDB Performance Summary ===" > "$summary_file"
    echo "Generated: $(date)" >> "$summary_file"
    echo "Analysis period: Last $days days" >> "$summary_file"
    echo "" >> "$summary_file"
    
    # Count benchmark files
    local total_runs=$(ls "$RESULTS_DIR"/benchmark_*.txt 2>/dev/null | wc -l)
    echo "Total benchmark runs: $total_runs" >> "$summary_file"
    echo "" >> "$summary_file"
    
    # List all benchmark files
    echo "Benchmark files found:" >> "$summary_file"
    for file in "$RESULTS_DIR"/benchmark_*.txt; do
        if [[ -f "$file" ]]; then
            local file_size=$(du -h "$file" | cut -f1)
            local file_date=$(stat -c %y "$file" | cut -d' ' -f1)
            echo "  - $(basename "$file") ($file_size, $file_date)" >> "$summary_file"
        fi
    done
    
    print_success "Performance summary generated: $summary_file"
}

# Function to generate performance trends
generate_trends() {
    local output_dir="$1"
    local days="$2"
    
    print_status "Generating performance trends analysis..."
    
    local trends_file="$output_dir/performance_trends.txt"
    
    echo "=== Performance Trends Analysis ===" > "$trends_file"
    echo "Generated: $(date)" >> "$trends_file"
    echo "Analysis period: Last $days days" >> "$trends_file"
    echo "" >> "$trends_file"
    
    # Find regression reports
    local regression_count=$(ls "$RESULTS_DIR"/regression_report_*.txt 2>/dev/null | wc -l)
    echo "Regression reports found: $regression_count" >> "$trends_file"
    echo "" >> "$trends_file"
    
    if [[ $regression_count -gt 0 ]]; then
        echo "Recent regression reports:" >> "$trends_file"
        for report in "$RESULTS_DIR"/regression_report_*.txt; do
            if [[ -f "$report" ]]; then
                local report_date=$(stat -c %y "$report" | cut -d' ' -f1)
                echo "  - $(basename "$report") ($report_date)" >> "$trends_file"
            fi
        done
    fi
    
    print_success "Performance trends generated: $trends_file"
}

# Function to generate HTML dashboard
generate_html_dashboard() {
    local output_dir="$1"
    local days="$2"
    
    print_status "Generating HTML dashboard..."
    
    local html_file="$output_dir/index.html"
    
    # Count files for dashboard
    local total_runs=$(ls "$RESULTS_DIR"/benchmark_*.txt 2>/dev/null | wc -l)
    local regression_reports=$(ls "$RESULTS_DIR"/regression_report_*.txt 2>/dev/null | wc -l)
    
    cat > "$html_file" << EOF
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TimeSeriesDB Performance Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007acc; padding-bottom: 10px; }
        .metric { display: flex; justify-content: space-between; margin: 10px 0; padding: 10px; background: #f8f9fa; border-radius: 5px; }
        .metric-label { font-weight: bold; }
        .metric-value { color: #007acc; }
        .status-badge { padding: 4px 8px; border-radius: 12px; font-size: 0.8em; font-weight: bold; }
        .status-badge.success { background: #d5f4e6; color: #27ae60; }
        .status-badge.warning { background: #fef9e7; color: #f39c12; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 TimeSeriesDB Performance Dashboard</h1>
        <p><small>Last updated: $(date)</small></p>
        
        <h2>📊 Performance Overview</h2>
        <div class="metric">
            <span class="metric-label">Total Benchmark Runs</span>
            <span class="metric-value">$total_runs</span>
        </div>
        <div class="metric">
            <span class="metric-label">Regression Reports</span>
            <span class="metric-value">$regression_reports</span>
        </div>
        <div class="metric">
            <span class="metric-label">Analysis Period</span>
            <span class="metric-value">$days days</span>
        </div>
        
        <h2>🔧 Quick Actions</h2>
        <p>Run benchmarks: <code>./scripts/run-benchmarks.sh</code></p>
        <p>Detect regressions: <code>./scripts/detect-regressions.sh</code></p>
        <p>View reports: <code>benchmark-results/</code></p>
    </div>
</body>
</html>
EOF

    print_success "HTML dashboard generated: $html_file"
}

# Function to generate complete dashboard
generate_dashboard() {
    local output_dir="$1"
    local days="$2"
    
    print_status "Generating complete performance dashboard..."
    
    # Create directory structure
    if [ ! -d "$output_dir" ]; then
        mkdir -p "$output_dir"
        print_status "Created dashboard directory: $output_dir"
    fi
    
    # Generate components
    generate_summary "$output_dir" "$days"
    generate_trends "$output_dir" "$days"
    generate_html_dashboard "$output_dir" "$days"
    
    print_success "Complete dashboard generated in: $output_dir"
    print_status "Open $output_dir/index.html in your browser to view the dashboard"
}

# Parse command line arguments
GENERATE_DASHBOARD=false
GENERATE_TRENDS=false
GENERATE_SUMMARY=false
OUTPUT_DIR="$DASHBOARD_DIR"
DAYS=30

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -g|--generate)
            GENERATE_DASHBOARD=true
            shift
            ;;
        -t|--trends)
            GENERATE_TRENDS=true
            shift
            ;;
        -s|--summary)
            GENERATE_SUMMARY=true
            shift
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --days)
            DAYS="$2"
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
    print_status "TimeSeriesDB Performance Monitoring Dashboard"
    
    # Check if we're in the right directory
    if [ ! -d "$RESULTS_DIR" ]; then
        print_error "Results directory not found: $RESULTS_DIR"
        print_status "Make sure you're running this script from the project root"
        exit 1
    fi
    
    # Set default action if none specified
    if [[ "$GENERATE_DASHBOARD" = false && "$GENERATE_TRENDS" = false && "$GENERATE_SUMMARY" = false ]]; then
        GENERATE_DASHBOARD=true
    fi
    
    # Execute requested actions
    if [ "$GENERATE_SUMMARY" = true ]; then
        generate_summary "$OUTPUT_DIR" "$DAYS"
    fi
    
    if [ "$GENERATE_TRENDS" = true ]; then
        generate_trends "$OUTPUT_DIR" "$DAYS"
    fi
    
    if [ "$GENERATE_DASHBOARD" = true ]; then
        generate_dashboard "$OUTPUT_DIR" "$DAYS"
    fi
    
    print_success "Dashboard generation completed!"
    print_status "Output directory: $OUTPUT_DIR"
    
    if [ "$GENERATE_DASHBOARD" = true ]; then
        print_status "Open $OUTPUT_DIR/index.html in your browser to view the dashboard"
    fi
}

# Run main function
main "$@"
