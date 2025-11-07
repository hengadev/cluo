# Database Migrations

This directory contains Goose database migrations for the document management system.

## Migration Structure

The migrations are organized by document domain for better readability and maintainability:

- `20250602201515_document_estimates_init_schema.sql` - Estimates table and related indexes
- `20250602201630_document_mandates_init_schema.sql` - Mandates table and related indexes
- `20250602201745_document_contracts_init_schema.sql` - Contracts table and related indexes
- `20250602201900_document_invoices_init_schema.sql` - Invoices table and related indexes
- `20250602202015_document_versions_init_schema.sql` - Document versions table for audit trail

## Using Goose

### Installation

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Running Migrations

1. **Set up your database connection string:**
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

2. **Run all pending migrations:**
```bash
goose -dir ./migrations postgres "$DATABASE_URL" up
```

3. **Run a specific migration:**
```bash
goose -dir ./migrations postgres "$DATABASE_URL" up 20250602201745
```

4. **Rollback migrations:**
```bash
goose -dir ./migrations postgres "$DATABASE_URL" down
```

5. **Rollback to a specific version:**
```bash
goose -dir ./migrations postgres "$DATABASE_URL" down-to 20250602201630
```

6. **Check migration status:**
```bash
goose -dir ./migrations postgres "$DATABASE_URL" status
```

7. **Create a new migration:**
```bash
goose -dir ./migrations create add_new_field sql
```

## Migration Order

The migrations should be run in the following order:

1. **Document Versions** - Creates the audit trail system
2. **Estimates** - Creates the estimates table
3. **Mandates** - Creates the mandates table (references estimates)
4. **Contracts** - Creates the contracts table (references mandates)
5. **Invoices** - Creates the invoices table (references contracts)

The timestamp-based naming ensures proper execution order when running `goose up`.

## Database Schema Overview

### Core Tables

#### `estimates`
- Stores price quotations for investigative services
- Contains line items, pricing, and acceptance status
- Links to cases and clients

#### `mandates`
- Stores legal authorization documents
- Contains signatures, scope of work, and validity periods
- Links to estimates for traceability

#### `contracts`
- Stores formal agreements between parties
- Contains multiple signatures, terms, and conditions
- Links to mandates for business flow continuity

#### `invoices`
- Stores billing documents for services rendered
- Contains payment processing and tax calculations
- Links to contracts for financial tracking

#### `document_versions`
- Stores audit trail for all document types
- Contains complete document history with author tracking
- Supports rollback and compliance requirements

### Key Features

- **UUID Primary Keys**: All tables use UUID primary keys for global uniqueness
- **Timestamp Triggers**: Automatic `updated_at` field updates
- **JSON Storage**: Flexible data storage for line items and signatures
- **Referential Integrity**: Foreign key constraints with proper cascade rules
- **Performance Indexes**: Optimized queries for common access patterns
- **Business Rules**: Database-level constraints for data integrity
- **Audit Trail**: Complete version history for compliance

## Environment-Specific Considerations

### Development
```bash
goose -dir ./migrations postgres "postgres://dev_user:dev_pass@localhost:5432/cluo_dev?sslmode=disable" up
```

### Staging
```bash
goose -dir ./migrations postgres "postgres://stage_user:stage_pass@staging-db:5432/cluo_staging?sslmode=require" up
```

### Production
```bash
goose -dir ./migrations postgres "postgres://prod_user:prod_pass@prod-db.internal:5432/cluo_prod?sslmode=require" up
```

## Troubleshooting

### Common Issues

1. **Migration Lock Issues**
   ```bash
   # Check for active locks
   SELECT * FROM goose_db_version;

   # Reset if needed (only in development!)
   TRUNCATE goose_db_version;
   ```

2. **Schema Validation Errors**
   ```bash
   # Check current schema version
   goose -dir ./migrations postgres "$DATABASE_URL" status
   ```

3. **Missing Extensions**
   ```sql
   -- Enable UUID extension manually if needed
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
   ```

### Data Migration Scripts

When modifying existing schemas, always create a new migration:

```bash
goose -dir ./migrations create add_document_metadata sql
```

Example migration content:
```sql
-- +goose Up
ALTER TABLE estimates ADD COLUMN metadata JSONB DEFAULT '{}';
CREATE INDEX idx_estimates_metadata ON estimates USING gin(metadata);

-- +goose Down
ALTER TABLE estimates DROP COLUMN metadata;
```

## Best Practices

1. **Always test migrations** on a copy of production data
2. **Use descriptive migration names** with timestamps
3. **Include both UP and DOWN** migrations
4. **Add comments** explaining complex changes
5. **Review foreign key cascades** carefully
6. **Test rollback procedures** regularly
7. **Document breaking changes** in the README

## Performance Considerations

- Large JSON fields (line_items, signatures) are indexed with GIN indexes
- Timestamp fields use TIMESTAMPTZ for timezone support
- Conditional indexes on frequently filtered columns
- Foreign key indexes for join performance

## Security Notes

- All migrations run with the database user's permissions
- No sensitive data is stored in migration files
- Consider using environment variables for connection strings
- Review generated SQL for potential security implications
