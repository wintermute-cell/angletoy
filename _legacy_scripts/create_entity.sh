#!/bin/bash

# Check if an argument is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 some_name"
    exit 1
fi

some_name="$1"

# Convert some_name to SomeName format
# This uses sed to capitalize each segment of the name separated by underscores
SomeName=$(echo "$some_name" | sed -r 's/(^|_)([a-z])/\U\2/g')

# Copy the template file to the new location
cp internal/core/entities/entity.template "pkg/entities/entity_${some_name}.go"

# Replace all occurrences of "Template" with "SomeName" in the new file
sed -i "s/Template/$SomeName/g" "pkg/entities/entity_${some_name}.go"

echo "Processed template for $SomeName and saved to pkg/entities/entity_${some_name}.go"

