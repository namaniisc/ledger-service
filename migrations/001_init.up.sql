CREATE TABLE customers (
    customer_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    balance NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
    transaction_id UUID PRIMARY KEY,
    customer_id UUID REFERENCES customers(customer_id),
    type TEXT NOT NULL CHECK (type IN ('credit','debit')),
    amount NUMERIC NOT NULL,
    client_transaction_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_transactions_client_txn_id ON transactions(client_transaction_id) WHERE client_transaction_id IS NOT NULL;
