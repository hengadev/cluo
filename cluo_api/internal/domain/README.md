# Document Domain Implementation

This directory contains the complete implementation of the document management system for investigative services, following the specification in `DOMAIN_DESIGN.md`.

## Architecture Overview

The document domain follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer                               │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   Handlers      │  │     Routes      │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Application Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │     Service     │  │  Business Logic │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                 Domain Layer                                │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   Entities      │  │    Value Objects│                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│               Infrastructure Layer                          │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  Repositories   │  │   Database      │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

## Document Types

### 1. Estimate (`domain/estimate.go`)
- **Purpose**: Price quotations for investigative services
- **Key Features**: Line items, validation, acceptance workflow
- **Status Flow**: Draft → Sent → Signed → Active/Rejected

### 2. Mandate (`domain/mandate.go`)
- **Purpose**: Legal authorization for investigation work
- **Key Features**: Client/investigator signatures, scope of work, validity period
- **Status Flow**: Draft → Sent → Signed → Active/Expired

### 3. Contract (`domain/contract.go`)
- **Purpose**: Formal agreements between parties
- **Key Features**: Multiple signatures, terms and conditions, renewal options
- **Status Flow**: Draft → Sent → Signed → Active/Expired

### 4. Invoice (`domain/invoice.go`)
- **Purpose**: Billing documents for services rendered
- **Key Features**: Payment processing, tax calculations, late fees
- **Status Flow**: Draft → Sent → Paid/Overdue/Refunded

## Core Components

### Base Types

#### DocumentBase (`domain/document_base.go`)
- Shared metadata for all document types
- ID, CaseID, ClientID, Status, timestamps
- Common operations and validation

#### Signature (`domain/signature.go`)
- Digital and wet signature representation
- Signer information, method, timestamps
- Audit trail capabilities

#### DocumentStatus (`domain/document_status.go`)
- Enum for document lifecycle states
- Status transition validation
- JSON serialization support

### Business Logic

#### Document Workflow
The typical document progression follows this pattern:

```
Estimate → Mandate → Contract → Invoice
    ↓         ↓         ↓        ↓
  Draft → Sent → Signed → Active → Archived
```

#### Document Linking
- Estimates can be linked to Mandates
- Mandates can be linked to Contracts
- Contracts can be linked to Invoices
- Maintains complete audit trail

### Persistence Layer

#### Database Schema (`sql/001_documents_schema.sql`)
- PostgreSQL tables for each document type
- Proper indexing and constraints
- JSON storage for flexible data structures
- Automated triggers for timestamp updates

#### Repositories (`internal/adapters/postgres/document/`)
- Type-safe database operations
- Comprehensive CRUD functionality
- Version history management
- Document linking queries

### Service Layer (`internal/application/document/`)

#### Core Services
- Document creation and management
- Workflow orchestration
- Business rule enforcement
- Version tracking

#### Specialized Services
- **Estimate Service**: Creation, acceptance, line item management
- **Mandate Service**: Signing, activation, linking
- **Contract Service**: Multi-party signatures, activation
- **Invoice Service**: Payment processing, status management

### HTTP Layer (`internal/adapters/http/document/`)

#### REST API Endpoints

**Generic Document Operations:**
```
GET    /documents                    # List documents with filtering
POST   /documents                    # Create new document
GET    /documents/{id}/{type}         # Get specific document
PATCH  /documents/{id}/{type}         # Update document
DELETE /documents/{id}/{type}         # Delete document
POST   /documents/{id}/{type}/send    # Send document
POST   /documents/{id}/{type}/sign    # Sign document
POST   /documents/{id}/{type}/archive # Archive document
GET    /documents/{id}/{type}/history # Get document history
```

**Estimate Operations:**
```
POST   /estimates                     # Create estimate
PATCH  /estimates/{id}                # Update estimate
POST   /estimates/{id}/accept         # Accept estimate → Create Mandate
```

**Mandate Operations:**
```
POST   /mandates                      # Create mandate
POST   /mandates/{id}/sign            # Sign mandate
POST   /mandates/{id}/activate        # Activate mandate
POST   /mandates/{id}/create-contract # Create Contract from Mandate
```

**Contract Operations:**
```
POST   /contracts                     # Create contract
POST   /contracts/{id}/sign           # Sign contract
POST   /contracts/{id}/activate       # Activate contract
POST   /contracts/{id}/create-invoice # Create Invoice from Contract
```

**Invoice Operations:**
```
POST   /invoices                      # Create invoice
POST   /invoices/{id}/pay             # Process payment
POST   /invoices/{id}/void            # Void invoice
GET    /invoices/overdue              # Get overdue invoices
```

## Key Features

### 1. Document Versioning
- Automatic version creation on all changes
- Complete audit trail with author tracking
- Reason logging for all modifications
- Rollback capabilities

### 2. Business Rules Enforcement
- Status transition validation
- Document linking constraints
- Financial calculations validation
- Signature requirements

### 3. Workflow Automation
- Estimate acceptance → Mandate creation
- Mandate activation → Contract generation
- Contract completion → Invoice creation
- Payment processing → Status updates

### 4. Comprehensive Validation
- Domain-level validation
- Database constraints
- API request validation
- Business rule validation

## Testing

### Integration Tests (`test/integration/document/`)

The test suite includes:

1. **Estimate Creation and Acceptance**
   - Line item management
   - Total calculations
   - Acceptance workflow

2. **Mandate Signing Flow**
   - Client signature capture
   - Mandate activation
   - Status transitions

3. **Contract Creation from Mandate**
   - Document linking
   - Data inheritance
   - Workflow continuity

4. **Invoice Creation from Contract**
   - Financial calculations
   - Tax processing
   - Payment handling

5. **Complete Document Workflow**
   - End-to-end testing
   - Multi-document integration
   - Business process validation

### Running Tests

```bash
# Run all document tests
go test ./test/integration/document/...

# Run specific test
go test -run TestCompleteDocumentWorkflow ./test/integration/document/...
```

## Usage Examples

### Creating an Estimate

```go
// Create line items
lineItems := []domain.EstimateItem{
    {
        Description: "Surveillance services (8 hours)",
        Quantity:    8,
        UnitPrice:   50.00,
    },
}

// Create estimate
estimate := domain.NewEstimate(caseID, clientID, "EST-2025-001", lineItems)

// Save estimate
createdEstimate, err := service.CreateEstimate(ctx, estimate)
```

### Accepting Estimate and Creating Mandate

```go
// Accept estimate
mandate, err := service.AcceptEstimate(ctx, estimateID, userID)

// Mandate is automatically created and linked to estimate
```

### Processing Invoice Payment

```go
paymentReq := &domain.PaymentRequest{
    Amount:        1000.00,
    PaymentMethod: "wire_transfer",
    Notes:         "Payment for investigation services",
}

updatedInvoice, err := service.ProcessPayment(ctx, invoiceID, paymentReq)
```

## Configuration

### Database Setup

1. Execute the schema:
```sql
\i sql/001_documents_schema.sql
```

2. Configure connection:
```go
pool, err := pgxpool.New(ctx, "your-connection-string")
```

### Service Initialization

```go
// Initialize repositories
repo := documentRepository.New(pool)
versionRepo := documentRepository.NewVersionRepository(pool)

// Initialize service
service := document.New(repo, versionRepo)

// Register HTTP routes
document.RegisterRoutes(router, service)
```

## Future Enhancements

### Planned Features
- [ ] PDF generation and storage
- [ ] Email/SMS notifications
- [ ] Advanced filtering and search
- [ ] Document templates
- [ ] Bulk operations
- [ ] Document export (PDF, Excel)
- [ ] Advanced reporting
- [ ] Integration with external signature services
- [ ] Automated payment processing
- [ ] Document collaboration features

### Performance Optimizations
- [ ] Caching layer
- [ ] Database query optimization
- [ ] Pagination improvements
- [ ] Background job processing

### Security Enhancements
- [ ] Role-based access control
- [ ] Document encryption
- [ ] Audit logging
- [ ] Data retention policies

## Dependencies

- `github.com/google/uuid` - UUID generation
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/go-chi/chi/v5` - HTTP routing
- `github.com/stretchr/testify` - Testing utilities

## Contributing

When contributing to the document domain:

1. Follow existing code patterns and architecture
2. Add comprehensive tests for new features
3. Update documentation for API changes
4. Ensure all validation rules are implemented
5. Test complete workflows when making changes

## Support

For questions or issues related to the document domain:
- Review the test cases for usage examples
- Check the API documentation
- Examine the domain models for business rules
- Consult the database schema for data relationships