#!/bin/bash

# Configuration
log_file="testDir/app.log"      # Base log file name
backup_count=5          # Maximum number of backup files to keep

# Function to simulate rolling
roll_logs() {
    # Start from the oldest file and remove it if the maximum backup count is reached
    if [[ -f "${log_file}.${backup_count}" ]]; then
        echo "Deleting oldest log: ${log_file}.${backup_count}"
        rm "${log_file}.${backup_count}"
    fi

    # Shift all existing log files (e.g., app.log.4 -> app.log.5)
    for (( i=$backup_count-1; i>=1; i-- )); do
        if [[ -f "${log_file}.${i}" ]]; then
            echo "Renaming ${log_file}.${i} to ${log_file}.$((i+1))"
            mv "${log_file}.${i}" "${log_file}.$((i+1))"
        fi
    done

    # Roll the current log file (e.g., app.log -> app.log.1)
    if [[ -f "$log_file" ]]; then
        echo "Renaming $log_file to ${log_file}.1"
        mv "$log_file" "${log_file}.1"
    fi

    # Create a new log file (simulating new log entries)
    echo "Creating new log file: $log_file"
    echo "Log entry at $(date)" > "$log_file"
}

# Main loop to simulate log writing and rolling
for i in {1..10}; do
    echo "Simulating log entry $i"
    echo "Log entry $i at $(date)" >> "$log_file"

    # Roll logs every 3 iterations (for demonstration purposes)
    if (( i % 3 == 0 )); then
        echo "=== Rolling logs ==="
        roll_logs
        echo "===================="
    fi

    sleep 1  # Simulate time delay
done

