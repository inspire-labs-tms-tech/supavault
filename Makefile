# Variables
SCHEMA ?= public

# Default target
.PHONY: diff
diff:
	@echo "Running supabase db diff on schema: $(SCHEMA)"
	supabase db diff -s $(SCHEMA)

# Help target
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make diff         # Run supabase db diff on the 'public' schema by default"
	@echo "  make diff SCHEMA=your_schema  # Specify a different schema to diff"