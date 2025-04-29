-- Criação da tabela de contas
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    balance NUMERIC DEFAULT 0
);

-- Inserção de dados iniciais
INSERT INTO accounts (name, balance) VALUES 
('Alice', 200.00),
('Bob', 150.00);

-- Índice para melhorar a performance de consultas por nome
CREATE INDEX IF NOT EXISTS idx_accounts_name ON accounts(name);

-- Função para registrar transações
CREATE OR REPLACE FUNCTION register_transaction(
    p_account_id INTEGER,
    p_amount NUMERIC,
    p_type TEXT
) RETURNS VOID AS $$
BEGIN
    IF p_type = 'deposit' THEN
        UPDATE accounts SET balance = balance + p_amount WHERE id = p_account_id;
    ELSIF p_type = 'withdraw' THEN
        UPDATE accounts SET balance = balance - p_amount WHERE id = p_account_id;
    END IF;
END;
$$ LANGUAGE plpgsql; 