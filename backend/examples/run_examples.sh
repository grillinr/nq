#!/bin/bash

# Examples Runner Script
# This script helps run the NQ backend examples

set -e

echo "NQ Backend Examples"
echo "=================="
echo

# Check if we're in the right directory
if [ ! -d "examples" ]; then
    echo "Error: This script must be run from the backend directory"
    echo "Usage: ./examples/run_examples.sh"
    exit 1
fi

# Function to run an example
run_example() {
    local example_name=$1
    local example_dir="examples/$example_name"
    
    if [ ! -d "$example_dir" ]; then
        echo "Error: Example '$example_name' not found in $example_dir"
        return 1
    fi
    
    echo "Running $example_name example..."
    echo "================================"
    cd "$example_dir"
    go run main.go
    cd - > /dev/null
    echo
}

# Show menu if no arguments provided
if [ $# -eq 0 ]; then
    echo "Available examples:"
    echo "  1. metadata     - Fetch metadata for movies, books, TV shows, games"
    echo "  2. integrations - Sync data from external services (Spotify, Steam, etc.)"
    echo "  3. all          - Run all examples"
    echo
    read -p "Enter your choice (1-3): " choice
    
    case $choice in
        1)
            run_example "metadata"
            ;;
        2)
            run_example "integrations"
            ;;
        3)
            echo "Running all examples..."
            echo
            run_example "metadata"
            run_example "integrations"
            ;;
        *)
            echo "Invalid choice: $choice"
            exit 1
            ;;
    esac
else
    # Run specific example(s) from command line arguments
    for example in "$@"; do
        case $example in
            "metadata"|"integrations")
                run_example "$example"
                ;;
            "all")
                run_example "metadata"
                run_example "integrations"
                ;;
            *)
                echo "Unknown example: $example"
                echo "Available examples: metadata, integrations, all"
                exit 1
                ;;
        esac
    done
fi

echo "Done!"