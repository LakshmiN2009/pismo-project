CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    document_number VARCHAR(20) NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS operation_types (
    operation_type_id SERIAL PRIMARY KEY,
    description VARCHAR(50) NOT NULL,
    transaction_type INT DEFAULT 0,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(account_id),
    operation_type_id INT REFERENCES operation_types(operation_type_id),
    amount float NOT NULL,
    event_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO operation_types(description) VALUES ('Normal Purchase'),('Purchase with installments'),('Withdrawal'),('Credit Voucher');