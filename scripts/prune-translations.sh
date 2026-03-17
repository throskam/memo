#!/bin/bash

# Prunes outdated translations from messages.gotext.json files
# Keeps only translations that are present in out.gotext.json

for dir in internal/translations/locales/*/; do
    messages_file="${dir}messages.gotext.json"
    out_file="${dir}out.gotext.json"
    
    if [[ ! -f "$messages_file" || ! -f "$out_file" ]]; then
        continue
    fi
    
    echo "Pruning $messages_file..."
    
    # Create temporary file
    tmp_file=$(mktemp)
    
    # Filter messages using both files
    jq -s '
        .[0] as $msgs | .[1] as $out_msgs | 
        [$out_msgs.messages[].id] as $out_ids |
        [$msgs.messages[] | select(.id | . as $id | $out_ids | index($id) // false)] as $filtered_msgs |
        {"language": $msgs.language, "messages": $filtered_msgs}
    ' "$messages_file" "$out_file" > "$tmp_file"
    
    if [[ $? -eq 0 ]]; then
        mv "$tmp_file" "$messages_file"
        echo "  Done"
    else
        echo "  Failed - keeping original file"
        rm -f "$tmp_file"
    fi
done
