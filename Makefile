include .env
export

ATLAS := atlas

# Apply migrations to local database
db-apply:
	@echo "💻 Applying migrations to local..."
	@$(ATLAS) schema apply --env gorm -u "mysql://myapp_user:your_strong_password@db:3306/myapp_db"
	@echo "✅ Migration applied successfully\n"

# Generate migration files from models
db-diff:
	@echo "💻 Generating migration files..."
	@$(ATLAS) migrate diff --env gorm

# Print environment variables for debugging
env:
	@echo "MYSQL_HOST: $(MYSQL_HOST)"
	@echo "MYSQL_PORT_INTERNAL: $(MYSQL_PORT_INTERNAL)"
	@echo "MYSQL_USER: $(MYSQL_USER)"
	@echo "MYSQL_DATABASE: $(MYSQL_DATABASE)" 